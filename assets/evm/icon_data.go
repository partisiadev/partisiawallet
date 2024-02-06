package evm

import (
	"bytes"
	"errors"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
	"strings"
)

type IconData struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format"`
}

func (i IconData) SmallEncodedImage() ([]byte, error) {
	smImg, err := i.SmallImage()
	if smImg != nil && err == nil {
		var b bytes.Buffer
		if i.Format == "png" {
			err = png.Encode(&b, smImg)
		} else {
			err = jpeg.Encode(&b, smImg, &jpeg.Options{Quality: 75})
		}
		if err == nil && len(b.Bytes()) > 0 {
			return b.Bytes(), nil
		}
	}
	if err == nil {
		err = ErrUnsupportedImage
	}
	return nil, err
}

func (i IconData) SmallImage() (image.Image, error) {
	img, err := i.OrigImage()
	if img != nil && err == nil {
		if img != nil {
			dst := image.NewRGBA(image.Rectangle{
				Max: image.Point{X: 256, Y: 256},
			})
			draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
			return dst, nil
		}
	}
	return nil, ErrUnsupportedImage
}

func (i IconData) OrigImage() (image.Image, error) {
	format := i.Format
	isSupported := format == "jpg" || format == "png" || format == "jpeg"
	filename, ok := strings.CutPrefix(i.URL, "ipfs://")
	if isSupported && ok {
		var img image.Image
		imgFs, err := IconsDownloadDir.Open(IconsDownloadDirName + "/" + filename)
		if err == nil {
			if format == "jpg" || format == "jpeg" {
				img, _ = jpeg.Decode(imgFs)
			} else {
				img, _ = png.Decode(imgFs)
			}
			_ = imgFs.Close()
			return img, nil
		}
	}
	return nil, ErrUnsupportedImage
}

func (i IconData) OrigEncodedImage() ([]byte, error) {
	format := i.Format
	isSupported := format == "jpg" || format == "png" || format == "jpeg"
	filename, ok := strings.CutPrefix(i.URL, "ipfs://")
	if isSupported && ok {
		encBs, err := IconsDownloadDir.ReadFile(IconsDownloadDirName + "/" + filename)
		if err != nil {
			return nil, err
		}
		return encBs, err
	}
	return nil, errors.New("unsupported image")
}
func (i IconData) Encoded() ([]byte, error) {
	filename, ok := strings.CutPrefix(i.URL, "ipfs://")
	if ok {
		return IconsDownloadDir.ReadFile(IconsDownloadDirName + "/" + filename)
	}
	return nil, ErrUnsupportedImage
}
