package numeric

type Function struct {
	Name string
	Fn   func(float64) float64
}

var functions = []func(x float64) float64{}

func GetFunction(index int) func(float64) float64 {
	return functions[index]
}
