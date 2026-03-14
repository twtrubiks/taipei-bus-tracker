/** Normalize stop name for cross-provider matching (trim, remove spaces, NFKC for full/half-width). */
export function normalizeName(name: string): string {
  return name.replace(/\s/g, "").normalize("NFKC");
}
