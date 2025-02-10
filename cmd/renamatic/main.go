package main

import (
	"flag"
	"log"

	"github.com/notJoon/renamatic/internal"
)

const defaultMappingPath = "../mapping.yml"

func main() {
	mappingPath := flag.String("mapping", defaultMappingPath, "Path to YAML file with function name mappings (default: mapping.yml)")
	dirPath := flag.String("dir", ".", "Directory to traverse (default: current directory)")
	flag.Parse()

	mapping, err := internal.LoadMapping(*mappingPath)
	if err != nil {
		log.Fatalf("Failed to load mapping file: %v", err)
	}

	if err := internal.ProcessDir(*dirPath, mapping); err != nil {
		log.Fatalf("Failed to process directory: %v", err)
	}
}
