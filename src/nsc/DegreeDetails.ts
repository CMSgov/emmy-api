import { z } from "zod";

export const DegreeDetailsSchema = z.object({
  degreeTitle: z.string().optional(),
  degreeLevel: z.string().optional(),
  degreeLevelCode: z.string().optional(),

  major: z.string().optional(),
  majors: z.array(z.string()).optional(),
  minors: z.array(z.string()).optional(),

  yearAwarded: z.union([z.string().regex(/^\d{4}$/), z.number().int()]).optional(),
  awardDate: z.iso.date().optional(),

  // Anything else the API returns that we haven't modeled yet
  raw: z.record(z.string(), z.unknown()).optional(),
});

export type DegreeDetails = z.infer<typeof DegreeDetailsSchema>;
