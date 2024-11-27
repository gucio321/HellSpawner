package common

import (
	"bytes"
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
)

// LoadTexture converts byte slice to image.Image and enqueues new request to giu.
func LoadTexture(fileData []byte, cb func(*g.Texture)) {
	fileReader := bytes.NewReader(fileData)

	rgba, err := convertToImage(fileReader)
	if err != nil {
		log.Fatal(err)
	}

	g.EnqueueNewTextureFromRgba(rgba, cb)
}

func convertToImage(file io.Reader) (*image.RGBA, error) {
	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding png file: %w", err)
	}

	switch trueImg := img.(type) {
	case *image.RGBA:
		return trueImg, nil
	default:
		rgba := image.NewRGBA(trueImg.Bounds())
		draw.Draw(rgba, trueImg.Bounds(), trueImg, image.Pt(0, 0), draw.Src)

		return rgba, nil
	}
}
