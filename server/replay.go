package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func SaveAsPng(path string, g *Game) error {
	b := g.Board

	img := image.NewRGBA(image.Rect(0, 0, b.Width, b.Height))

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			id := b.fields[x][y]
			img.Set(x, y, color.Gray{uint8(100 + id*(150/len(g.players)))})
		}
	}

	f, _ := os.Create(path)

	defer f.Close()
	return png.Encode(f, img)
}

func CreateUniqueFolder(base string) (string, error) {
	path := base
	i := 1
	for true {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path, os.Mkdir(path, 0755)
		}
		path = fmt.Sprintf("%s_%d", base, i)
		i++
	}
	return path, nil
}
