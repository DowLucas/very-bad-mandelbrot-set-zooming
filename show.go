package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"mandelbrotset/mandlebrot"
	"math"
	"sync"
)

const detail = 3

type Point struct {
	x float64
	y float64
}

type ScreenLimits struct {
	left float64
	top float64
	right float64
	bottom float64
}


const width = 500
const height = 300
const zoomXSpeed = 1.1
const zoomYSpeed = 1.1
var zoomX float64 = 0.99
var zoomY float64 = 0.99
var doUpdate = true
var completedZoom = false
var zooming = false


func GetMidPoint() *Point {
	midP := Point{width/2, height/2}
	return &midP
}

func GetScreenLimits(mid Point) *ScreenLimits {
	//midX := mid.x
	//midY := mid.y

	left := -2.5 * zoomX
	right := 1 * zoomX
	top :=  1 * zoomY
	bottom :=  -1 * zoomY

	return &ScreenLimits{left, top, right, bottom}
}

func MapNumbers(a float64, min float64, max float64, from float64, to float64) float64 {
	return (a-min) / (max - min) * (to - from) + from
}

func TranslateScreenCoordinates(x float64, y float64, mid *Point) mandelbrot.ComplexNumber {
	return mandelbrot.ComplexNumber{
		Re: MapNumbers(x, -mid.x*zoomX, mid.x*zoomX, -2.5, 1),
		Im: MapNumbers(y, -mid.y*zoomY, mid.y*zoomY, -1, 1)}
}

func DrawPoints(pixelValues [][]bool, imd *imdraw.IMDraw) {
	for dy := range pixelValues {
		for dx := range pixelValues[dy] {
			if pixelValues[dy][dx] {
				imd.Push(pixel.V(float64(dx), float64(dy)))
			}
		}
	}
}

func Update(mid* Point, ch chan<- Point, wg* sync.WaitGroup) {
	defer wg.Done()

	for dy := -mid.y; dy < mid.y; dy += detail {
		for dx := -mid.x; dx < mid.x; dx += detail {
			complexPoint := TranslateScreenCoordinates(dx, -dy, mid)
			if mandelbrot.InMandelbrotSet(&complexPoint) {
				arr := Point{dx+mid.x, dy+mid.y}
				ch <- arr
			}
		}
	}
}

func DrawPoint(ch <-chan Point, imd *imdraw.IMDraw) {
	for i := range ch {
		imd.Push(pixel.V(i.x, i.y))
	}
}

func run() {
	mid := GetMidPoint()
	wg := new(sync.WaitGroup)
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	var ch chan Point

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		imd.Color = pixel.RGB(0, 0, 1)

		if adjustedZoom := win.MouseScroll().Y; adjustedZoom != 0.0 {
			zoomX *= math.Pow(zoomXSpeed, adjustedZoom)
			zoomY *= math.Pow(zoomYSpeed, adjustedZoom)

			doUpdate = true
			zooming = true
		} else if zooming {
			zooming = false
			doUpdate = true
		}


		if doUpdate {
			ch = make(chan Point)
			wg.Add(1)
			go DrawPoint(ch, imd)
			go Update(mid, ch, wg)
			wg.Wait()
			imd.Circle(1, 1)
			doUpdate = false
		}

		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}