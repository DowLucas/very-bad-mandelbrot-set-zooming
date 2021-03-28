package mandelbrot

type ComplexNumber struct {
	Re float64
	Im float64
}

const IterationLimit = 1000
const Threshold = 2e64

func SquareComplexNumber(z *ComplexNumber) {
	re := z.Re * z.Re - z.Im * z.Im
	im := 2 * z.Re * z.Im
	z.Re = re
	z.Im = im
}

func InMandelbrotSet(c* ComplexNumber) bool {
	x0 := c.Re
	y0 := c.Im
	x := 0.0
	y := 0.0
	x2 := 0.0
	y2 := 0.0

	iteration := 0
	for  ;iteration < IterationLimit && x2 + y2 <= 4; {
		y = 2 * x * y + y0
		x = x2 - y2 + x0
		x2 = x * x
		y2 = y * y

		iteration++
	}

	return iteration == 1000
}

func InSet(dx int, dy int, c* ComplexNumber, pixelValues [][]bool) bool {
	return InMandelbrotSet(c)
}

func FromFloats(re float64, im float64) bool {
	c := ComplexNumber{re, im}

	return InMandelbrotSet(&c)
}


