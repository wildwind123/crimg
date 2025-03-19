package crimg

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/tiff"
)

// isCMYK checks if the image is in CMYK format.
func IsCMYK(img image.Image) bool {
	_, ok := img.(*image.CMYK)
	return ok
}

// convertCMYKToRGB converts a CMYK image to RGB.
func ConvertCMYKToRGB(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	rgbImg := image.NewNRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			if cmyk, ok := c.(color.CMYK); ok {
				r, g, b := color.CMYKToRGB(cmyk.C, cmyk.M, cmyk.Y, cmyk.K)
				rgbImg.Set(x, y, color.NRGBA{R: r, G: g, B: b, A: 255})
			} else {
				// If the pixel is not CMYK, copy it as-is.
				rgbImg.Set(x, y, c)
			}
		}
	}

	return rgbImg
}

// decodeImage decodes an image from a file, supporting multiple formats.
func DecodeImage(file *os.File) (image.Image, string, error) {
	// Try to decode the image format.
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", fmt.Errorf("unsupported image format: %v", err)
	}
	return img, format, nil
}

// encodeImage encodes an image to a file based on the format.
func EncodeImage(outputFile *os.File, img image.Image, format string) error {
	switch format {
	case "jpeg":
		return jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 100})
	case "png":
		return png.Encode(outputFile, img)
	case "tiff":
		return tiff.Encode(outputFile, img, nil)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}
