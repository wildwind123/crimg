package crimg

import (
	"image"
	_ "image/jpeg"
	"os"

	"github.com/go-faster/errors"
	_ "golang.org/x/image/webp"

	_ "image/png"

	_ "golang.org/x/image/tiff"
)

// dIG *DefaultImageGetter github.com/wildwind123/crimg.ImageInfoGetter
type DefaultImageGetter struct {
}

func (dIG *DefaultImageGetter) GetImageInfo(req *ReqGetImageInfo) (*ImageInfo, error) {
	var r = req.InputReader
	if req.InputFilePath != "" {
		reader, err := os.Open(req.InputFilePath)
		if err != nil {
			return nil, errors.Wrap(err, "cant open file")
		}
		defer reader.Close()
		r = reader
	}

	im, _, err := image.DecodeConfig(r)
	if err != nil {
		return nil, errors.Wrap(err, "cant decode config")
	}
	return &ImageInfo{
		Width:  im.Width,
		Height: im.Height,
	}, nil
}
