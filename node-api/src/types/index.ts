import { z } from 'zod';

export const matrixSchema = z.array(z.array(z.number())).min(1);
export type Matrix = z.infer<typeof matrixSchema>;

export const statsPayloadSchema = z.object({
  q: matrixSchema,
  r: matrixSchema,
});
export type StatsPayload = z.infer<typeof statsPayloadSchema>;

export interface StatsResponse {
  max: number;
  min: number;
  average: number;
  sum: number;
  anyDiagonal: boolean;
}

export const tokenRequestSchema = z.object({
  client_id: z.string().min(1),
  client_secret: z.string().min(1),
});
export type TokenRequest = z.infer<typeof tokenRequestSchema>;

export interface TokenResponse {
  access_token: string;
  expires_in: number;
}
