import { z } from "zod";
import { FlexibleDateCodec, SsnSchema } from "./helpers";

export const StudentInfoProvidedSchema = z.object({
  // Flags indicating which student info fields were provided
  ssn: SsnSchema.optional(),
  dateOfBirth: FlexibleDateCodec.optional(),
  firstName: z.string().optional(),
  middleName: z.string().optional(),
  lastName: z.string().optional(),
  studentId: z.string().optional(),
  schoolName: z.string().optional(),
  schoolCode: z.string().optional(),
});

export type StudentInfoProvided = z.infer<typeof StudentInfoProvidedSchema>;
