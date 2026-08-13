import { Router } from 'express';
import { postStats } from '../controllers/stats.controller';
import { verifyJwt } from '../middleware/auth.middleware';
import { validateBody } from '../middleware/validate.middleware';
import { statsPayloadSchema } from '../types';

const router = Router();

router.post('/stats', verifyJwt, validateBody(statsPayloadSchema), postStats);

export default router;
