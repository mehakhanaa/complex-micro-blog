package types

type ImageFileType int

const (
	IMAGE_FILE_TYPE_UNKNOWN ImageFileType = -1

	IMAGE_FILE_TYPE_WEBP ImageFileType = iota

	IMAGE_FILE_TYPE_JPEG

	IMAGE_FILE_TYPE_PNG
)
