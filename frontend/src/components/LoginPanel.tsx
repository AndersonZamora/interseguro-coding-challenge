import { useState } from 'react';
import { fetchToken } from '../api';

interface Props {
  token: string | null;
  onToken: (token: string) => void;
}

export default function LoginPanel({ token, onToken }: Props) {
  const [clientId, setClientId] = useState('interseguro-demo');
  const [clientSecret, setClientSecret] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleLogin() {
    setLoading(true);
    setError(null);
    try {
      const res = await fetchToken(clientId, clientSecret);
      onToken(res.access_token);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'error desconocido');
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="rounded-lg border border-slate-200 p-4">
      <h2 className="mb-3 text-lg font-semibold text-slate-800">1. Autenticación (JWT demo)</h2>
      <div className="flex flex-wrap items-end gap-3">
        <label className="flex flex-col text-sm text-slate-600">
          Client ID
          <input
            className="mt-1 rounded border border-slate-300 px-2 py-1"
            value={clientId}
            onChange={(e) => setClientId(e.target.value)}
          />
        </label>
        <label className="flex flex-col text-sm text-slate-600">
          Client Secret
          <input
            type="password"
            className="mt-1 rounded border border-slate-300 px-2 py-1"
            value={clientSecret}
            onChange={(e) => setClientSecret(e.target.value)}
            placeholder="DEMO_CLIENT_SECRET configurado en node-api"
          />
        </label>
        <button
          className="rounded bg-blue-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
          onClick={handleLogin}
          disabled={loading || clientSecret === ''}
        >
          {loading ? 'Obteniendo…' : 'Obtener token'}
        </button>
        {token && <span className="text-sm text-green-700">✓ Token obtenido</span>}
      </div>
      {error && <p className="mt-2 text-sm text-red-600">{error}</p>}
    </section>
  );
}
