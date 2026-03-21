import { z } from 'zod';

const EnvSchema = z.object({
  PORT: z.coerce.number().int().positive().default(3000),

  NSC_BASE_URL: z.url(),
  NSC_ACCESS_TOKEN_URL: z.url(),
  NSC_ACCOUNT_ID: z.string(),
  NSC_CLIENT_ID: z.string(),
  NSC_CLIENT_SECRET: z.string(),
})

export type Env = z.infer<typeof EnvSchema>;

export const ENV: Env = EnvSchema.parse(process.env);
