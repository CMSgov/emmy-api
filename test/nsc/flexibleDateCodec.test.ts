import fc from 'fast-check';
import { describe, expect, it } from 'vitest';
import { z } from 'zod';
import { FlexibleDateCodec } from '../../src/nsc/helpers';

describe('FlexibleDateCodec', () => {
  it.each([
    ['YYYYMMDD', '20240209', '2024-02-09'],
    ['YYYY-MM-DD', '2024-02-09', '2024-02-09'],
    ['YYYY/MM/DD', '2024/02/09', '2024-02-09'],
  ])('should decode %s into a Date', (_label, input, expectedIsoDate) => {
    // Act
    const result = FlexibleDateCodec.parse(input);

    // Assert
    expect(result).toBeInstanceOf(Date);
    expect(z.encode(FlexibleDateCodec, result)).toBe(expectedIsoDate);
  });

  it('should reject invalid calendar dates', () => {
    const invalid = ['20240230', '2024-02-30', '2024/02/30', '2024-13-01', '2024-00-01'];

    for (const input of invalid) {
      expect(() => FlexibleDateCodec.parse(input)).toThrow(z.ZodError);
    }
  });

  it('should reject unsupported formats', () => {
    const invalid = ['2024.02.09', '02/09/2024', '2024-2-9', '2024029', 'abcd', ''];

    for (const input of invalid) {
      expect(() => FlexibleDateCodec.parse(input)).toThrow(z.ZodError);
    }
  });

  it('should always encode a Date to YYYY-MM-DD', () => {
    // Arrange (UTC date to keep the expected day stable)
    const d = new Date('2024-02-09T00:00:00.000Z');

    // Act
    const out = z.encode(FlexibleDateCodec, d);

    // Assert
    expect(out).toBe('2024-02-09');
  });

  it('should decode YYYY-MM-DD, YYYY/MM/DD, YYYYMMDD consistently', () => {
    fc.assert(
      fc.property(validYmdArb, ({ yyyy, mm, dd }) => {
        const iso = `${yyyy}-${mm}-${dd}`;
        const slash = `${yyyy}/${mm}/${dd}`;
        const basic = `${yyyy}${mm}${dd}`;

        const a = z.decode(FlexibleDateCodec, iso);
        const b = z.decode(FlexibleDateCodec, slash);
        const c = z.decode(FlexibleDateCodec, basic);

        const expected = `${yyyy}-${mm}-${dd}`;
        expect(z.encode(FlexibleDateCodec, a)).toBe(expected);
        expect(z.encode(FlexibleDateCodec, b)).toBe(expected);
        expect(z.encode(FlexibleDateCodec, c)).toBe(expected);
      }),
    );
  });
});

// Helpers
const pad2 = (n: number): string => String(n).padStart(2, '0');

// Limit years to > 2000 to match domain expectations and avoid 0-prefixed years.
const dateAfter2000Arb = fc.date({
  noInvalidDate: true,
  min: new Date('2001-01-01T00:00:00.000Z'),
  max: new Date('9999-12-31T23:59:59.999Z'),
});

const validYmdArb = dateAfter2000Arb.map((d) => {
  const yyyy = String(d.getUTCFullYear());
  const mm = pad2(d.getUTCMonth() + 1);
  const dd = pad2(d.getUTCDate());
  return { yyyy, mm, dd };
});
