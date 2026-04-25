package algo

import (
	"comp-math-5/internal/numeric"
	"fmt"
)

type InterpolationResult struct {
	Method string          `json:"method"`
	XValue float64         `json:"xValue"` // Точка, введенная пользователем
	YValue float64         `json:"yValue"` // Результат в этой точке
	Table  [][]float64     `json:"table"`  // Таблица конечных/разделенных разностей
	Curve  []numeric.Point `json:"curve"`  // Точки для отрисовки графика (линия)
	Nodes  []numeric.Point `json:"nodes"`  // Исходные узлы (для точек на графике)
	Error  string          `json:"error,omitempty"`
}

func Interpolate(points []numeric.Point, x1, x2 float64) ([]InterpolationResult, error) {
	var results []InterpolationResult

	// --- Многочлен Лагранжа ---
	resL1 := LagrangeInterpolation(points, x1)
	results = append(results, InterpolationResult{Method: "Lagrange", XValue: x1, YValue: resL1})

	resL2 := LagrangeInterpolation(points, x2)
	results = append(results, InterpolationResult{Method: "Lagrange", XValue: x2, YValue: resL2})

	// --- Многочлен Ньютона с разделенными разностями ---
	resND1 := NewtonDividedForwardInterpolation(points, x1)
	results = append(results, InterpolationResult{Method: "Newton Forward (Divided)", XValue: x1, YValue: resND1})

	resND2 := NewtonDividedBackwardInterpolation(points, x2)
	results = append(results, InterpolationResult{Method: "Newton Backward (Divided)", XValue: x2, YValue: resND2})

	// --- Многочлены Гаусса ---
	resG1, errG1 := GaussForwardInterpolation(points, x1)
	if errG1 != nil {
		fmt.Printf("Gauss Forward for X1 skipped: %v\n", errG1)
		results = append(results, InterpolationResult{Method: "Gauss Forward", XValue: x1, Error: errG1.Error()})
	} else {
		results = append(results, InterpolationResult{Method: "Gauss Forward", XValue: x1, YValue: resG1})
	}

	resG2, errG2 := GaussBackwardInterpolation(points, x2)
	if errG2 != nil {
		fmt.Printf("Gauss Backward for X2 skipped: %v\n", errG2)
		results = append(results, InterpolationResult{Method: "Gauss Backward", XValue: x2, Error: errG2.Error()})
	} else {
		results = append(results, InterpolationResult{Method: "Gauss Backward", XValue: x2, YValue: resG2})
	}

	// --- Схема Стирлинга ---
	resS1, errS1 := StirlingInterpolation(points, x1)
	if errS1 != nil {
		fmt.Printf("Stirling for X1 skipped: %v\n", errS1)
		results = append(results, InterpolationResult{Method: "Stirling", XValue: x1, Error: errS1.Error()})
	} else {
		results = append(results, InterpolationResult{Method: "Stirling", XValue: x1, YValue: resS1})
	}

	resS2, errS2 := StirlingInterpolation(points, x2)
	if errS2 != nil {
		fmt.Printf("Stirling for X2 skipped: %v\n", errS2)
		results = append(results, InterpolationResult{Method: "Stirling", XValue: x2, Error: errS2.Error()})
	} else {
		results = append(results, InterpolationResult{Method: "Stirling", XValue: x2, YValue: resS2})
	}

	// --- Схема Бесселя ---
	resB1, errB1 := BesselInterpolation(points, x1)
	if errB1 != nil {
		fmt.Printf("Bessel for X1 skipped: %v\n", errB1)
		results = append(results, InterpolationResult{Method: "Bessel", XValue: x1, Error: errB1.Error()})
	} else {
		results = append(results, InterpolationResult{Method: "Bessel", XValue: x1, YValue: resB1})
	}

	resB2, errB2 := BesselInterpolation(points, x2)
	if errB2 != nil {
		fmt.Printf("Bessel for X2 skipped: %v\n", errB2)
		results = append(results, InterpolationResult{Method: "Bessel", XValue: x2, Error: errB2.Error()})
	} else {
		results = append(results, InterpolationResult{Method: "Bessel", XValue: x2, YValue: resB2})
	}

	return results, nil
}
