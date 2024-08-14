package service

import (
	"bytes"
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"image"
	"io"
	"io/ioutil"
	"strconv"
)

var name = "svg"
var magicString = "<svg "

func init() {
	image.RegisterFormat(name, magicString, Decode, decodeSVGConfig)
}
func decodeSVG(input []byte, width int, height int) (image.Image, error) {
	in := bytes.NewReader(input)
	icon, err := oksvg.ReadIconStream(in)
	if err != nil {
		return nil, err
	}

	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	if width > 0 && height > 0 {
		w, h = width, height
		icon.SetTarget(0, 0, float64(w), float64(h))
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	scannerGV := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scannerGV)
	icon.Draw(raster, 1.0)
	return img, nil
}
func Decode(r io.Reader) (image.Image, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return decodeSVG(b, 0, 0)
}

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
}

func decodeSVGConfig(reader io.Reader) (image.Config, error) {
	logger := log.WithField("prefix", "SVG decoder::decodeSVGConfig")

	byteValue, readError := io.ReadAll(reader)
	if readError != nil {
		logger.Errorf("failed to read image content: %v", readError)
		return image.Config{}, fmt.Errorf("failed to read image content: %w", readError)
	}
	svg := SVG{}
	xml.Unmarshal(byteValue, &svg)

	width, widthParseError := strconv.Atoi(svg.Width)
	if widthParseError != nil {
		logger.Errorf("failed to parse image width: %v", readError)
		return image.Config{}, fmt.Errorf("failed to parse image width: %w", readError)
	}

	height, heightParseError := strconv.Atoi(svg.Height)
	if heightParseError != nil {
		logger.Errorf("failed to parse image height: %v", readError)
		return image.Config{}, fmt.Errorf("failed to parse image height: %w", readError)
	}

	return image.Config{
		ColorModel: nil,
		Width:      width,
		Height:     height,
	}, nil
}
