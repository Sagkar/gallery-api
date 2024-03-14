package model

import "mime/multipart"

type PhotoEntity struct {
	File *multipart.FileHeader
}
