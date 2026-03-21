import { z } from "zod";

// Response identity details (not the request identityDetails shape necessarily)
export const IdentityDetailsSchema = z.object({
  email: z.array(z.email()).optional(),
  phone: z.array(z.string()).optional(),

  address1: z.string().optional(),
  address2: z.string().optional(),
  city: z.string().optional(),
  state: z.string().optional(),
  zipCode: z.string().optional(),

  raw: z.record(z.string(), z.unknown()).optional(),
});

export type IdentityDetails = z.infer<typeof IdentityDetailsSchema>;
