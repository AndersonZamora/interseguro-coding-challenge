import { NextFunction, Request, Response } from 'express';
import { AnyZodObject } from 'zod';

/**
 * Valida req.body contra un schema de zod y lo reemplaza por la version
 * parseada (con tipos y defaults aplicados). Los errores de validacion se
 * delegan al errorHandler central via next(err).
 */
export function validateBody(schema: AnyZodObject) {
  return (req: Request, _res: Response, next: NextFunction): void => {
    try {
      req.body = schema.parse(req.body);
      next();
    } catch (err) {
      next(err);
    }
  };
}
