package algo

import (
	"comp-math-5/internal/numeric"
	"fmt"
)

// LagrangeInterpolation вычисляет значение функции в точке x с помощью многочлена Лагранжа.
func LagrangeInterpolation(points []numeric.Point, x float64) float64 {
	n := len(points)
	var result float64

	for i := 0; i < n; i++ {
		term := points[i].Y
		for j := 0; j < n; j++ {
			if i != j {
				term *= (x - points[j].X) / (points[i].X - points[j].X)
			}
		}
		result += term
	}
	return result
}

// finiteDifferencesTable строит таблицу конечных разностей.
func finiteDifferencesTable(points []numeric.Point) [][]float64 {
	n := len(points)
	table := make([][]float64, n)
	for i := range table {
		table[i] = make([]float64, n)
		table[i][0] = points[i].Y
	}

	for j := 1; j < n; j++ {
		for i := 0; i < n-j; i++ {
			table[i][j] = table[i+1][j-1] - table[i][j-1]
		}
	}
	return table
}

// GaussForwardInterpolation вычисляет значение функции в точке x с помощью первой интерполяционной формулы Гаусса.
func GaussForwardInterpolation(points []numeric.Point, x float64) (float64, error) {
	n := len(points)
	if n < 2 {
		return 0, fmt.Errorf("нужно как минимум 2 точки для интерполяции")
	}

	h := points[1].X - points[0].X
	for i := 1; i < n-1; i++ {
		if (points[i+1].X - points[i].X) != h {
			return 0, fmt.Errorf("точки должны быть равноотстоящими для формул Гаусса")
		}
	}

	// Находим центральный узел (или ближайший к x, если n четное)
	midIndex := (n - 1) / 2
	x0 := points[midIndex].X
	t := (x - x0) / h

	table := finiteDifferencesTable(points)

	result := table[midIndex][0] // y_mid

	tProd := 1.0
	for k := 1; k < n; k++ {
		if k%2 != 0 { // Нечетные разности
			tProd *= (t - float64(k/2))
			result += tProd / factorial(k) * table[midIndex-(k/2)][k]
		} else { // Четные разности
			tProd *= (t + float64(k/2-1))
			result += tProd / factorial(k) * table[midIndex-k/2][k]
		}
	}
	return result, nil
}

// GaussBackwardInterpolation вычисляет значение функции в точке x с помощью второй интерполяционной формулы Гаусса.
func GaussBackwardInterpolation(points []numeric.Point, x float64) (float64, error) {
	n := len(points)
	if n < 2 {
		return 0, fmt.Errorf("нужно как минимум 2 точки для интерполяции")
	}

	h := points[1].X - points[0].X
	for i := 1; i < n-1; i++ {
		if (points[i+1].X - points[i].X) != h {
			return 0, fmt.Errorf("точки должны быть равноотстоящими для формул Гаусса")
		}
	}

	// Находим центральный узел (или ближайший к x, если n четное)
	midIndex := (n - 1) / 2
	x0 := points[midIndex].X
	t := (x - x0) / h

	table := finiteDifferencesTable(points)

	result := table[midIndex][0] // y_mid

	tProd := 1.0
	for k := 1; k < n; k++ {
		if k%2 != 0 { // Нечетные разности
			tProd *= (t + float64(k/2))
			result += tProd / factorial(k) * table[midIndex-(k/2)][k]
		} else { // Четные разности
			tProd *= (t - float64(k/2-1))
			result += tProd / factorial(k) * table[midIndex-k/2][k]
		}
	}
	return result, nil
}

// --- 1. Многочлен Ньютона с разделенными разностями ---

// dividedDifferencesTable строит таблицу разделенных разностей.
func dividedDifferencesTable(points []numeric.Point) [][]float64 {
	n := len(points)
	table := make([][]float64, n)
	for i := range table {
		table[i] = make([]float64, n)
		table[i][0] = points[i].Y
	}

	for j := 1; j < n; j++ {
		for i := 0; i < n-j; i++ {
			table[i][j] = (table[i+1][j-1] - table[i][j-1]) / (points[i+j].X - points[i].X)
		}
	}
	return table
}

// NewtonDividedForwardInterpolation - первая формула Ньютона (с разделенными разностями).
func NewtonDividedForwardInterpolation(points []numeric.Point, x float64) float64 {
	table := dividedDifferencesTable(points)
	n := len(points)
	result := table[0][0]
	prod := 1.0

	for i := 1; i < n; i++ {
		prod *= x - points[i-1].X
		result += prod * table[0][i]
	}
	return result
}

func NewtonDividedBackwardInterpolation(points []numeric.Point, x float64) float64 {
	table := dividedDifferencesTable(points)
	n := len(points)
	result := table[n-1][0]
	prod := 1.0

	for i := 1; i < n; i++ {
		prod *= x - points[n-i].X // Идем с конца: x_n, x_{n-1}, ...
		result += prod * table[n-1-i][i]
	}
	return result
}

func calcStirlingTerm(t float64, k int) float64 {
	m := k / 2
	prod := 1.0
	if k%2 != 0 { // Нечетная степень
		prod = t
		for i := 1; i <= m; i++ {
			prod *= t*t - float64(i*i)
		}
	} else { // Четная степень
		prod = t * t
		for i := 1; i < m; i++ {
			prod *= t*t - float64(i*i)
		}
	}
	return prod
}

func StirlingInterpolation(points []numeric.Point, x float64) (float64, error) {
	n := len(points)
	if n < 3 {
		return 0, fmt.Errorf("нужно как минимум 3 точки для интерполяции Стирлинга")
	}

	h := points[1].X - points[0].X
	for i := 1; i < n-1; i++ {
		if (points[i+1].X - points[i].X) != h {
			return 0, fmt.Errorf("точки должны быть равноотстоящими для схемы Стирлинга")
		}
	}

	midIndex := (n - 1) / 2
	x0 := points[midIndex].X
	t := (x - x0) / h
	table := finiteDifferencesTable(points)

	result := table[midIndex][0]

	for k := 1; k < n; k++ {
		term := calcStirlingTerm(t, k)
		if k%2 != 0 {
			// Нечетные разности: берем среднее арифметическое двух центральных
			idx1 := midIndex - k/2
			idx2 := midIndex - k/2 - 1
			if idx1 >= 0 && idx2 >= 0 && idx1 < len(table) && k < len(table[idx1]) {
				avgDiff := (table[idx1][k] + table[idx2][k]) / 2.0
				result += (term / factorial(k)) * avgDiff
			}
		} else {
			// Четные разности: берем одну центральную
			idx := midIndex - k/2
			if idx >= 0 && idx < len(table) && k < len(table[idx]) {
				result += (term / factorial(k)) * table[idx][k]
			}
		}
	}
	return result, nil
}

func calcBesselTerm(t float64, k int) float64 {
	if k == 1 {
		return t - 0.5
	}
	prod := 1.0
	m := k / 2
	for i := 0; i < m; i++ {
		prod *= (t + float64(i)) * (t - float64(i) - 1.0)
	}
	if k%2 != 0 {
		prod *= t - 0.5
	}
	return prod
}

func BesselInterpolation(points []numeric.Point, x float64) (float64, error) {
	n := len(points)
	if n < 4 {
		return 0, fmt.Errorf("нужно как минимум 4 точки для интерполяции Бесселя")
	}

	h := points[1].X - points[0].X
	for i := 1; i < n-1; i++ {
		if (points[i+1].X - points[i].X) != h {
			return 0, fmt.Errorf("точки должны быть равноотстоящими для схемы Бесселя")
		}
	}

	midIndex := (n - 1) / 2
	if midIndex+1 >= n {
		return 0, fmt.Errorf("недостаточно точек для вычисления центрального узла")
	}

	x0 := points[midIndex].X
	t := (x - x0) / h
	table := finiteDifferencesTable(points)

	// k = 0: среднее арифметическое двух центральных значений y
	result := (table[midIndex][0] + table[midIndex+1][0]) / 2.0

	for k := 1; k < n; k++ {
		term := calcBesselTerm(t, k)
		if k%2 != 0 {
			// Нечетные разности: берем одну разность
			m := (k - 1) / 2
			idx := midIndex - m
			if idx >= 0 && idx < len(table) && k < len(table[idx]) {
				result += (term / factorial(k)) * table[idx][k]
			}
		} else {
			// Четные разности: берем среднее арифметическое
			m := k / 2
			idx1 := midIndex - m
			idx2 := midIndex - m + 1
			if idx1 >= 0 && idx2 >= 0 && idx1 < len(table) && k < len(table[idx1]) {
				avgDiff := (table[idx1][k] + table[idx2][k]) / 2.0
				result += (term / factorial(k)) * avgDiff
			}
		}
	}
	return result, nil
}

func factorial(k int) float64 {
	res := 1.0
	for i := 2; i <= k; i++ {
		res *= float64(i)
	}
	return res
}

func generateCurve(points []numeric.Point, f func(float64) float64) []numeric.Point {
	if len(points) < 2 {
		return nil
	}

	minX := points[0].X
	maxX := points[len(points)-1].X
	step := (maxX - minX) / 100

	var curve []numeric.Point
	for i := 0; i <= 100; i++ {
		curX := minX + float64(i)*step
		curve = append(curve, numeric.Point{X: curX, Y: f(curX)})
	}
	return curve
}
