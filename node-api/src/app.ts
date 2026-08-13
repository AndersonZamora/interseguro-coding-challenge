import express, { Express } from 'express';
import helmet from 'helmet';
import cors from 'cors';
import morgan from 'morgan';
import { env } from './config/env';
import authRoutes from './routes/auth.routes';
import statsRoutes from './routes/stats.routes';
import { errorHandler } from './middleware/error.middleware';

/**
 * Ensambla la aplicacion Express con todos los middlewares y rutas
 * registrados. Exportada por separado de index.ts para poder testearla con
 * supertest sin levantar un puerto real.
 */
export function createApp(): Express {
  const app = express();

  app.use(helmet());
  app.use(cors({ origin: env.allowedOrigin }));
  app.use(express.json({ limit: '1mb' }));
  app.use(morgan('combined'));

  app.get('/health', (_req, res) => res.json({ status: 'ok' }));
  app.use('/auth', authRoutes);
  app.use('/api/v1', statsRoutes);

  app.use(errorHandler);

  return app;
}
