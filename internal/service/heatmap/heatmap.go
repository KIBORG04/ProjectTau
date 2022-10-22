package heatmap

import (
	"github.com/PerformLine/go-stockutil/colorutil"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

func image2RGBA(src image.Image) *image.RGBA {
	b := src.Bounds()
	img := image.NewRGBA(b)
	draw.Draw(img, b, src, image.Point{}, draw.Src)
	return img
}

func tileCoords2Point(x, y int) image.Point {
	// 4 is size of 1 tile on nano map image
	return image.Point{
		X: (x * 4) - 5,
		Y: 1025 - (y * 4),
	}
}

func drawTiles(src *image.RGBA, coords image.Point, color *image.Uniform) {
	x := coords.X
	y := coords.Y
	draw.Draw(src, image.Rect(x, y, x+4, y+4), color, image.Point{}, draw.Over)
}

func getHeatColor(x, max float64) *image.Uniform {
	x = math.Log(x)
	max = math.Log(max)
	coef := math.Sqrt(x / max)
	y := math.Abs(240 - 240*coef)
	r, g, b := colorutil.HslToRgb(y, 1, 0.5)
	return &image.Uniform{C: color.RGBA{uint8(r), uint8(g), uint8(b), 255}}
}

func New(name string, data []image.Point) {
	file, err := os.Open("nanomap_exodus_1.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tempFile, _ := os.Create(name + ".png")
	defer tempFile.Close()

	mapPng, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	mapImage := image2RGBA(mapPng)

	duplicates := make([][]int, 255)
	for i := 0; i < len(duplicates); i++ {
		duplicates[i] = make([]int, 255)
	}

	max := 0
	for _, p := range data {
		duplicates[p.X][p.Y] += 1
		if duplicates[p.X][p.Y] > max {
			max = duplicates[p.X][p.Y]
		}
	}

	for i := 0; i < len(duplicates); i++ {
		for j := 0; j < len(duplicates); j++ {
			if duplicates[i][j] == 0 {
				continue
			}
			point := tileCoords2Point(i, j)
			clr := getHeatColor(float64(duplicates[i][j]), float64(max))
			drawTiles(mapImage, point, clr)
		}
	}

	err = png.Encode(tempFile, mapImage)
	if err != nil {
		panic(err)
	}
}
