interface Env {
  port: number;
  jwtSecret: string;
  demoClientId: string;
  demoClientSecret: string;
  allowedOrigin: string;
}

function required(name: string): string {
  const value = process.env[name];
  if (!value) {
    throw new Error(`${name} es obligatorio`);
  }
  return value;
}

// env se resuelve al importar el modulo: si falta una variable obligatoria,
// el proceso falla al arrancar en vez de quedar en un estado inconsistente
// (mismo criterio que config.Load() en la API de Go).
export const env: Env = {
  port: Number(process.env.PORT ?? 3000),
  jwtSecret: required('JWT_SECRET'),
  demoClientId: process.env.DEMO_CLIENT_ID ?? 'interseguro-demo',
  demoClientSecret: required('DEMO_CLIENT_SECRET'),
  allowedOrigin: process.env.ALLOWED_ORIGIN ?? '*',
};
