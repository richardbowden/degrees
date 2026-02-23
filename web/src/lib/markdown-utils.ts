/**
 * Extract the first N bullet-point items from a Markdown string.
 * Returns plain text strings (no Markdown formatting).
 */
export function extractFeaturePreview(markdown: string, limit = 4): string[] {
  const matches = markdown.match(/^[-*]\s+(.+)$/gm);
  if (!matches) return [];
  return matches.slice(0, limit).map(m => m.replace(/^[-*]\s+/, '').replace(/[*_`]/g, ''));
}

/**
 * Extract ## and ### heading texts from Markdown.
 */
export function extractSections(markdown: string): string[] {
  const matches = markdown.match(/^#{2,3}\s+(.+)$/gm);
  if (!matches) return [];
  return matches.map(m => m.replace(/^#{2,3}\s+/, ''));
}
