package signature

import (
	"strconv"
	"strings"
)

type PatternOwner struct {
	Name    string
	Pattern *PatternEdge
}

type PatternEdge struct {
	owner      *PatternOwner
	Value      uint8
	IsWildcard bool
	Next       *PatternEdge
}

func New(name string, pattern string) *PatternOwner {
	patternBytes := strings.Split(pattern, " ")

	patternOwner := &PatternOwner{
		Name:    name,
		Pattern: nil,
	}
	currentEdge := &PatternEdge{
		owner: patternOwner,
		Next:  nil,
	}
	patternOwner.Pattern = currentEdge

	for i, patternByte := range patternBytes {
		isWildcard := patternByte == "??"
		var parsedValue uint8

		if !isWildcard {
			value, _ := strconv.ParseInt(patternByte, 16, 64)
			parsedValue = uint8(value)
		}

		currentEdge.Value = parsedValue
		currentEdge.IsWildcard = isWildcard

		if i >= len(patternBytes)-1 {
			break
		}

		nextEdge := &PatternEdge{
			owner: patternOwner,
			Next:  nil,
		}
		currentEdge.Next = nextEdge
		currentEdge = nextEdge
	}

	return patternOwner
}
