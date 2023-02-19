package main

import (
	"bytes"
	"encoding/json"
	"github.com/HunterPie/Longinus/core"
	"os"
)

type PatternDefinition struct {
	Name    string
	Pattern string
}

func main() {
	definitions := []PatternDefinition{
		{
			Name:    "STAGE_ADDRESS",
			Pattern: "48 8B 05 ?? ?? ?? ?? ?? 8B ?? ?? ?? ?? 83 78 60 ?? 74 ??",
		},
		{
			Name:    "MONSTER_ADDRESS",
			Pattern: "48 8B 05 ?? ?? ?? ?? 48 8B 88 ?? ?? ?? ?? 4A 8B 44 C9 20",
		},
		{
			Name:    "CRC_FUNC_1",
			Pattern: "48 89 5C 24 08 48 89 6C 24 10 48 89 74 24 18 48 89 7C 24 20 41 56 48 83 EC 20 44 8B 05 ?? ?? ?? ?? 48 8B EA 8B 05 ?? ?? ?? ??",
		},
		{
			Name:    "CRC_FUNC_2",
			Pattern: "48 89 5C 24 08 48 89 6C 24 10 48 89 74 24 18 57 48 83 EC 20 48 8B 5A 10 48 8B F2 48 8B 52 20 48 8B E9",
		},
		{
			Name:    "CRC_FUNC_3",
			Pattern: "48 89 5C 24 08 48 89 6C 24 10 48 89 74 24 18 57 48 83 EC 20 4C 8B 05 ?? ?? ?? ?? 33 DB",
		},
	}

	patterns := make([]*core.PatternOwner, 0)

	for _, definition := range definitions {
		patterns = append(patterns, core.NewSignatureFrom(definition.Name, definition.Pattern))
	}

	graph := core.NewGraph(patterns...)

	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(graph)
	os.WriteFile("output.json", buffer.Bytes(), 0644)
}
