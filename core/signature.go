package core

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

func NewSignatureFrom(name string, pattern string) *PatternOwner {
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

	for _, patternByte := range patternBytes {
		isWildcard := patternByte == "??"
		var parsedValue uint8

		if !isWildcard {
			value, _ := strconv.ParseInt(patternByte, 16, 64)
			parsedValue = uint8(value)
		}

		currentEdge.Value = parsedValue
		currentEdge.IsWildcard = isWildcard
		nextEdge := &PatternEdge{
			owner: patternOwner,
			Next:  nil,
		}
		currentEdge.Next = nextEdge
		currentEdge = nextEdge
	}

	return patternOwner
}
