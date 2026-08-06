package services

import (
	"fmt"
	"math"
)

const Epsilon = 1e-12

// ComputeQR calculates the QR factorization of an m x n rectangular matrix A (with m >= n)
// using the Modified Gram-Schmidt (MGS) orthogonalization algorithm.
//
// Mathematically, A = Q * R, where:
// - Q is an m x n matrix with orthonormal columns (Q^T * Q = I_n).
// - R is an n x n upper triangular matrix (R_ij = 0 for i > j).
//
// Algorithm Rationale:
// The Modified Gram-Schmidt (MGS) method is chosen over Classical Gram-Schmidt because
// it is numerically stable. In MGS, vector projections are subtracted from intermediate
// column vectors iteratively as soon as each orthonormal basis vector is computed,
// minimizing catastrophic floating-point cancellation errors.
//
// Computational Complexity:
// - Time Complexity: O(m * n^2) operations. For each of the n columns, a vector norm O(m)
//   and projection updates across remaining columns O(m * (n - k)) are performed.
// - Space Complexity: O(m * n + n^2) to allocate the resulting Q (m x n) and R (n x n) matrices.
func ComputeQR(matrix [][]float64) (q [][]float64, r [][]float64, err error) {
	if err := ValidateMatrix(matrix); err != nil {
		return nil, nil, err
	}

	m := len(matrix)
	n := len(matrix[0])

	if m < n {
		return nil, nil, fmt.Errorf("%w: el número de filas (%d) debe ser >= al número de columnas (%d)", ErrInsufficientRowsQR, m, n)
	}

	// Initialize V as a copy of the columns of A (m x n)
	// V[i][j] represents row i, column j
	v := make([][]float64, m)
	for i := 0; i < m; i++ {
		v[i] = make([]float64, n)
		copy(v[i], matrix[i])
	}

	// Q is an m x n matrix
	q = make([][]float64, m)
	for i := 0; i < m; i++ {
		q[i] = make([]float64, n)
	}

	// R is an n x n upper triangular matrix
	r = make([][]float64, n)
	for i := 0; i < n; i++ {
		r[i] = make([]float64, n)
	}

	for k := 0; k < n; k++ {
		// Calculate the L2-norm of column k of V
		var normSq float64
		for i := 0; i < m; i++ {
			normSq += v[i][k] * v[i][k]
		}
		norm := math.Sqrt(normSq)

		r[k][k] = norm

		if norm < Epsilon {
			// Column k is linearly dependent or zero vector
			for i := 0; i < m; i++ {
				q[i][k] = 0.0
			}
		} else {
			// Compute orthonormal column vector q_k = v_k / norm
			for i := 0; i < m; i++ {
				q[i][k] = v[i][k] / norm
			}
		}

		// Orthogonalize remaining columns v_j against q_k
		for j := k + 1; j < n; j++ {
			// Compute dot product R[k][j] = q_k . v_j
			var dot float64
			for i := 0; i < m; i++ {
				dot += q[i][k] * v[i][j]
			}
			r[k][j] = dot

			// Subtract projection: v_j = v_j - R[k][j] * q_k
			for i := 0; i < m; i++ {
				v[i][j] -= dot * q[i][k]
			}
		}
	}

	// Ensure strict upper triangularity for R (explicitly zero out lower triangular elements)
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			r[i][j] = 0.0
		}
	}

	return q, r, nil
}
