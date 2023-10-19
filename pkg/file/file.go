package file

import (
	"github.com/HunterPie/Longinus/core/reader"
	"os"
)

type Reader struct {
	bytes []uint8
}

func (r Reader) Read() []uint8 {
	return r.bytes
}

func New(path string) reader.ByteDataSource {
	bytes, _ := os.ReadFile(path)

	return &Reader{bytes}
}
