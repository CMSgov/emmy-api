import { z } from "zod";

export const ClientDataSchema = z.object({
  // The response docs can refine these (max lengths, etc). For now keep permissive.
  clientName: z.string().optional(),
  clientCode: z.string().optional(),
  secondaryClient: z.string().optional(),
  endClient: z.string().optional(),
});

export type ClientData = z.infer<typeof ClientDataSchema>;
