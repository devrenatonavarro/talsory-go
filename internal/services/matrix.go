package services

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var (
	ErrEmptyMatrix        = errors.New("la matriz no puede estar vacía")
	ErrEmptyRow           = errors.New("la matriz no puede contener filas vacías")
	ErrNonRectangular     = errors.New("todas las filas de la matriz deben tener la misma longitud")
	ErrInvalidNumber      = errors.New("la matriz contiene valores numéricos no válidos (NaN o Infinito)")
	ErrInvalidDirection   = errors.New("la dirección de rotación debe ser 'clockwise' o 'counterclockwise'")
	ErrInsufficientRowsQR = errors.New("para la factorización QR, el número de filas (m) debe ser mayor o igual al número de columnas (n)")
)

// ValidateMatrix checks whether a matrix is rectangular, non-empty, and contains valid numbers.
func ValidateMatrix(matrix [][]float64) error {
	if len(matrix) == 0 {
		return ErrEmptyMatrix
	}

	cols := len(matrix[0])
	if cols == 0 {
		return ErrEmptyRow
	}

	for i, row := range matrix {
		if len(row) == 0 {
			return ErrEmptyRow
		}
		if len(row) != cols {
			return fmt.Errorf("%w: fila %d tiene longitud %d, se esperaba %d", ErrNonRectangular, i, len(row), cols)
		}
		for j, val := range row {
			if math.IsNaN(val) || math.IsInf(val, 0) {
				return fmt.Errorf("%w en la posición [%d][%d]", ErrInvalidNumber, i, j)
			}
		}
	}

	return nil
}

// RotateMatrix rotates an m x n matrix by 90 degrees in the specified direction.
// Supported directions are "clockwise" (default) and "counterclockwise".
// Time complexity: O(m * n) where m is rows and n is columns.
// Space complexity: O(m * n) to allocate the rotated matrix.
func RotateMatrix(matrix [][]float64, direction string) ([][]float64, error) {
	if err := ValidateMatrix(matrix); err != nil {
		return nil, err
	}

	dir := strings.ToLower(strings.TrimSpace(direction))
	if dir == "" {
		dir = "clockwise"
	}

	if dir != "clockwise" && dir != "counterclockwise" {
		return nil, ErrInvalidDirection
	}

	m := len(matrix)
	n := len(matrix[0])

	// The rotated matrix will have dimensions n x m
	rotated := make([][]float64, n)
	for i := range rotated {
		rotated[i] = make([]float64, m)
	}

	if dir == "clockwise" {
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				rotated[j][m-1-i] = matrix[i][j]
			}
		}
	} else { // counterclockwise
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				rotated[n-1-j][i] = matrix[i][j]
			}
		}
	}

	return rotated, nil
}
