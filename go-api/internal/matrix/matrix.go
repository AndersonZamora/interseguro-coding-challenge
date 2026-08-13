// Package matrix contiene utilidades de validación y factorización de matrices.
package matrix

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyMatrix    = errors.New("la matriz no puede estar vacia")
	ErrRaggedMatrix   = errors.New("todas las filas de la matriz deben tener la misma longitud")
	ErrTooFewRows     = errors.New("la matriz debe tener al menos tantas filas como columnas (m >= n) para calcular una factorizacion QR reducida")
	ErrRankDeficient  = errors.New("la matriz tiene columnas linealmente dependientes (rango incompleto); no se puede calcular una factorizacion QR reducida")
)

// Validate verifica que la matriz sea rectangular (todas las filas con igual
// longitud) y no este vacia. Devuelve las dimensiones (filas, columnas).
func Validate(a [][]float64) (rows, cols int, err error) {
	rows = len(a)
	if rows == 0 {
		return 0, 0, ErrEmptyMatrix
	}
	cols = len(a[0])
	if cols == 0 {
		return 0, 0, ErrEmptyMatrix
	}
	for i, row := range a {
		if len(row) != cols {
			return 0, 0, fmt.Errorf("%w: fila 0 tiene %d columnas, fila %d tiene %d", ErrRaggedMatrix, cols, i, len(row))
		}
	}
	return rows, cols, nil
}
