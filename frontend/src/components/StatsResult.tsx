import { StatsResponse } from '../types';

interface Props {
  stats: StatsResponse;
}

function Stat({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded border border-slate-200 px-3 py-2">
      <div className="text-xs uppercase tracking-wide text-slate-500">{label}</div>
      <div className="text-lg font-semibold text-slate-800">{value}</div>
    </div>
  );
}

export default function StatsResult({ stats }: Props) {
  return (
    <div className="grid grid-cols-2 gap-3 sm:grid-cols-5">
      <Stat label="Máximo" value={stats.max.toFixed(4)} />
      <Stat label="Mínimo" value={stats.min.toFixed(4)} />
      <Stat label="Promedio" value={stats.average.toFixed(4)} />
      <Stat label="Suma total" value={stats.sum.toFixed(4)} />
      <Stat label="¿Alguna diagonal?" value={stats.anyDiagonal ? 'Sí' : 'No'} />
    </div>
  );
}
