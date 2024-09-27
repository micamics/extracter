package models

import (
	"io"
)

type File struct {
	Reader io.Reader
}
