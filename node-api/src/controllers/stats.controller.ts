import { Request, Response } from 'express';
import { computeStats } from '../services/stats.service';
import { StatsPayload } from '../types';

export function postStats(req: Request, res: Response): void {
  const { q, r } = req.body as StatsPayload;
  const stats = computeStats([q, r]);
  res.status(200).json(stats);
}
