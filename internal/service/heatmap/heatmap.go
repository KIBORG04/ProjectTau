package heatmap

import (
	"github.com/PerformLine/go-stockutil/colorutil"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"os"
)

const MapSize = 255

type Mapper struct {
	src  image.Image
	dest *image.NRGBA
}

func NewMapper(_src image.Image) Mapper {
	mapper := Mapper{src: _src}
	mapper.dest = mapper.image2RGBA()
	return mapper
}

func (m *Mapper) image2RGBA() *image.NRGBA {
	b := m.src.Bounds()
	img := image.NewNRGBA(b)
	//draw.Draw(img, b, &image.Uniform{C: color.Transparent}, image.Point{}, draw.Src)
	draw.Draw(img, b, m.src, image.Point{}, draw.Src)
	return img
}

func (m *Mapper) tileCoords2Point(x, y int) image.Point {
	b := m.src.Bounds()
	tileSize := b.Max.X / MapSize
	someShitX := 0
	someShitY := 0
	if tileSize == 4 {
		someShitX = 5
		someShitY = 1025
	} else if tileSize == 32 {
		someShitX = tileSize
		someShitY = tileSize * MapSize
	}
	return image.Point{
		X: (x * tileSize) - someShitX,
		Y: someShitY - (y * tileSize),
	}
}

func (m *Mapper) drawTile(coords image.Point, color *image.Uniform) {
	x := coords.X
	y := coords.Y
	b := m.src.Bounds()
	tileSize := b.Max.X / MapSize
	draw.Draw(m.dest, image.Rect(x, y, x+tileSize, y+tileSize), color, image.Point{}, draw.Over)
}

func (m *Mapper) getHeatColor(x, max float64) *image.Uniform {
	x = math.Log(x)
	max = math.Log(max)
	coef := math.Sqrt(x / max)
	y := math.Abs(240 - 240*coef)
	r, g, b := colorutil.HslToRgb(y, 1, 0.5)
	return &image.Uniform{C: color.NRGBA{r, g, b, 128}}
}

func (m *Mapper) SaveTo(img io.Writer) error {
	err := png.Encode(img, m.dest)
	if err != nil {
		return err
	}
	return nil
}

func Create(outName string, data []image.Point) error {
	file, err := os.Open("boxstation-1.png")
	if err != nil {
		return err
	}
	defer file.Close()

	heatMapOverlay, _ := os.Create(outName + ".png")
	defer heatMapOverlay.Close()

	mapPng, err := png.Decode(file)
	if err != nil {
		return err
	}
	mapper := NewMapper(mapPng)

	mapData := make([][]int, 255)
	for i := 0; i < len(mapData); i++ {
		mapData[i] = make([]int, 255)
	}

	max := 0
	for _, p := range data {
		mapData[p.X][p.Y] += 1
		if mapData[p.X][p.Y] > max {
			max = mapData[p.X][p.Y]
		}
	}

	for i := 0; i < len(mapData); i++ {
		for j := 0; j < len(mapData); j++ {
			if mapData[i][j] == 0 {
				continue
			}
			point := mapper.tileCoords2Point(i, j)
			clr := mapper.getHeatColor(float64(mapData[i][j]), float64(max))
			mapper.drawTile(point, clr)
		}
	}

	err = mapper.SaveTo(heatMapOverlay)
	if err != nil {
		return err
	}

	return nil
}
