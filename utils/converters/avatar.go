package converters

import (
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

func ResizeAvatar(fileType types.ImageFileType, file *multipart.File) ([]byte, error) {

	var (
		img image.Image
		err error
	)
	switch fileType {
	case types.IMAGE_FILE_TYPE_WEBP:
		img, err = webp.Decode(*file)

	case types.IMAGE_FILE_TYPE_JPEG:
		img, err = jpeg.Decode(*file)

	case types.IMAGE_FILE_TYPE_PNG:
		img, err = png.Decode(*file)
	}
	if err != nil {
		return nil, err
	}

	resizedImg := resize.Resize(
		consts.STANDERED_AVATAR_SIZE,
		consts.STANDERED_AVATAR_SIZE,
		img, resize.Lanczos3,
	)

	imgData, err := webpEncoder.EncodeRGBA(resizedImg, consts.AVATAR_QUALITY)
	if err != nil {
		return nil, err
	}

	_, err = (*file).Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return imgData, nil
}
