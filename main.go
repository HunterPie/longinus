package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/HunterPie/Longinus/cli"
	"github.com/HunterPie/Longinus/core/signature"
	"github.com/HunterPie/Longinus/core/tree"
	"github.com/HunterPie/Longinus/pkg/yaml"
	"os"
)

func main() {
	args, err := cli.ParseArguments()

	if err != nil {
		fmt.Println(err)
		return
	}

	config, err := yaml.LoadConfig(args.Configuration)

	if err != nil {
		fmt.Println(err)
		return
	}

	patterns := make([]*signature.PatternOwner, 0)

	for _, pattern := range config.Executables[0].Signatures {
		patterns = append(patterns, signature.New(pattern.Name, pattern.Signature))
	}

	patternGraph := tree.New(patterns...)

	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(patternGraph)
	os.WriteFile("output.json", buffer.Bytes(), 0644)
}
