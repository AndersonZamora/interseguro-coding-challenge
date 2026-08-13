import { useState } from 'react';
import { fetchQr } from '../api';
import { QrResponse } from '../types';

const DEFAULT_MATRIX = '[[12, -51, 4], [6, 167, -68], [-4, 24, -41]]';

interface Props {
  token: string | null;
  onResult: (result: QrResponse) => void;
}

function parseMatrix(text: string): number[][] {
  const value = JSON.parse(text);
  if (!Array.isArray(value) || value.length === 0 || !value.every((row) => Array.isArray(row))) {
    throw new Error('el texto debe ser un array de arrays de números, ej. [[1,2],[3,4]]');
  }
  return value as number[][];
}

export default function MatrixInput({ token, onResult }: Props) {
  const [text, setText] = useState(DEFAULT_MATRIX);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit() {
    setError(null);
    let matrix: number[][];
    try {
      matrix = parseMatrix(text);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'JSON inválido');
      return;
    }

    if (!token) {
      setError('primero obtén un token de autenticación');
      return;
    }

    setLoading(true);
    try {
      const result = await fetchQr(matrix, token);
      onResult(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'error desconocido');
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="rounded-lg border border-slate-200 p-4">
      <h2 className="mb-3 text-lg font-semibold text-slate-800">2. Matriz de entrada</h2>
      <textarea
        className="w-full rounded border border-slate-300 p-2 font-mono text-sm"
        rows={4}
        value={text}
        onChange={(e) => setText(e.target.value)}
      />
      <button
        className="mt-3 rounded bg-blue-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
        onClick={handleSubmit}
        disabled={loading}
      >
        {loading ? 'Calculando…' : 'Calcular QR y estadísticas'}
      </button>
      {error && <p className="mt-2 text-sm text-red-600">{error}</p>}
    </section>
  );
}
