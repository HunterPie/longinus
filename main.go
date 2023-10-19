package main

import (
	"fmt"
	"github.com/HunterPie/Longinus/cli"
	"github.com/HunterPie/Longinus/core/reader"
	"github.com/HunterPie/Longinus/core/signature"
	"github.com/HunterPie/Longinus/core/tree"
	"github.com/HunterPie/Longinus/pkg/file"
	"github.com/HunterPie/Longinus/pkg/yaml"
	"github.com/jedib0t/go-pretty/v6/table"
	"log"
	"os"
)

func printResults(results []*reader.ScanResult) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Address"})

	for _, result := range results {
		t.AppendRow(table.Row{
			result.Owner.Name, fmt.Sprintf("0x%08X", result.Offset),
		})
	}
	t.Render()
}

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
		patterns = append(patterns, signature.New(pattern.Name, pattern.Signature, pattern.InstructionOffset, pattern.IsRelative))
	}

	log.Printf("Loaded %d patterns", len(patterns))

	patternTree := tree.New(patterns...)

	dataSource := file.New(args.Executable)
	log.Printf("Finished loading up datasource")

	scanner := reader.New(dataSource, patternTree)
	results := scanner.Execute()

	log.Printf("Found %d patterns", len(results))

	printResults(results)
}
