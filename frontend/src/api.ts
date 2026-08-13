import axios, { AxiosError } from 'axios';
import { ApiErrorBody, QrResponse, TokenResponse } from './types';

const GO_API_URL = import.meta.env.VITE_GO_API_URL ?? 'http://localhost:8080';
const NODE_API_URL = import.meta.env.VITE_NODE_API_URL ?? 'http://localhost:3000';

function extractErrorMessage(err: unknown, fallback: string): string {
  const axiosErr = err as AxiosError<ApiErrorBody>;
  return axiosErr.response?.data?.error ?? fallback;
}

export async function fetchToken(clientId: string, clientSecret: string): Promise<TokenResponse> {
  try {
    const res = await axios.post<TokenResponse>(`${NODE_API_URL}/auth/token`, {
      client_id: clientId,
      client_secret: clientSecret,
    });
    return res.data;
  } catch (err) {
    throw new Error(extractErrorMessage(err, 'error obteniendo token'));
  }
}

export async function fetchQr(matrix: number[][], token: string): Promise<QrResponse> {
  try {
    const res = await axios.post<QrResponse>(
      `${GO_API_URL}/api/v1/qr`,
      { matrix },
      { headers: { Authorization: `Bearer ${token}` } },
    );
    return res.data;
  } catch (err) {
    throw new Error(extractErrorMessage(err, 'error calculando la factorización QR'));
  }
}
