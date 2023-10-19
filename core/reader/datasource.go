package reader

type ByteDataSource interface {
	Read() []uint8
}
