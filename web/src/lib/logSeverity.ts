export type LogSeverity = 'error' | 'warning' | null;

const ansiEscapeRegex = /\x1B\[[0-?]*[ -/]*[@-~]/g;
const buildKitPrefixRegex = /^#\d+\s+(?:\d+(?:\.\d+)?\s+)?/;
const zeroProblemRegex = /\b(?:0\s+(?:errors?|warnings?|failures?)|no\s+(?:errors?|warnings?|failures?))\b/i;

const errorPatterns = [
  /^(?:##\[error\]|\[(?:error|fatal)\])(?:\s|$)/i,
  /^(?:error|fatal|panic|exception)(?:\b|\s*[:[])/i,
  /^FAIL(?:ED)?(?:\s|:|$)/,
  /\berror\s+TS\d+:/i,
  /^(?:npm|yarn|pnpm)\s+(?:ERR!|error\b)/i,
  /^(?:[^\r\n]*?:\d+(?::\d+)?(?::|\s))?\s*(?:error|fatal)(?:\s+[A-Z][A-Z0-9_-]*|\s+TS\d+)?\s*:/i,
  /\b(?:exit(?:ed)?|returned)\s+(?:with\s+)?(?:code|status)\s+[1-9]\d*\b/i,
  /^(?:command|process)\b.*\b(?:failed|exited)\b.*\b(?:code|status)\s+[1-9]\d*\b/i,
  /^(?:found\s+)?[1-9]\d*\s+errors?(?:\s|$)/i,
  /^✖\s+.*\b[1-9]\d*\s+errors?\b/i,
];

const warningPatterns = [
  /^(?:##\[warning\]|\[(?:warning|warn)\])(?:\s|$)/i,
  /^(?:warning|warn)(?:\b|\s*:)/i,
  /^(?:npm|yarn|pnpm)\s+WARN(?:ING)?\b/i,
  /^(?:Deprecation|Experimental|Runtime|Syntax)Warning(?:\b|:)/,
  /^(?:[^\r\n]*?:\d+(?::\d+)?(?::|\s))?\s*warning(?:\s+[A-Z][A-Z0-9_-]*)?\s*:/i,
  /^(?:found\s+)?[1-9]\d*\s+warnings?(?:\s|$)/i,
  /^✖\s+.*\b[1-9]\d*\s+warnings?\b/i,
];

export function stripLogAnsi(text: string): string {
  return text.replace(ansiEscapeRegex, '');
}

function normalizeLogLine(rawText: string): string {
  const text = stripLogAnsi(rawText).trim();
  return text.replace(buildKitPrefixRegex, '').trim();
}

export function classifyLogLine(rawText: string): LogSeverity {
  const text = normalizeLogLine(rawText);
  if (!text || text.startsWith('+ ') || zeroProblemRegex.test(text)) return null;

  if (errorPatterns.some((pattern) => pattern.test(text))) return 'error';
  if (warningPatterns.some((pattern) => pattern.test(text))) return 'warning';
  return null;
}
