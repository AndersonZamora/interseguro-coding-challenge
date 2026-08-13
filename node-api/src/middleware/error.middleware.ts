import { NextFunction, Request, Response } from 'express';
import { ZodError } from 'zod';

// Firma de 4 parametros requerida por Express para reconocer esto como error
// handler; debe registrarse al final de la cadena de middlewares.
export function errorHandler(err: unknown, _req: Request, res: Response, _next: NextFunction): void {
  if (err instanceof ZodError) {
    res.status(400).json({ error: 'invalid request payload', details: err.errors });
    return;
  }

  console.error(err);
  res.status(500).json({ error: 'internal server error' });
}
