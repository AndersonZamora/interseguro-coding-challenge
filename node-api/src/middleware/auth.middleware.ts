import { NextFunction, Request, Response } from 'express';
import jwt from 'jsonwebtoken';
import { env } from '../config/env';

// Mismo mensaje uniforme que usa la API de Go para 401, asi el frontend
// maneja un solo caso sin importar cual de las dos APIs lo devuelve.
const UNAUTHORIZED_MESSAGE = 'missing or invalid authorization token';

export function verifyJwt(req: Request, res: Response, next: NextFunction): void {
  const header = req.header('Authorization');
  if (!header?.startsWith('Bearer ')) {
    res.status(401).json({ error: UNAUTHORIZED_MESSAGE });
    return;
  }

  const token = header.slice('Bearer '.length).trim();
  try {
    jwt.verify(token, env.jwtSecret);
    next();
  } catch {
    res.status(401).json({ error: UNAUTHORIZED_MESSAGE });
  }
}
