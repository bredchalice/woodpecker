import { describe, expect, it } from 'vitest';

import { classifyLogLine } from './logSeverity';

describe('classifyLogLine', () => {
  it('does not treat ordinary build output as an error or warning', () => {
    const lines = [
      '#18 93.43 ../public/assets/index-CM6Za5o-.css 23.17 kB | gzip: 5.46 kB',
      '#18 DONE 93.7s',
      '#17 [linux/arm64 stage-1 3/7] RUN apt-get update && apt-get install -y ca-certificates',
      'The failed requests are retried automatically',
      'failure policy: continue',
      'deprecated package metadata is retained for compatibility',
      'https://example.invalid/error/failed-warning',
      '{"error":"none","warning":"none"}',
      'copy error-report.txt /tmp/error-report.txt',
    ];

    for (const line of lines) expect(classifyLogLine(line)).toBeNull();
  });

  it('does not classify zero-problem summaries', () => {
    expect(classifyLogLine('0 errors, 0 warnings')).toBeNull();
    expect(classifyLogLine('Found no errors')).toBeNull();
    expect(classifyLogLine('0 failures')).toBeNull();
  });

  it('recognizes explicit errors', () => {
    expect(classifyLogLine('ERROR: build failed')).toBe('error');
    expect(classifyLogLine('#18 ERROR: process "/bin/sh -c yarn build" did not complete successfully: exit code: 1')).toBe(
      'error',
    );
    expect(classifyLogLine('src/index.ts(12,3): error TS2322: Type string is not assignable')).toBe('error');
    expect(classifyLogLine('npm ERR! lifecycle script failed')).toBe('error');
    expect(classifyLogLine('Command exited with code 2')).toBe('error');
    expect(classifyLogLine('FAIL src/example.test.ts')).toBe('error');
  });

  it('recognizes explicit warnings', () => {
    expect(classifyLogLine('WARNING: cache is disabled')).toBe('warning');
    expect(classifyLogLine('npm WARN deprecated package@1.0.0')).toBe('warning');
    expect(classifyLogLine('DeprecationWarning: old API')).toBe('warning');
    expect(classifyLogLine('src/main.go:12:3: warning: unused value')).toBe('warning');
  });

  it('ignores command echo lines even if the command contains severity words', () => {
    expect(classifyLogLine('+ echo ERROR: this is test data')).toBeNull();
    expect(classifyLogLine('+ grep warning build.log')).toBeNull();
  });
});
