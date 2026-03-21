import { z } from "zod";

const IsoDateString = z.iso.date();

export const EnrollmentDetailsSchema = z.object({
  schoolName: z.string().optional(),
  schoolCode: z.string().optional(),

  enrollmentStatus: z.string().optional(),
  enrollmentStatusCode: z.string().optional(),

  startDate: IsoDateString.optional(),
  endDate: IsoDateString.optional(),
  asOfDate: IsoDateString.optional(),

  attendance: z.string().optional(),
  attendanceCode: z.string().optional(),

  raw: z.record(z.string(), z.unknown()).optional(),
});

export type EnrollmentDetails = z.infer<typeof EnrollmentDetailsSchema>;
