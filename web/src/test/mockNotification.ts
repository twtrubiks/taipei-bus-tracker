import { vi } from "vitest";

/**
 * Sets up a mock Notification API on window and mock ServiceWorker registration.
 * Returns the spy function that records showNotification calls.
 */
export function setupMockNotification(
  permission: string = "granted",
  requestResult: string = "granted",
): ReturnType<typeof vi.fn> {
  const spy = vi.fn();

  class MockNotification {
    static permission = permission;
    static requestPermission = vi.fn().mockResolvedValue(requestResult);
    constructor(...args: unknown[]) {
      (spy as (...a: unknown[]) => void)(...args);
    }
  }

  Object.defineProperty(window, "Notification", {
    writable: true,
    configurable: true,
    value: MockNotification,
  });

  // Mock ServiceWorker registration with showNotification
  const mockRegistration = { showNotification: spy };
  Object.defineProperty(navigator, "serviceWorker", {
    writable: true,
    configurable: true,
    value: { getRegistration: vi.fn().mockResolvedValue(mockRegistration) },
  });

  return spy;
}
