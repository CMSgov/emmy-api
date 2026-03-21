import { z } from "zod";

// Helpers can live in-module
const IsoDateString = z.iso.date();

export const TransactionDetailsSchema = z.object({
  transactionId: z.string().optional(),
  correlationId: z.string().optional(),
  caseReferenceId: z.string().optional(),

  // Common timestamps
  createdAt: z.union([IsoDateString, z.iso.datetime()]).optional(),
  submittedAt: z.union([IsoDateString, z.iso.datetime()]).optional(),

  status: z.string().optional(),
  message: z.string().optional(),
});

export type TransactionDetails = z.infer<typeof TransactionDetailsSchema>;
