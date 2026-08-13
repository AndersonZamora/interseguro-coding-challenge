import { Router, Request, Response } from 'express';
import jwt from 'jsonwebtoken';
import { env } from '../config/env';
import { tokenRequestSchema, TokenResponse } from '../types';

const router = Router();

// Token de demo de 1 hora: no hay modelo de usuarios en el enunciado, asi
// que se autentica contra un client_id/client_secret fijo por env var. Este
// es el unico emisor de JWT del sistema; la API de Go firma sus propios
// tokens de servicio con el mismo secreto para su llamada interna, sin pasar
// por este endpoint (ver go-api/internal/auth/jwt.go).
const TOKEN_TTL_SECONDS = 60 * 60;

router.post('/token', (req: Request, res: Response) => {
  const parsed = tokenRequestSchema.safeParse(req.body);
  if (!parsed.success) {
    res.status(400).json({ error: 'invalid request payload', details: parsed.error.errors });
    return;
  }

  const { client_id, client_secret } = parsed.data;
  if (client_id !== env.demoClientId || client_secret !== env.demoClientSecret) {
    res.status(401).json({ error: 'invalid client credentials' });
    return;
  }

  const access_token = jwt.sign({ sub: client_id, iss: 'interseguro-challenge' }, env.jwtSecret, {
    expiresIn: TOKEN_TTL_SECONDS,
  });

  const response: TokenResponse = { access_token, expires_in: TOKEN_TTL_SECONDS };
  res.status(200).json(response);
});

export default router;
