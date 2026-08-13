import { computeStats, isDiagonal } from '../../src/services/stats.service';

describe('isDiagonal', () => {
  it('detecta una matriz diagonal real', () => {
    expect(isDiagonal([[2, 0], [0, 3]])).toBe(true);
  });

  it('rechaza una matriz cuadrada no diagonal', () => {
    expect(isDiagonal([[2, 1], [0, 3]])).toBe(false);
  });

  it('rechaza una matriz rectangular (el concepto no aplica)', () => {
    expect(isDiagonal([[1, 0, 0], [0, 1, 0]])).toBe(false);
  });

  it('trata un solo elemento como diagonal', () => {
    expect(isDiagonal([[5]])).toBe(true);
  });

  it('tolera ruido numerico por debajo del epsilon', () => {
    expect(isDiagonal([[2, 1e-12], [1e-12, 3]])).toBe(true);
  });
});

describe('computeStats', () => {
  it('calcula max, min, promedio y suma sobre una sola matriz', () => {
    const stats = computeStats([[[1, 2], [3, 4]]]);
    expect(stats.max).toBe(4);
    expect(stats.min).toBe(1);
    expect(stats.sum).toBe(10);
    expect(stats.average).toBe(2.5);
  });

  it('combina estadisticas sobre multiples matrices de distinto tamano', () => {
    const q = [[1, -2], [3, 4], [5, 6]];
    const r = [[10, 0], [0, 20]];
    const stats = computeStats([q, r]);
    expect(stats.max).toBe(20);
    expect(stats.min).toBe(-2);
    expect(stats.sum).toBe(1 - 2 + 3 + 4 + 5 + 6 + 10 + 0 + 0 + 20);
  });

  it('maneja valores negativos mezclados con positivos', () => {
    const stats = computeStats([[[-5, -1], [3, 7]]]);
    expect(stats.max).toBe(7);
    expect(stats.min).toBe(-5);
  });

  it('un solo elemento produce max=min=avg=sum y anyDiagonal=true', () => {
    const stats = computeStats([[[5]]]);
    expect(stats.max).toBe(5);
    expect(stats.min).toBe(5);
    expect(stats.average).toBe(5);
    expect(stats.sum).toBe(5);
    expect(stats.anyDiagonal).toBe(true);
  });

  it('anyDiagonal es true si al menos una matriz del conjunto es diagonal', () => {
    const notDiagonal = [[1, 2], [3, 4]];
    const diagonal = [[9, 0], [0, 9]];
    const stats = computeStats([notDiagonal, diagonal]);
    expect(stats.anyDiagonal).toBe(true);
  });

  it('anyDiagonal es false si ninguna matriz del conjunto es diagonal', () => {
    const a = [[1, 2], [3, 4]];
    const b = [[1, 2, 3], [4, 5, 6]];
    const stats = computeStats([a, b]);
    expect(stats.anyDiagonal).toBe(false);
  });
});
