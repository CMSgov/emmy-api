import { z } from 'zod'

const ssnRaw = z.string().trim();

// allow: 123-45-6789 | 123 45 6789 | 123456789 | 6789
const ssnAccepted = ssnRaw.regex(
  /^(?:\d{3}[- ]\d{2}[- ]\d{4}|\d{9}|\d{4})$/,
  "Invalid ssn format (expected XXX-XX-XXXX, XXX XX XXXX, XXXXXXXXX, or last4)"
);

// optional normalization to digits only
export const SsnSchema = ssnAccepted.transform((s) => s.replace(/[^\d]/g, ""));

const YyyyMmDdBasic = z
  .string()
  .regex(/^\d{8}$/, { error: "Expected date in YYYYMMDD format" });

const YyyyMmDdSlash = z
  .string()
  .regex(/^\d{4}\/\d{2}\/\d{2}$/, { error: "Expected date in YYYY/MM/DD format" });

const FlexibleDateInput = z.union([z.iso.date(), YyyyMmDdBasic, YyyyMmDdSlash]);

const parseYmdToDateUtc = (
  yyyy: string,
  mm: string,
  dd: string,
  ctx: z.core.ParseContext,
  rawInput: string
): Date | typeof z.NEVER => {
  const y = Number(yyyy);
  const m = Number(mm);
  const d = Number(dd);

  if (!Number.isInteger(y) || !Number.isInteger(m) || !Number.isInteger(d)) {
    ctx.issues.push({
      code: "invalid_format",
      format: "date",
      input: rawInput,
      message: "Invalid date components",
    });
    return z.NEVER;
  }

  // Construct in UTC to avoid timezone-induced day drift
  const date = new Date(Date.UTC(y, m - 1, d));

  // Reject normalization (e.g. 2024-02-31 turning into March)
  if (
    date.getUTCFullYear() !== y ||
    date.getUTCMonth() !== m - 1 ||
    date.getUTCDate() !== d
  ) {
    ctx.issues.push({
      code: "invalid_format",
      format: "date",
      input: rawInput,
      message: "Invalid calendar date",
    });
    return z.NEVER;
  }

  return date;
};

const decodeFlexibleDate = (s: string, ctx: z.core.ParseContext): Date | typeof z.NEVER => {
  // At this point, `s` is guaranteed by FlexibleDateInput to match one of:
  // - YYYY-MM-DD
  // - YYYYMMDD
  // - YYYY/MM/DD

  const dash = /^(\d{4})-(\d{2})-(\d{2})$/.exec(s);
  if (dash) return parseYmdToDateUtc(dash[1], dash[2], dash[3], ctx, s);

  const basic = /^(\d{4})(\d{2})(\d{2})$/.exec(s);
  if (basic) return parseYmdToDateUtc(basic[1], basic[2], basic[3], ctx, s);

  const slash = /^(\d{4})\/(\d{2})\/(\d{2})$/.exec(s);
  if (slash) return parseYmdToDateUtc(slash[1], slash[2], slash[3], ctx, s);

  // Should be unreachable because FlexibleDateInput validates formats first,
  // but keep a defensive error in case the input schema changes.
  ctx.issues.push({
    code: "invalid_format",
    format: "date",
    input: s,
    message: "Expected YYYYMMDD, YYYY-MM-DD, or YYYY/MM/DD",
  });
  return z.NEVER;
};

const encodeIsoDateOnlyUtc = (d: Date, ctx: z.core.ParseContext): string | typeof z.NEVER => {
  if (!(d instanceof Date) || Number.isNaN(d.getTime())) {
    ctx.issues.push({
      code: "invalid_type",
      expected: "date",
      input: d as any,
      message: "Expected a valid Date",
    });
    return z.NEVER;
  }

  // Canonical ISO date-only string
  return d.toISOString().slice(0, 10);
};

/**
 * FlexibleDateCodec
 * - decode: accepts "YYYYMMDD" | "YYYY-MM-DD" | "YYYY/MM/DD" and produces Date
 * - encode: always produces canonical ISO date-only "YYYY-MM-DD"
 */
export const FlexibleDateCodec = z.codec(FlexibleDateInput, z.date(), {
  decode: decodeFlexibleDate,
  encode: encodeIsoDateOnlyUtc,
});

export type FlexibleDate = z.infer<typeof FlexibleDateCodec>;
