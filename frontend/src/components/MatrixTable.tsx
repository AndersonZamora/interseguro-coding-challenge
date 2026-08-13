import { Matrix } from '../types';

interface Props {
  matrix: Matrix;
}

export default function MatrixTable({ matrix }: Props) {
  return (
    <table className="border-collapse text-sm">
      <tbody>
        {matrix.map((row, i) => (
          <tr key={i}>
            {row.map((value, j) => (
              <td key={j} className="border border-slate-300 px-2 py-1 text-right tabular-nums">
                {value.toFixed(4)}
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
