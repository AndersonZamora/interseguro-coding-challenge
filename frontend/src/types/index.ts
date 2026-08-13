// Tipos de dominio duplicados (a proposito, sin paquete compartido) desde
// node-api/src/types — ver decision documentada en docs/decisiones-arquitectura.md.
export type Matrix = number[][];

export interface StatsResponse {
  max: number;
  min: number;
  average: number;
  sum: number;
  anyDiagonal: boolean;
}

export interface QrResponse {
  q: Matrix;
  r: Matrix;
  stats: StatsResponse;
}

export interface TokenResponse {
  access_token: string;
  expires_in: number;
}

export interface ApiErrorBody {
  error: string;
}
