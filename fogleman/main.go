package main

import "github.com/fogleman/gg"

func main() {
	const S = 1000
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.Translate(S/2, S/2)
	dc.Scale(40, 40)

	var x0, y0, x1, y1, x2, y2, x3, y3 float64
	x0, y0 = -10, 0
	x1, y1 = -8, -8
	x2, y2 = 8, 8
	x3, y3 = 10, 0

	dc.MoveTo(x0, y0)
	dc.CubicTo(x1, y1, x2, y2, x3, y3)
	dc.SetRGBA(0, 0, 0, 0.2)
	dc.SetLineWidth(8)
	dc.FillPreserve()
	dc.SetRGB(0, 0, 0)
	dc.SetDash(16, 24)
	dc.Stroke()

	dc.SavePNG("out1.png")
}
