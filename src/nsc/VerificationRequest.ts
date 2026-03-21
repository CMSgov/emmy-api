import { z } from "zod";

// -----------------------------
// helpers
// -----------------------------

// ISO date (YYYY-MM-DD)
const IsoDateString = z.iso.date();

// ISO date codec
// - decode: YYYY-MM-DD -> Date
// - encode: Date -> YYYY-MM-DD
const IsoDateToDate = z.codec(IsoDateString, z.date(), {
  decode: (isoDateString) => new Date(isoDateString),
  encode: (date) => date.toISOString().slice(0, 10),
});

// Request fields: accept Date | YYYY-MM-DD, but when encoding for the wire they become YYYY-MM-DD.
const IsoDateRequestField = z.union([IsoDateString, z.date()]).pipe(IsoDateToDate);

// SSN: keep only digits (9)
const SsnSchema = z.string().regex(/^\d{9}$/, "Expected 9 digits");

// -----------------------------
// main schema
// -----------------------------

export const VerificationServicesRequestSchema = z.object({
  accountId: z.string(),

  organizationName: z.string().max(50).optional(),
  caseReferenceId: z.string().max(40).optional(),
  correlationId: z.string().max(40).optional(),

  contactEmail: z.email().max(256).optional(),

  ssn: SsnSchema.optional(),
  dateOfBirth: IsoDateRequestField.optional(),

  lastName: z.string().max(30),
  firstName: z.string().max(30),
  middleName: z.string().max(30).optional(),

  previousNames: z
    .array(
      z.object({
        firstName: z.string().max(30),
        lastName: z.string().max(30),
        middleName: z.string().max(30).optional(),
      }),
    )
    .max(5)
    .optional(),

  studentId: z.string().optional(),

  schoolCode: z.string().regex(/^\d{6}$/).optional(),
  schoolName: z.string().optional(),

  degreeTitle: z.string().optional(),
  yearAwarded: z.union([z.string().regex(/^\d{4}$/), z.number().int()]).optional(),
  major: z.string().optional(),

  degreeLevelCode: z.enum(["P", "M", "B", "A"]).optional(),

  applyLikeSchoolMatching: z.enum(["Y", "N"]).optional(),

  asOfDate: IsoDateRequestField.optional(),

  terms: z.literal("Y"),

  identityDetails: z
    .object({
      email: z.array(z.email().max(50)).optional(),
      phone: z.array(z.string().regex(/^\d{10}$/)).optional(),
      address1: z.string().optional(),
      address2: z.string().optional(),
      city: z.string().optional(),
      state: z.string().regex(/^[A-Z]{2}$/).optional(),
      zipCode: z.string().regex(/^\d{5}$/).optional(),
    })
    .optional(),

  endClient: z.string().max(30),
  secondaryClient: z.string().max(50).optional(),
  naicsCode: z.union([z.string().regex(/^\d{2,6}$/), z.number().int()]).optional(),

  startDate: IsoDateRequestField.optional(),
  jobTitle: z.string().max(30).optional(),
});

// For building requests: accept Dates (input type)
export type VerificationServicesRequest = z.input<typeof VerificationServicesRequestSchema>;

// For sending over the wire: encoded form (output type)
export type VerificationServicesRequestEncoded = z.output<
  typeof VerificationServicesRequestSchema
>;
