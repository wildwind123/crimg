package crimg

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"

	_ "golang.org/x/image/webp"
)

func TestImageInfoGetter(t *testing.T) {
	var imageInfoGetter ImageInfoGetter = &DefaultImageGetter{}

	for _, v := range []struct {
		FilePath     string
		ResultWidth  int
		ResultHeight int
	}{
		{
			FilePath:     "fixture/4.webp",
			ResultWidth:  1024,
			ResultHeight: 772,
		},
		{
			FilePath:     "fixture/example_1.png",
			ResultWidth:  912,
			ResultHeight: 513,
		},
		{
			FilePath:     "fixture/example_2.jpg",
			ResultWidth:  2048,
			ResultHeight: 1365,
		},
	} {
		reader, err := os.Open(v.FilePath)
		if err != nil {
			t.Error(err)
			return
		}
		defer reader.Close()
		// test reader
		imgInfo, err := imageInfoGetter.GetImageInfo(&ReqGetImageInfo{
			InputReader: reader,
		})
		if err != nil {
			t.Error(err)
			return
		}
		if imgInfo.Height != v.ResultHeight {
			t.Error("wrong img height")
		}
		if imgInfo.Width != v.ResultWidth {
			t.Error("wrong img width")
		}

		// test file path
		imgInfo, err = imageInfoGetter.GetImageInfo(&ReqGetImageInfo{
			InputFilePath: v.FilePath,
		})
		if err != nil {
			t.Error(err)
			return
		}
		if imgInfo.Height != v.ResultHeight {
			t.Error("wrong img height")
		}
		if imgInfo.Width != v.ResultWidth {
			t.Error("wrong img width")
		}
	}
}

func TestCWebpCompressor(t *testing.T) {
	t.Skip("manual test")
	var c CWebpCompressor = CWebpCompressor{
		CWebpBinaryPath: "/home/ganbatte/apps/bins/libwebp-1.5.0-linux-x86-64/bin/cwebp",
		ImageInfoGetter: &DefaultImageGetter{},
	}

	compressedFilePath, err := c.CompressImage(&ReqCompressImage{
		Format:        WebpFormat,
		InputFilePath: "fixture/example_2.jpg",
		ImageResize: ImageResize{
			Height: 500,
			Width:  500,
		},
		ReturnByte: false,
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("path", compressedFilePath.CompressedFilePath)
}

func TestCWebpCompressorReader(t *testing.T) {
	// t.Skip("manual test")
	var c CWebpCompressor = CWebpCompressor{
		CWebpBinaryPath: "/home/ganbatte/apps/bins/libwebp-1.5.0-linux-x86-64/bin/cwebp",
		ImageInfoGetter: &DefaultImageGetter{},
	}

	rr, err := os.Open("fixture/example_2.jpg")
	if err != nil {
		t.Error(err)
		return
	}

	compressedFilePath, err := c.CompressImage(&ReqCompressImage{
		Format:      WebpFormat,
		InputReader: rr,
		ImageResize: ImageResize{
			Height: 500,
			Width:  500,
		},
		ReturnByte: false,
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("path", compressedFilePath.CompressedFilePath)
}
