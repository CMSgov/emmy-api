import { z } from "zod";

// "Derived" fields are usually enrichment / computed flags.
export const DerivedDetailsSchema = z.object({
  isMatch: z.boolean().optional(),
  matchConfidence: z.number().min(0).max(1).optional(),

  // Common derived identifiers
  personId: z.string().optional(),
  studentId: z.string().optional(),

  raw: z.record(z.string(), z.unknown()).optional(),
});

export type DerivedDetails = z.infer<typeof DerivedDetailsSchema>;
