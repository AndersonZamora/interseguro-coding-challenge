package matrix

import (
	"errors"
	"math"
	"testing"
)

// frobeniusDiff calcula la norma de Frobenius de (a - b), asumiendo que
// ambas matrices tienen las mismas dimensiones.
func frobeniusDiff(a, b [][]float64) float64 {
	var sum float64
	for i := range a {
		for j := range a[i] {
			d := a[i][j] - b[i][j]
			sum += d * d
		}
	}
	return math.Sqrt(sum)
}

func multiply(a, b [][]float64) [][]float64 {
	rows, inner, cols := len(a), len(b), len(b[0])
	out := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		out[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			var sum float64
			for k := 0; k < inner; k++ {
				sum += a[i][k] * b[k][j]
			}
			out[i][j] = sum
		}
	}
	return out
}

func transpose(a [][]float64) [][]float64 {
	rows, cols := len(a), len(a[0])
	out := make([][]float64, cols)
	for j := 0; j < cols; j++ {
		out[j] = make([]float64, rows)
		for i := 0; i < rows; i++ {
			out[j][i] = a[i][j]
		}
	}
	return out
}

const tolerance = 1e-9

func assertReconstructs(t *testing.T, a [][]float64, res Result) {
	t.Helper()

	reconstructed := multiply(res.Q, res.R)
	if diff := frobeniusDiff(a, reconstructed); diff > tolerance {
		t.Errorf("‖A - QR‖ = %g, se esperaba < %g\nA=%v\nQR=%v", diff, tolerance, a, reconstructed)
	}

	n := len(res.Q[0])
	identityN := identity(n)
	qtq := multiply(transpose(res.Q), res.Q)
	if diff := frobeniusDiff(identityN, qtq); diff > tolerance {
		t.Errorf("‖QᵀQ - I‖ = %g, se esperaba < %g (Q no es ortogonal)\nQᵀQ=%v", diff, tolerance, qtq)
	}
}

func TestDecompose_SquareMatrix(t *testing.T) {
	a := [][]float64{
		{12, -51, 4},
		{6, 167, -68},
		{-4, 24, -41},
	}
	res, err := Decompose(a)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	assertReconstructs(t, a, res)
}

func TestDecompose_RectangularMatrix(t *testing.T) {
	a := [][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
	}
	res, err := Decompose(a)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(res.Q) != 4 || len(res.Q[0]) != 2 {
		t.Fatalf("dimensiones de Q incorrectas: got %dx%d, want 4x2", len(res.Q), len(res.Q[0]))
	}
	if len(res.R) != 2 || len(res.R[0]) != 2 {
		t.Fatalf("dimensiones de R incorrectas: got %dx%d, want 2x2", len(res.R), len(res.R[0]))
	}
	assertReconstructs(t, a, res)
}

func TestDecompose_IdentityMatrix(t *testing.T) {
	a := identity(3)
	res, err := Decompose(a)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	assertReconstructs(t, a, res)
}

func TestDecompose_EmptyMatrix(t *testing.T) {
	_, err := Decompose([][]float64{})
	if !errors.Is(err, ErrEmptyMatrix) {
		t.Fatalf("got err=%v, want ErrEmptyMatrix", err)
	}
}

func TestDecompose_RaggedMatrix(t *testing.T) {
	a := [][]float64{
		{1, 2, 3},
		{4, 5},
	}
	_, err := Decompose(a)
	if !errors.Is(err, ErrRaggedMatrix) {
		t.Fatalf("got err=%v, want ErrRaggedMatrix", err)
	}
}

func TestDecompose_TooFewRows(t *testing.T) {
	a := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}
	_, err := Decompose(a)
	if !errors.Is(err, ErrTooFewRows) {
		t.Fatalf("got err=%v, want ErrTooFewRows", err)
	}
}

func TestDecompose_RankDeficient(t *testing.T) {
	// La segunda columna es 2x la primera: columnas linealmente dependientes.
	a := [][]float64{
		{1, 2},
		{2, 4},
		{3, 6},
	}
	_, err := Decompose(a)
	if !errors.Is(err, ErrRankDeficient) {
		t.Fatalf("got err=%v, want ErrRankDeficient", err)
	}
}

func TestDecompose_SingleColumn(t *testing.T) {
	a := [][]float64{{3}, {4}}
	res, err := Decompose(a)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	assertReconstructs(t, a, res)
	// ‖a‖ = 5, R debe ser [[-5]] o [[5]] segun el signo elegido para evitar
	// cancelacion; en cualquier caso |R[0][0]| debe ser 5.
	if math.Abs(math.Abs(res.R[0][0])-5) > tolerance {
		t.Errorf("R[0][0] = %v, se esperaba magnitud 5", res.R[0][0])
	}
}
