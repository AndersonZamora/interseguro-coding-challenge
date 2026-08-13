// Variables de entorno minimas para que los modulos de config puedan cargar
// durante los tests, sin depender de un .env real.
process.env.JWT_SECRET = process.env.JWT_SECRET ?? 'test-secret';
process.env.DEMO_CLIENT_ID = process.env.DEMO_CLIENT_ID ?? 'test-client';
process.env.DEMO_CLIENT_SECRET = process.env.DEMO_CLIENT_SECRET ?? 'test-secret-value';
process.env.ALLOWED_ORIGIN = process.env.ALLOWED_ORIGIN ?? '*';
