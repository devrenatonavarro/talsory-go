package services

import (
	"math"
	"testing"
)

// Helper function to multiply an m x n matrix Q by an n x n matrix R resulting in an m x n matrix
func multiplyMatrices(q [][]float64, r [][]float64) [][]float64 {
	m := len(q)
	n := len(r)
	result := make([][]float64, m)
	for i := 0; i < m; i++ {
		result[i] = make([]float64, len(r[0]))
		for j := 0; j < len(r[0]); j++ {
			var sum float64
			for k := 0; k < n; k++ {
				sum += q[i][k] * r[k][j]
			}
			result[i][j] = sum
		}
	}
	return result
}

// Helper function to compute Q^T * Q for an m x n matrix Q, resulting in an n x n matrix
func computeQTQ(q [][]float64) [][]float64 {
	m := len(q)
	n := len(q[0])
	qtq := make([][]float64, n)
	for i := 0; i < n; i++ {
		qtq[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			var sum float64
			for k := 0; k < m; k++ {
				sum += q[k][i] * q[k][j]
			}
			qtq[i][j] = sum
		}
	}
	return qtq
}

func TestComputeQR_KnownMatrices(t *testing.T) {
	testCases := []struct {
		name   string
		matrix [][]float64
	}{
		{
			name: "Matriz 3x2 clásica",
			matrix: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
			},
		},
		{
			name: "Matriz 3x3 bien condicionada (Ejemplo de Householder/Gram-Schmidt)",
			matrix: [][]float64{
				{12, -51, 4},
				{6, 167, -68},
				{-4, 24, -41},
			},
		},
		{
			name: "Matriz 2x2 identidad",
			matrix: [][]float64{
				{1, 0},
				{0, 1},
			},
		},
		{
			name: "Matriz 4x2 rectangular",
			matrix: [][]float64{
				{1, -1},
				{1, 1},
				{1, 2},
				{1, 3},
			},
		},
	}

	const tol = 1e-6

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, r, err := ComputeQR(tc.matrix)
			if err != nil {
				t.Fatalf("ComputeQR devolvió un error inesperado: %v", err)
			}

			m := len(tc.matrix)
			n := len(tc.matrix[0])

			// 1. Verificar dimensiones de Q (m x n) y R (n x n)
			if len(q) != m || len(q[0]) != n {
				t.Errorf("Dimensiones de Q incorrectas. Esperado (%d x %d), obtenido (%d x %d)", m, n, len(q), len(q[0]))
			}
			if len(r) != n || len(r[0]) != n {
				t.Errorf("Dimensiones de R incorrectas. Esperado (%d x %d), obtenido (%d x %d)", n, n, len(r), len(r[0]))
			}

			// 2. Verificar que Q * R == A (Reconstrucción dentro de epsilon)
			qr := multiplyMatrices(q, r)
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					if math.Abs(qr[i][j]-tc.matrix[i][j]) > tol {
						t.Errorf("Reconstrucción Q*R falló en [%d][%d]: esperado %f, obtenido %f", i, j, tc.matrix[i][j], qr[i][j])
					}
				}
			}

			// 3. Verificar que Q^T * Q == I_n (Ortogonalidad)
			qtq := computeQTQ(q)
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					expected := 0.0
					if i == j {
						expected = 1.0
					}
					if math.Abs(qtq[i][j]-expected) > tol {
						t.Errorf("Ortogonalidad Q^T * Q falló en [%d][%d]: esperado %f, obtenido %f", i, j, expected, qtq[i][j])
					}
				}
			}

			// 4. Verificar que R es triangular superior (R[i][j] == 0 para i > j)
			for i := 0; i < n; i++ {
				for j := 0; j < i; j++ {
					if math.Abs(r[i][j]) > tol {
						t.Errorf("R no es triangular superior en [%d][%d]: valor es %f", i, j, r[i][j])
					}
				}
			}
		})
	}
}

func TestComputeQR_InsufficientRows(t *testing.T) {
	matrix := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}
	_, _, err := ComputeQR(matrix)
	if err == nil {
		t.Error("Se esperaba un error para matriz con m < n, pero se obtuvo nil")
	}
}
