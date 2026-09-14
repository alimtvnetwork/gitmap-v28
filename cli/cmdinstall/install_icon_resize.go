package cmdinstall

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func decodeSourceImage(sourceBytes []byte) (image.Image, error) {
	if len(sourceBytes) == 0 {
		return nil, apperror.NewSimple("empty icon source bytes", "E1000")
	}

	img, err := png.Decode(bytes.NewReader(sourceBytes))
	if err != nil {
		return nil, apperror.WrapSimple(err, "decode.iconPng")
	}

	return img, nil
}

func clampCoord(val, minVal, maxVal int) int {
	if val < minVal {
		return minVal
	}
	if val > maxVal {
		return maxVal
	}

	return val
}

func interpolateChannel(c00, c10, c01, c11, fx, fy float64) uint8 {
	top := c00*(1.0-fx) + c10*fx
	bot := c01*(1.0-fx) + c11*fx
	val := top*(1.0-fy) + bot*fy

	return uint8(math.Round(val))
}

func sampleBilinearPixel(src image.Image, x0, y0, x1, y1 int, fx, fy float64) color.RGBA {
	r00, g00, b00, a00 := src.At(x0, y0).RGBA()
	r10, g10, b10, a10 := src.At(x1, y0).RGBA()
	r01, g01, b01, a01 := src.At(x0, y1).RGBA()
	r11, g11, b11, a11 := src.At(x1, y1).RGBA()

	return color.RGBA{
		R: interpolateChannel(float64(r00>>8), float64(r10>>8), float64(r01>>8), float64(r11>>8), fx, fy),
		G: interpolateChannel(float64(g00>>8), float64(g10>>8), float64(g01>>8), float64(g11>>8), fx, fy),
		B: interpolateChannel(float64(b00>>8), float64(b10>>8), float64(b01>>8), float64(b11>>8), fx, fy),
		A: interpolateChannel(float64(a00>>8), float64(a10>>8), float64(a01>>8), float64(a11>>8), fx, fy),
	}
}

func computeSampleCoords(
	x int,
	y int,
	targetSize int,
	bounds image.Rectangle,
) (int, int, int, int, float64, float64) {
	srcX := (float64(x)+0.5)*float64(bounds.Dx())/float64(targetSize) - 0.5
	srcY := (float64(y)+0.5)*float64(bounds.Dy())/float64(targetSize) - 0.5
	x0 := clampCoord(int(math.Floor(srcX)), bounds.Min.X, bounds.Max.X-1)
	y0 := clampCoord(int(math.Floor(srcY)), bounds.Min.Y, bounds.Max.Y-1)
	x1 := clampCoord(x0+1, bounds.Min.X, bounds.Max.X-1)
	y1 := clampCoord(y0+1, bounds.Min.Y, bounds.Max.Y-1)

	return x0, y0, x1, y1, math.Max(0, math.Min(1, srcX-float64(x0))), math.Max(0, math.Min(1, srcY-float64(y0)))
}

func resizeImageSquare(src image.Image, targetSize int) *image.RGBA {
	dest := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))
	bounds := src.Bounds()

	for y := 0; y < targetSize; y++ {
		for x := 0; x < targetSize; x++ {
			x0, y0, x1, y1, fx, fy := computeSampleCoords(x, y, targetSize, bounds)
			dest.SetRGBA(x, y, sampleBilinearPixel(src, x0, y0, x1, y1, fx, fy))
		}
	}

	return dest
}

func encodePngToBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, apperror.WrapSimple(err, "encode.iconPng")
	}

	return buf.Bytes(), nil
}

func resizeToSingleResolution(src image.Image, size int) (ResizedIcon, error) {
	resizedImg := resizeImageSquare(src, size)
	encodedData, err := encodePngToBytes(resizedImg)
	if err != nil {
		return ResizedIcon{}, err
	}

	return ResizedIcon{Size: size, Data: encodedData}, nil
}

func resizeAllResolutions(srcImg image.Image) ([]ResizedIcon, error) {
	icons := make([]ResizedIcon, 0, len(StandardIconResolutions))
	for _, size := range StandardIconResolutions {
		icon, err := resizeToSingleResolution(srcImg, size)
		if err != nil {
			return nil, err
		}

		icons = append(icons, icon)
	}

	return icons, nil
}

func GenerateMultiResolutionIcons(sourceBytes []byte) ([]ResizedIcon, error) {
	srcImg, err := decodeSourceImage(sourceBytes)
	if err != nil {
		return nil, err
	}

	return resizeAllResolutions(srcImg)
}
