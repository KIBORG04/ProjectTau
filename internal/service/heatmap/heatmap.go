package heatmap

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/PerformLine/go-stockutil/colorutil"
	"golang.org/x/exp/slices"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"net/http"
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
	base64.StdEncoding.EncodeToString(m.dest.Pix)
	err := png.Encode(img, m.dest)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mapper) Create(ctx context.Context, outName string, data []image.Point) error {
	file, err := os.Open("boxstation-1.png")
	if err != nil {
		return err
	}
	defer file.Close()

	mapData := make([][]int, 255)
	for i := 0; i < len(mapData); i++ {
		mapData[i] = make([]int, 255)
		if ContextIsDone(ctx) {
			return fmt.Errorf("context is done")
		}
	}

	max := 0
	for _, p := range data {
		mapData[p.X][p.Y] += 1
		if mapData[p.X][p.Y] > max {
			max = mapData[p.X][p.Y]
		}
		if ContextIsDone(ctx) {
			return fmt.Errorf("context is done")
		}
	}

	for i := 0; i < len(mapData); i++ {
		for j := 0; j < len(mapData); j++ {
			if mapData[i][j] == 0 {
				continue
			}
			point := m.tileCoords2Point(i, j)
			clr := m.getHeatColor(float64(mapData[i][j]), float64(max))
			m.drawTile(point, clr)
			if ContextIsDone(ctx) {
				return fmt.Errorf("context is done")
			}
		}
	}

	if ContextIsDone(ctx) {
		return fmt.Errorf("context is done")
	}
	heatMapOverlayName := outName + ".png"
	heatMapOverlay, _ := os.Create(heatMapOverlayName)
	defer heatMapOverlay.Close()

	if ContextIsDone(ctx) {
		os.Remove(heatMapOverlayName)
		return fmt.Errorf("context is done")
	}
	err = m.SaveTo(heatMapOverlay)
	if err != nil {
		return err
	}

	return nil
}

func ContextIsDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		println("ctx gg done aye")
		return true
	default:
		return false
	}
}

func GetMap(ctx context.Context, mapName, mapType, mapResolution string) (string, error) {
	map2nanomap := map[string]string{
		"boxStation": "https://cdn.jsdelivr.net/gh/TauCetiStation/TauCetiClassic@latest/nano/images/nanomap_exodus_1.png",
		"falcon":     "https://cdn.jsdelivr.net/gh/TauCetiStation/TauCetiClassic@latest/nano/images/nanomap_falcon_1.png",
		"gamma":      "https://cdn.jsdelivr.net/gh/TauCetiStation/TauCetiClassic@latest/nano/images/nanomap_gamma_1.png",
		"prometheus": "https://cdn.jsdelivr.net/gh/TauCetiStation/TauCetiClassic@latest/nano/images/nanomap_prometheus_1.png",
		//"asteroid":   "https://cdn.jsdelivr.net/gh/TauCetiStation/TauCetiClassic@latest/nano/images/", верим, надеемся
	}

	map2webmap := map[string]string{
		"boxStation": "https://mocha.affectedarc07.co.uk/webmap/tcc/box/boxstation-1.png",
		"falcon":     "https://mocha.affectedarc07.co.uk/webmap/tcc/falcon/falcon-1.png",
		"gamma":      "https://mocha.affectedarc07.co.uk/webmap/tcc/gamma/gamma-1.png",
		// "prometheus": "https://mocha.affectedarc07.co.uk/webmap/tcc/", мда.....
		"asteroid": "https://mocha.affectedarc07.co.uk/webmap/tcc/asteroid/asteroid-1.png",
	}

	var mapLink string
	switch mapResolution {
	case "nano":
		mapLink = map2nanomap[mapName]
	case "webmap":
		mapLink = map2webmap[mapName]
	default:
		return "", fmt.Errorf("the mapresolution can be nano or webmap")
	}

	validMapTypes := []string{
		"deaths", "explosions",
	}
	if !slices.Contains(validMapTypes, mapType) {
		return "", fmt.Errorf("the maptype can be deaths or explosions")
	}

	get, err := http.Get(mapLink)
	if err != nil {
		return "", err
	}

	_ = get
	/*
		mapPng, err := png.Decode(file)
		if err != nil {
			return err
		}

		//mapper := NewMapper(mapPng)
	*/
	return "poebat'", nil
}
