package crimg

import (
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/go-faster/errors"
)

// https://developers.google.com/speed/webp/download
// c *CWebpCompressor github.com/wildwind123/crimg.ImageCompressor
type CWebpCompressor struct {
	CWebpBinaryPath string
	ImageInfoGetter ImageInfoGetter
}

func (c *CWebpCompressor) CompressImage(req *ReqCompressImage) (*CompressedImageInfo, error) {
	if req.Format != WebpFormat {
		return nil, errors.Errorf("support only webp")
	}

	// set fileForCompress
	fileForCompress := req.InputFilePath
	if req.InputFilePath != "" {
		if !path.IsAbs(fileForCompress) {
			ex, err := os.Executable()
			if err != nil {
				return nil, errors.Wrap(err, "cant os.Executable()")
			}

			exPath := filepath.Dir(ex)
			fileForCompress = path.Join(exPath, req.InputFilePath)
		}
	}
	// cwebp -resize 500 0
	argWidth := 0
	argHeight := 0
	commandResizeArgs := []string{}
	if req.ImageResize.Height != 0 && req.ImageResize.Width != 0 {
		iInfo, err := c.ImageInfoGetter.GetImageInfo(&ReqGetImageInfo{
			InputFilePath: fileForCompress,
		})
		if err != nil {
			return nil, errors.Wrap(err, "cant GetImageInfo")
		}
		argWidth, argHeight = calculate(iInfo.Width, iInfo.Height, req.ImageResize.Height, req.ImageResize.Width)

		commandResizeArgs = append(commandResizeArgs, "-resize")
		commandResizeArgs = append(commandResizeArgs, fmt.Sprintf("%d", argWidth))
		commandResizeArgs = append(commandResizeArgs, fmt.Sprintf("%d", argHeight))
	}
	// cwebp -resize 500 0

	// set outFilePath
	outFilePath := path.Join(os.TempDir(), fmt.Sprintf("%d_%s.webp", time.Now().UnixMilli(), path.Base(fileForCompress)))

	// build command line args
	args := []string{fileForCompress, "-o", outFilePath}
	if len(commandResizeArgs) > 0 {
		args = append(args, commandResizeArgs...)
	}

	// cwebp input.png -o output.webp
	cmd := exec.Command(c.CWebpBinaryPath, args...)

	// Capture output
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "cant run cli command")
	}
	_ = output

	r := &CompressedImageInfo{}
	if req.ReturnByte {
		reader, err := os.Open(outFilePath)
		if err != nil {
			return nil, errors.Wrap(err, "cant open file")
		}
		b, err := io.ReadAll(reader)
		if err != nil {
			return nil, errors.Wrap(err, "cant io.ReadAll")
		}
		r.CompressedFileByte = b
		os.Remove(outFilePath)
	} else {
		r.CompressedFilePath = outFilePath
	}

	return r, nil
}

func calculate(sourceWidth, sourceHeight, maxWidth, maxHeight int) (newWidth, newHeight int) {
	// Convert to float64 for precise ratio calculation
	srcW := float64(sourceWidth)
	srcH := float64(sourceHeight)
	maxW := float64(maxWidth)
	maxH := float64(maxHeight)

	// Calculate aspect ratio of source image
	aspectRatio := srcW / srcH

	// First try scaling based on width
	newW := maxW
	newH := maxW / aspectRatio

	// If height exceeds maxHeight, scale based on height instead
	if newH > maxH {
		newH = maxH
		newW = maxH * aspectRatio
	}

	if newW > newH {
		return int(math.Round(newW)), 0
	} else {
		return 0, int(math.Round(newH))
	}
}
