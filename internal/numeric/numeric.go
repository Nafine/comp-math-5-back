package numeric

import (
	"fmt"
	"math"
	"strings"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ApproximationResult struct {
	Name      string    `json:"name"`
	Coeffs    []float64 `json:"coeffs"`
	S         float64   `json:"s"`                 // Среднеквадратичное отклонение
	R2        float64   `json:"r2"`                // Коэффициент детерминации
	Pearson   *float64  `json:"pearson,omitempty"` // Только для линейной
	Message   string    `json:"message"`           // Сообщение по R2
	XValues   []float64 `json:"x_values"`
	YValues   []float64 `json:"y_values"`
	PhiValues []float64 `json:"phi_values"`
	EpsValues []float64 `json:"eps_values"`
}

func (ar ApproximationResult) String() string {
	var sb strings.Builder

	// Заголовок
	sb.WriteString(fmt.Sprintf("=== %s ===\n", ar.Name))
	sb.WriteString(fmt.Sprintf("Коэффициенты: %v\n", ar.Coeffs))
	sb.WriteString(fmt.Sprintf("Среднеквадратичное отклонение (S): %.6f\n", ar.S))
	sb.WriteString(fmt.Sprintf("Коэффициент детерминации (R²): %.6f\n", ar.R2))
	sb.WriteString(fmt.Sprintf("Сообщение: %s\n", ar.Message))

	// Pearson для линейной функции
	if ar.Pearson != nil {
		sb.WriteString(fmt.Sprintf("Коэффициент Пирсона: %.6f\n", *ar.Pearson))
	}

	// Таблица значений
	sb.WriteString("\n┌─────────┬──────────┬───────────┬───────────┐\n")
	sb.WriteString("│    x    │    y     │  φ(x)     │   ε = y-φ │\n")
	sb.WriteString("├─────────┼──────────┼───────────┼───────────┤\n")

	for i := 0; i < len(ar.XValues); i++ {
		x := ar.XValues[i]
		y := ar.YValues[i]
		phi := ar.PhiValues[i]
		eps := ar.EpsValues[i]

		sb.WriteString(fmt.Sprintf("│ %7.4f │ %8.4f │ %9.4f │ %9.4f │\n", x, y, phi, eps))
	}
	sb.WriteString("└─────────┴──────────┴───────────┴───────────┘\n")

	// Статистика по погрешностям
	if len(ar.EpsValues) > 0 {
		maxEps := ar.EpsValues[0]
		minEps := ar.EpsValues[0]
		sumEps := 0.0
		sumAbsEps := 0.0

		for _, eps := range ar.EpsValues {
			sumEps += eps
			sumAbsEps += math.Abs(eps)
			if eps > maxEps {
				maxEps = eps
			}
			if eps < minEps {
				minEps = eps
			}
		}
		meanEps := sumEps / float64(len(ar.EpsValues))
		meanAbsEps := sumAbsEps / float64(len(ar.EpsValues))

		sb.WriteString("\nСтатистика погрешностей:\n")
		sb.WriteString(fmt.Sprintf("  Средняя ошибка (MSE): %.6f\n", meanEps))
		sb.WriteString(fmt.Sprintf("  Средняя абсолютная ошибка (MAE): %.6f\n", meanAbsEps))
		sb.WriteString(fmt.Sprintf("  Максимальная ошибка: %.6f\n", maxEps))
		sb.WriteString(fmt.Sprintf("  Минимальная ошибка: %.6f\n", minEps))
		sb.WriteString(fmt.Sprintf("  Размах ошибок: %.6f\n", maxEps-minEps))
	}

	// Оценка качества
	sb.WriteString("\nОценка качества аппроксимации:\n")
	if ar.R2 > 0.95 {
		sb.WriteString("  ✓ Отличная аппроксимация (R² > 0.95)\n")
	} else if ar.R2 > 0.8 {
		sb.WriteString("  ✓ Хорошая аппроксимация (R² > 0.8)\n")
	} else if ar.R2 > 0.6 {
		sb.WriteString("  ⚠ Удовлетворительная аппроксимация (R² > 0.6)\n")
	} else {
		sb.WriteString("  ✗ Слабая аппроксимация (R² ≤ 0.6)\n")
	}

	return sb.String()
}
