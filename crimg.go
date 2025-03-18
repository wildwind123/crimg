package crimg

import (
	"io"
)

type ImgFormat string

const WebpFormat ImgFormat = "webp"

type ImageInfo struct {
	Width  int
	Height int
}

type ReqCompressImage struct {
	Format        ImgFormat
	InputReader   io.Reader
	InputFilePath string
	ReturnByte    bool
	ImageResize   ImageResize
}

type ImageResize struct {
	Height int
	Width  int
}

type ReqGetImageInfo struct {
	InputReader   io.Reader
	InputFilePath string
}

type CompressedImageInfo struct {
	CompressedFilePath string
	CompressedFileByte []byte
}

type Imager interface {
	ImageInfoGetter
	ImageCompressor
}

type ImageInfoGetter interface {
	GetImageInfo(req *ReqGetImageInfo) (*ImageInfo, error)
}

type ImageCompressor interface {
	CompressImage(req *ReqCompressImage) (*CompressedImageInfo, error)
}
