import { Matrix } from '../types';
import MatrixTable from './MatrixTable';

interface Props {
  q: Matrix;
  r: Matrix;
}

export default function QRResult({ q, r }: Props) {
  return (
    <div className="flex flex-wrap gap-8">
      <div>
        <h3 className="mb-2 font-medium text-slate-700">Q (ortogonal)</h3>
        <MatrixTable matrix={q} />
      </div>
      <div>
        <h3 className="mb-2 font-medium text-slate-700">R (triangular superior)</h3>
        <MatrixTable matrix={r} />
      </div>
    </div>
  );
}
