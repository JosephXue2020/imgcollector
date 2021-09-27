package imgresize

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"

	"github.com/nfnt/resize"
)

func openFile(pth string) (*os.File, string, error) {
	f, err := os.Open(pth)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	_, fname := filepath.Split(pth)
	ext := path.Ext(fname)

	return f, ext, err
}

func openImageFile(pth string) (*image.Image, error) {
	f, err := os.Open(pth)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, fname := filepath.Split(pth)
	ext := path.Ext(fname)

	if ext == ".jpg" || ext == ".jpeg" {
		img, err := jpeg.Decode(f)
		imgAddr := &img
		return imgAddr, err
	}

	if ext == ".png" {
		img, err := png.Decode(f)
		imgAddr := &img
		return imgAddr, err
	}

	if ext == ".gif" {
		img, err := gif.Decode(f)
		imgAddr := &img
		return imgAddr, err
	}

	return nil, fmt.Errorf("Unsupport image type.")
}

func getTargetSize(img image.Image, maxWidth int, maxHeight int) (int, int, bool) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	ratio := width / height
	targetRatio := maxWidth / maxHeight

	if ratio >= targetRatio {
		// 只需比较宽度
		if width > maxWidth {
			return maxWidth, 0, true
		} else {
			// 不用rescale
			return 0, 0, false
		}
	} else {
		// 只需比较高度
		if height > maxHeight {
			return 0, maxHeight, true
		} else {
			// 不用rescale
			return 0, 0, false
		}
	}
}

func ResizeImage(img *image.Image, outpath string) error {
	// 限宽400，限高300
	w, h, flag := getTargetSize(*img, 400, 300)
	var m image.Image
	if !flag {
		m = *img
	} else {
		m = resize.Resize(uint(w), uint(h), *img, resize.Lanczos3)
	}

	out, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	err = jpeg.Encode(out, m, nil)

	return err
}

func ResizeImageFile(pth string, outpath string) error {
	img, err := openImageFile(pth)
	if err != nil {
		return err
	}

	err = ResizeImage(img, outpath)
	return err
}
