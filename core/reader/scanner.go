package reader

import (
	"encoding/binary"
	"github.com/HunterPie/Longinus/core/signature"
	"github.com/HunterPie/Longinus/core/tree"
)

const PE_OFFSET = 0xA00

func min[T int](a, b T) T {
	if a < b {
		return a
	}
	return b
}

type Scanner struct {
	reader ByteDataSource
	tree   *tree.PatternTree
}

type ScanResult struct {
	Owner  *signature.PatternOwner
	Offset int64
}

func NewResult(offset int64, owner *signature.PatternOwner) *ScanResult {
	return &ScanResult{
		Owner:  owner,
		Offset: offset,
	}
}

func (s *Scanner) Execute() []*ScanResult {
	bytes := s.reader.Read()
	found := make([]*ScanResult, 0)

	for i := range bytes {
		end := min(i+s.tree.Depth, len(bytes))

		for _, owner := range s.tree.FindPattern(bytes[i:end]) {
			var addressFound int64

			if !owner.IsRelative {
				addressFound = int64(i)
			} else {
				targetStart := i + owner.TargetOffset
				targetEnd := targetStart + 0x4
				addressFound = calculateRelativeOffset(i, owner.TargetOffset, bytes[targetStart:targetEnd])
			}

			found = append(found, NewResult(addressFound+PE_OFFSET, owner))
		}
	}

	return found
}

func calculateRelativeOffset(currentOffset, targetOffset int, targetRaw []byte) int64 {
	target := binary.LittleEndian.Uint32(targetRaw)

	return int64(currentOffset) + int64(targetOffset) + int64(target) + 0x4
}

func New(reader ByteDataSource, tree *tree.PatternTree) *Scanner {
	return &Scanner{
		reader: reader,
		tree:   tree,
	}
}
