export function statusColor(eta: number): string {
  if (eta >= 0 && eta <= 180) return "text-green-600 font-bold";
  if (eta > 180) return "text-blue-600";
  return "text-gray-400";
}
