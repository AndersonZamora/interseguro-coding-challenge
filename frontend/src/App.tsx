import { useState } from 'react';
import LoginPanel from './components/LoginPanel';
import MatrixInput from './components/MatrixInput';
import QRResult from './components/QRResult';
import StatsResult from './components/StatsResult';
import { QrResponse } from './types';

export default function App() {
  const [token, setToken] = useState<string | null>(null);
  const [result, setResult] = useState<QrResponse | null>(null);

  return (
    <div className="mx-auto max-w-4xl space-y-6 p-6">
      <header>
        <h1 className="text-2xl font-bold text-slate-900">Interseguro — Reto técnico</h1>
        <p className="text-sm text-slate-600">
          Factorización QR (API Go) + estadísticas (API Node) sobre una matriz de entrada.
        </p>
      </header>

      <LoginPanel token={token} onToken={setToken} />
      <MatrixInput token={token} onResult={setResult} />

      {result && (
        <section className="space-y-4 rounded-lg border border-slate-200 p-4">
          <h2 className="text-lg font-semibold text-slate-800">3. Resultado</h2>
          <QRResult q={result.q} r={result.r} />
          <StatsResult stats={result.stats} />
        </section>
      )}
    </div>
  );
}
