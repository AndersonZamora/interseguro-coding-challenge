package matrix

import "math"

// rankEpsilon es la tolerancia relativa usada para detectar columnas
// linealmente dependientes durante la factorizacion (rango incompleto).
const rankEpsilon = 1e-10

// Result contiene la factorizacion QR reducida de una matriz A (m x n, m >= n),
// tal que A = Q * R, con Q (m x n) de columnas ortonormales y R (n x n)
// triangular superior.
type Result struct {
	Q [][]float64 // m x n
	R [][]float64 // n x n
}

// Decompose calcula la factorizacion QR reducida de "a" usando reflexiones de
// Householder (preferidas sobre Gram-Schmidt por su mejor estabilidad
// numerica, ver docs/decisiones-arquitectura.md). Requiere que "a" sea
// rectangular, no vacia, y con rows >= cols. Devuelve ErrEmptyMatrix,
// ErrRaggedMatrix, ErrTooFewRows o ErrRankDeficient segun corresponda; en
// cualquier otro caso de exito, err es nil.
func Decompose(a [][]float64) (Result, error) {
	rows, cols, err := Validate(a)
	if err != nil {
		return Result{}, err
	}
	if rows < cols {
		return Result{}, ErrTooFewRows
	}

	// r comienza como una copia de "a" y se transforma en la matriz
	// triangular superior R (m x n) aplicando una reflexion de Householder
	// por columna. qAcc acumula el producto de las reflexiones para obtener
	// la Q completa (m x m); al final se reduce a las primeras "cols" columnas.
	r := cloneMatrix(a)
	qAcc := identity(rows)

	for k := 0; k < cols; k++ {
		// x es la sub-columna k, desde la fila k hasta el final.
		x := make([]float64, rows-k)
		for i := range x {
			x[i] = r[k+i][k]
		}

		normX := norm(x)
		if normX < rankEpsilon {
			return Result{}, ErrRankDeficient
		}

		// alpha tiene signo opuesto a x[0] para evitar cancelacion numerica.
		alpha := -math.Copysign(normX, x[0])

		v := make([]float64, len(x))
		copy(v, x)
		v[0] -= alpha
		normV := norm(v)
		if normV < rankEpsilon {
			// x ya es (aproximadamente) un multiplo de e1: no hace falta
			// reflejar en esta columna, se continua con la siguiente.
			continue
		}
		for i := range v {
			v[i] /= normV
		}

		applyHouseholderLeft(r, v, k)
		applyHouseholderRight(qAcc, v, k)
	}

	q := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		q[i] = qAcc[i][:cols]
	}
	rTri := make([][]float64, cols)
	for i := 0; i < cols; i++ {
		rTri[i] = r[i][:cols]
	}

	if rankDeficient(rTri) {
		return Result{}, ErrRankDeficient
	}

	return Result{Q: q, R: rTri}, nil
}

// applyHouseholderLeft aplica la reflexion H = I - 2vv^T (definida sobre las
// filas k..rows-1) por la izquierda de m, es decir m[k:, k:] = H * m[k:, k:].
func applyHouseholderLeft(m [][]float64, v []float64, k int) {
	cols := len(m[0])
	for j := k; j < cols; j++ {
		var dot float64
		for i := range v {
			dot += v[i] * m[k+i][j]
		}
		for i := range v {
			m[k+i][j] -= 2 * v[i] * dot
		}
	}
}

// applyHouseholderRight aplica la misma reflexion por la derecha de m sobre
// las columnas k..rows-1, acumulando Q = Q * H.
func applyHouseholderRight(m [][]float64, v []float64, k int) {
	rows := len(m)
	for i := 0; i < rows; i++ {
		var dot float64
		for j := range v {
			dot += m[i][k+j] * v[j]
		}
		for j := range v {
			m[i][k+j] -= 2 * dot * v[j]
		}
	}
}

func norm(v []float64) float64 {
	var sum float64
	for _, x := range v {
		sum += x * x
	}
	return math.Sqrt(sum)
}

func identity(n int) [][]float64 {
	m := make([][]float64, n)
	for i := range m {
		m[i] = make([]float64, n)
		m[i][i] = 1
	}
	return m
}

func cloneMatrix(a [][]float64) [][]float64 {
	m := make([][]float64, len(a))
	for i, row := range a {
		m[i] = make([]float64, len(row))
		copy(m[i], row)
	}
	return m
}

// rankDeficient detecta si la diagonal de una matriz triangular superior
// tiene algun valor practicamente nulo, lo que indica columnas dependientes
// que no se detectaron durante la reduccion (caso limite numerico).
func rankDeficient(r [][]float64) bool {
	for i := range r {
		if math.Abs(r[i][i]) < rankEpsilon {
			return true
		}
	}
	return false
}
