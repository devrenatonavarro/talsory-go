package services

import (
	"math"
	"reflect"
	"testing"
)

func TestValidateMatrix(t *testing.T) {
	tests := []struct {
		name    string
		matrix  [][]float64
		wantErr bool
	}{
		{
			name: "Matriz rectangular válida 3x2",
			matrix: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
			},
			wantErr: false,
		},
		{
			name:    "Matriz vacía",
			matrix:  [][]float64{},
			wantErr: true,
		},
		{
			name:    "Matriz con fila vacía",
			matrix:  [][]float64{{}},
			wantErr: true,
		},
		{
			name: "Filas de distinta longitud (irregular)",
			matrix: [][]float64{
				{1, 2, 3},
				{4, 5},
			},
			wantErr: true,
		},
		{
			name: "Matriz con valor NaN",
			matrix: [][]float64{
				{1, math.NaN()},
				{3, 4},
			},
			wantErr: true,
		},
		{
			name: "Matriz con valor Infinito",
			matrix: [][]float64{
				{1, 2},
				{math.Inf(1), 4},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMatrix(tt.matrix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMatrix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRotateMatrix(t *testing.T) {
	input := [][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
	}

	t.Run("Rotación en sentido horario (clockwise 90deg)", func(t *testing.T) {
		expected := [][]float64{
			{5, 3, 1},
			{6, 4, 2},
		}
		got, err := RotateMatrix(input, "clockwise")
		if err != nil {
			t.Fatalf("RotateMatrix() devolvió un error inesperado: %v", err)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("RotateMatrix() got = %v, want %v", got, expected)
		}
	})

	t.Run("Rotación por defecto (dirección vacía -> clockwise)", func(t *testing.T) {
		expected := [][]float64{
			{5, 3, 1},
			{6, 4, 2},
		}
		got, err := RotateMatrix(input, "")
		if err != nil {
			t.Fatalf("RotateMatrix() devolvió un error inesperado: %v", err)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("RotateMatrix() got = %v, want %v", got, expected)
		}
	})

	t.Run("Rotación en sentido antihorario (counterclockwise 90deg)", func(t *testing.T) {
		expected := [][]float64{
			{2, 4, 6},
			{1, 3, 5},
		}
		got, err := RotateMatrix(input, "counterclockwise")
		if err != nil {
			t.Fatalf("RotateMatrix() devolvió un error inesperado: %v", err)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("RotateMatrix() got = %v, want %v", got, expected)
		}
	})

	t.Run("Dirección de rotación no válida", func(t *testing.T) {
		_, err := RotateMatrix(input, "diagonal")
		if err == nil {
			t.Error("Se esperaba error con dirección no válida, pero se obtuvo nil")
		}
	})
}
