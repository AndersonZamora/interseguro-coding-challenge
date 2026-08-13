import { Matrix, StatsResponse } from '../types';

const DIAGONAL_EPSILON = 1e-9;

/**
 * Determina si una matriz es diagonal: cuadrada, y todo elemento fuera de la
 * diagonal principal es ~0 dentro de una tolerancia. Una matriz no cuadrada
 * no puede ser diagonal por definicion.
 */
export function isDiagonal(matrix: Matrix, epsilon = DIAGONAL_EPSILON): boolean {
  const rows = matrix.length;
  const cols = matrix[0]?.length ?? 0;
  if (rows !== cols) {
    return false;
  }
  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      if (i !== j && Math.abs(matrix[i][j]) > epsilon) {
        return false;
      }
    }
  }
  return true;
}

/**
 * Calcula estadisticas agregadas sobre una lista de matrices: valor maximo,
 * minimo, promedio y suma de todos los valores combinados, y si alguna de
 * las matrices individuales es diagonal.
 */
export function computeStats(matrices: Matrix[]): StatsResponse {
  let max = -Infinity;
  let min = Infinity;
  let sum = 0;
  let count = 0;

  for (const matrix of matrices) {
    for (const row of matrix) {
      for (const value of row) {
        if (value > max) max = value;
        if (value < min) min = value;
        sum += value;
        count++;
      }
    }
  }

  const average = count === 0 ? 0 : sum / count;
  const anyDiagonal = matrices.some((matrix) => isDiagonal(matrix));

  return { max, min, average, sum, anyDiagonal };
}
