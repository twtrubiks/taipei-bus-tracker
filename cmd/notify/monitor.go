package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/twtrubiks/taipei-bus-tracker/internal/handler"
	"github.com/twtrubiks/taipei-bus-tracker/internal/model"
)

// calcInterval returns the polling interval based on current ETA and threshold.
func calcInterval(etaSeconds, thresholdMinutes int) time.Duration {
	if etaSeconds < 0 {
		return 60 * time.Second
	}
	thresholdSec := thresholdMinutes * 60
	switch {
	case etaSeconds > thresholdSec*2:
		return 60 * time.Second
	case etaSeconds > thresholdSec:
		return 30 * time.Second
	default:
		return 15 * time.Second
	}
}

// notifyState tracks threshold-crossing for notification deduplication.
type notifyState struct {
	wasAboveThreshold bool
}

func newNotifyState() *notifyState {
	return &notifyState{wasAboveThreshold: true}
}

// check returns true if a notification should be sent.
// It detects the transition from above-threshold to within-threshold (positive ETA only).
// Resets on: ETA > threshold, arriving (0), or special status (negative).
func (s *notifyState) check(etaSeconds, thresholdMinutes int) bool {
	thresholdSec := thresholdMinutes * 60

	// Reset conditions: arriving (0), special status (<0), above threshold
	if etaSeconds <= 0 || etaSeconds > thresholdSec {
		s.wasAboveThreshold = true
		return false
	}

	// ETA within threshold (0 < eta <= threshold)
	if s.wasAboveThreshold {
		s.wasAboveThreshold = false
		return true
	}
	return false
}

// findStopETA returns the ETA in seconds for the target stop.
// Matches by StopID first (TDX), then by Sequence (eBus).
func findStopETA(etas []model.StopETA, stop model.Stop) int {
	if stop.StopID != "" {
		for _, e := range etas {
			if e.StopID == stop.StopID {
				return e.ETA
			}
		}
	}
	for _, e := range etas {
		if e.Sequence > 0 && e.Sequence == stop.Sequence {
			return e.ETA
		}
	}
	return model.ETANotDeparted
}

// formatETA returns a human-readable string for the given ETA value.
func formatETA(etaSeconds int) string {
	switch etaSeconds {
	case 0:
		return "進站中"
	case model.ETANotDeparted:
		return "未發車"
	case model.ETALastBusLeft:
		return "末班駛離"
	case model.ETANoStop:
		return "交管不停靠"
	case model.ETANotOperating:
		return "未營運"
	default:
		if etaSeconds < 0 {
			return "未知狀態"
		}
		return fmt.Sprintf("ETA %d 分", etaSeconds/60)
	}
}

// logETA prints a timestamped log line for the current polling result.
func logETA(etaSeconds int, notified bool) {
	now := time.Now().Format("15:04:05")
	etaStr := formatETA(etaSeconds)
	if notified {
		fmt.Printf("%s  %s  🔔 已通知！\n", now, etaStr)
	} else {
		fmt.Printf("%s  %s\n", now, etaStr)
	}
}

// sendNotification sends a desktop notification using the specified tool.
// Returns nil if tool is empty (no notification tool available).
func sendNotification(tool, routeName, stopName string, etaMinutes int) error {
	title := fmt.Sprintf("🚌 %s", routeName)
	body := fmt.Sprintf("%s - %d 分鐘到站", stopName, etaMinutes)

	var cmd *exec.Cmd
	switch tool {
	case "notify-send":
		cmd = exec.Command("notify-send", "--urgency=critical", title, body)
	case "kdialog":
		cmd = exec.Command("kdialog", "--passivepopup", fmt.Sprintf("%s\n%s", title, body), "10")
	default:
		return nil
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() { _ = cmd.Wait() }()
	playSound(hasPaplay)
	return nil
}

// hasPaplay is set at startup by detectSoundTool().
var hasPaplay bool

// detectSoundTool checks if paplay is available (call once at startup).
func detectSoundTool() {
	_, err := exec.LookPath("paplay")
	hasPaplay = err == nil
}

// playSound plays a notification sound via paplay (non-blocking).
func playSound(available bool) {
	if !available {
		return
	}
	cmd := exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga")
	if err := cmd.Start(); err == nil {
		go func() { _ = cmd.Wait() }()
	}
}

// runMonitor starts the ETA monitoring loop. It blocks until ctx is cancelled.
func runMonitor(ctx context.Context, svc *handler.FallbackService, routeID string, direction int, stop model.Stop, routeName string, threshold int, notifyCmd string) int {
	state := newNotifyState()
	startTime := time.Now()
	notifyCount := 0

	interval := 30 * time.Second
	timer := time.NewTimer(0) // fire immediately on start
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			elapsed := time.Since(startTime)
			fmt.Printf("\n─────────────────────────────────────────────\n")
			fmt.Printf("監控結束 | 時間: %s | 通知: %d 次\n", formatDuration(elapsed), notifyCount)
			return notifyCount
		case <-timer.C:
			etas, err := svc.GetETA(ctx, defaultCity, routeID, direction)
			if err != nil {
				now := time.Now().Format("15:04:05")
				log.Printf("%s  ⚠ 取得 ETA 失敗: %v", now, err)
				timer.Reset(interval) // maintain previous interval on error
				continue
			}

			etaSec := findStopETA(etas, stop)
			interval = calcInterval(etaSec, threshold)
			shouldNotify := state.check(etaSec, threshold)

			if shouldNotify {
				notifyCount++
			}

			logETA(etaSec, shouldNotify)

			if shouldNotify && notifyCmd != "" {
				etaMin := etaSec / 60
				if err := sendNotification(notifyCmd, routeName, stop.Name, etaMin); err != nil {
					log.Printf("通知發送失敗: %v", err)
				}
			}

			timer.Reset(interval)
		}
	}
}

// formatDuration formats a duration as human-readable text.
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
