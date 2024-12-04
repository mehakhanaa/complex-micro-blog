package validers

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"

	"golang.org/x/image/webp"

	"github.com/mehakhanaa/complex-micro-blog/types"
)

func ValidImageFile(fileHeader *multipart.FileHeader, file *multipart.File, minWidth, minHeight int, maxSize int64) (types.ImageFileType, error) {
	fileType := types.IMAGE_FILE_TYPE_UNKNOWN

	if fileHeader.Size > maxSize {
		return fileType, errors.New("image file size too large")
	}

	var (
		imgConfig image.Config
		err       error
	)

	switch fileHeader.Header.Get("Content-Type") {
	case "image/jpeg":
		fileType = types.IMAGE_FILE_TYPE_JPEG
		imgConfig, err = jpeg.DecodeConfig(*file)
	case "image/png":
		fileType = types.IMAGE_FILE_TYPE_PNG
		imgConfig, err = png.DecodeConfig(*file)
	case "image/webp":
		fileType = types.IMAGE_FILE_TYPE_WEBP
		imgConfig, err = webp.DecodeConfig(*file)
	default:
		return fileType, errors.New("image file type not supported")
	}
	if err != nil {
		return fileType, err
	}

	if imgConfig.Width < minWidth || imgConfig.Height < minHeight {
		return fileType, errors.New("image size too small")
	}

	_, err = (*file).Seek(0, 0)
	if err != nil {
		return fileType, err
	}

	return fileType, nil
}
