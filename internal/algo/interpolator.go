package algo

import (
	"comp-math-5/internal/numeric"
)

type InterpolationResult struct {
	Method string          `json:"method"`
	XValue float64         `json:"xValue"`
	YValue float64         `json:"yValue"`
	Table  [][]float64     `json:"table"`
	Curve  []numeric.Point `json:"curve"`
	Error  string          `json:"error,omitempty"`
}

func Interpolate(points []numeric.Point, xUser float64) ([]InterpolationResult, error) {
	var results []InterpolationResult

	resL := LagrangeInterpolation(points, xUser)
	curveL := generateCurve(points, func(x float64) float64 {
		return LagrangeInterpolation(points, x)
	})
	results = append(results, InterpolationResult{
		Method: "Lagrange",
		XValue: xUser,
		YValue: resL,
		Table:  finiteDifferencesTable(points),
		Curve:  curveL,
	})

	tableFin := finiteDifferencesTable(points)

	tableDiv := dividedDifferencesTable(points)
	resND1 := NewtonDividedForwardInterpolation(points, xUser)
	curveND1 := generateCurve(points, func(x float64) float64 {
		return NewtonDividedForwardInterpolation(points, x)
	})
	results = append(results, InterpolationResult{
		Method: "Newton Forward (Divided Differences)",
		XValue: xUser,
		YValue: resND1,
		Table:  tableDiv,
		Curve:  curveND1,
	})

	resND2 := NewtonDividedBackwardInterpolation(points, xUser)
	curveND2 := generateCurve(points, func(x float64) float64 {
		return NewtonDividedBackwardInterpolation(points, x)
	})
	results = append(results, InterpolationResult{
		Method: "Newton Backward (Divided Differences)",
		XValue: xUser,
		YValue: resND2,
		Table:  tableDiv,
		Curve:  curveND2,
	})

	resGF1, errGF1 := GaussForwardInterpolation(points, xUser)
	if errGF1 == nil {
		curveGF := generateCurve(points, func(x float64) float64 {
			v, _ := GaussForwardInterpolation(points, x)
			return v
		})
		results = append(results, InterpolationResult{
			Method: "Gauss Forward",
			XValue: xUser,
			YValue: resGF1,
			Table:  tableFin,
			Curve:  curveGF,
		})
	}

	resGF2, errGF2 := GaussBackwardInterpolation(points, xUser)
	if errGF2 == nil {
		curveGF := generateCurve(points, func(x float64) float64 {
			v, _ := GaussBackwardInterpolation(points, x)
			return v
		})
		results = append(results, InterpolationResult{
			Method: "Gauss Backward",
			XValue: xUser,
			YValue: resGF2,
			Table:  tableFin,
			Curve:  curveGF,
		})
	}

	resS, errS := StirlingInterpolation(points, xUser)
	if errS == nil {
		curveS := generateCurve(points, func(x float64) float64 {
			v, _ := StirlingInterpolation(points, x)
			return v
		})
		results = append(results, InterpolationResult{
			Method: "Stirling",
			XValue: xUser,
			YValue: resS,
			Table:  tableFin,
			Curve:  curveS,
		})
	}

	resB, errB := BesselInterpolation(points, xUser)
	if errB == nil {
		curveB := generateCurve(points, func(x float64) float64 {
			v, _ := BesselInterpolation(points, x)
			return v
		})
		results = append(results, InterpolationResult{
			Method: "Bessel",
			XValue: xUser,
			YValue: resB,
			Table:  tableFin,
			Curve:  curveB,
		})
	}

	return results, nil
}
