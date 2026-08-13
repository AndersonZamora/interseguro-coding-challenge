// Cache simple del access_token en localStorage para que sobreviva a un
// refresh de pagina. No es parte de lo pedido por el reto (que solo exige
// JWT protegiendo las consultas), es una mejora de UX.
const STORAGE_KEY = 'interseguro_access_token';

interface DecodedPayload {
  exp?: number;
}

function decodePayload(token: string): DecodedPayload | null {
  try {
    const [, payload] = token.split('.');
    return JSON.parse(atob(payload));
  } catch {
    return null;
  }
}

function isExpired(token: string): boolean {
  const payload = decodePayload(token);
  if (!payload?.exp) return true;
  return Date.now() >= payload.exp * 1000;
}

/** Devuelve el token cacheado si existe y no esta expirado; si esta vencido, lo limpia. */
export function loadCachedToken(): string | null {
  const token = localStorage.getItem(STORAGE_KEY);
  if (!token || isExpired(token)) {
    localStorage.removeItem(STORAGE_KEY);
    return null;
  }
  return token;
}

export function cacheToken(token: string): void {
  localStorage.setItem(STORAGE_KEY, token);
}

export function clearCachedToken(): void {
  localStorage.removeItem(STORAGE_KEY);
}
