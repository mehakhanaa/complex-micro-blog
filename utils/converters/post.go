package converters

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"

	"github.com/KononK/resize"
	webpEncoder "github.com/chai2010/webp"
	"golang.org/x/image/webp"

	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/types"
)

func ResizePostImage(fileType types.ImageFileType, file *multipart.File) ([]byte, error) {

	var (
		img       image.Image
		imgConfig image.Config
		err       error
	)
	switch fileType {
	case types.IMAGE_FILE_TYPE_WEBP:
		img, err = webp.Decode(*file)
		if err != nil {
			return nil, err
		}

		_, err = (*file).Seek(0, 0)
		if err != nil {
			return nil, err
		}
		imgConfig, err = webp.DecodeConfig(*file)

	case types.IMAGE_FILE_TYPE_JPEG:
		img, err = jpeg.Decode(*file)
		if err != nil {
			return nil, err
		}

		_, err = (*file).Seek(0, 0)
		if err != nil {
			return nil, err
		}
		imgConfig, err = jpeg.DecodeConfig(*file)

	case types.IMAGE_FILE_TYPE_PNG:
		img, err = png.Decode(*file)
		if err != nil {
			return nil, err
		}

		_, err = (*file).Seek(0, 0)
		if err != nil {
			return nil, err
		}
		imgConfig, err = png.DecodeConfig(*file)
	}
	if err != nil {
		return nil, err
	}

	var resizedImg image.Image
	if imgConfig.Width > imgConfig.Height && imgConfig.Width > consts.POST_IMAGE_WIDTH_THRESHOLD {
		resizedImg = resize.Resize(
			consts.POST_IMAGE_WIDTH_THRESHOLD,
			uint(imgConfig.Height)*consts.POST_IMAGE_WIDTH_THRESHOLD/uint(imgConfig.Width),
			img,
			resize.Lanczos3,
		)
	} else if imgConfig.Width <= imgConfig.Height && imgConfig.Height > consts.POST_IMAGE_HEIGHT_THRESHOLD {
		resizedImg = resize.Resize(
			uint(imgConfig.Width)*consts.POST_IMAGE_HEIGHT_THRESHOLD/uint(imgConfig.Height),
			consts.POST_IMAGE_HEIGHT_THRESHOLD,
			img,
			resize.Lanczos3,
		)
	} else {
		resizedImg = img
	}

	if resizedImg == nil {
		return nil, errors.New("resize image failed")
	}

	imgData, err := webpEncoder.EncodeRGBA(resizedImg, consts.POST_IMAGE_QUALITY)
	if err != nil {
		return nil, err
	}

	_, err = (*file).Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return imgData, nil
}
