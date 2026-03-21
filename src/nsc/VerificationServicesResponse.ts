import { z } from "zod";

import { ClientDataSchema } from "./ClientData";
import { StudentInfoProvidedSchema } from "./StudentInfoProvided";
import { TransactionDetailsSchema } from "./TransactionDetails";
import { DegreeDetailsSchema } from "./DegreeDetails";
import { EnrollmentDetailsSchema } from "./EnrollmentDetails";
import { DerivedDetailsSchema } from "./DerivedDetails";
import { IdentityDetailsSchema } from "./IdentityDetails";

// Top-level response wrapper (placeholder until we align with the official response shape)
export const VerificationServicesResponseSchema = z.object({
  transactionDetails: TransactionDetailsSchema.optional(),

  clientData: ClientDataSchema.optional(),
  studentInfoProvided: StudentInfoProvidedSchema.optional(),

  degreeDetails: z.array(DegreeDetailsSchema).optional(),
  enrollmentDetails: z.array(EnrollmentDetailsSchema).optional(),

  derivedDetails: DerivedDetailsSchema.optional(),
  identityDetails: IdentityDetailsSchema.optional(),

  // Common envelope-ish fields
  correlationId: z.string().optional(),
  caseReferenceId: z.string().optional(),

  // Raw payload passthrough (keep unknown until modeled)
  data: z.unknown().optional(),
});

export type VerificationServicesResponse = z.infer<
  typeof VerificationServicesResponseSchema
>;
