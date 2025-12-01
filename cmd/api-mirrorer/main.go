package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/sourcehawk/operator-api-mirror/pkg"
)

func main() {
	var (
		operatorsFilePath = "operators.yaml"
		mirrorRootPath    = "mirrors"
		moduleRoot        = "github.com/sourcehawk/operator-api-mirror"
		mirrorTarget      = ""
	)
	flag.StringVar(&operatorsFilePath, "config", operatorsFilePath, "path to operators.yaml")
	flag.StringVar(&mirrorRootPath, "mirrorsPath", mirrorRootPath, "path to mirroring root directory")
	flag.StringVar(&moduleRoot, "rootModuleName", moduleRoot, "go module root")
	flag.StringVar(&mirrorTarget, "target", mirrorTarget, "optional target (slug) from config file")
	flag.Parse()

	if mirrorTarget != "" {
		log.Printf("Mirroring for target: %s", mirrorTarget)
	}

	operators, err := readOperatorsFile(operatorsFilePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", operators)

	err = operators.Process(mirrorRootPath, moduleRoot, mirrorTarget)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Success")
}

func readOperatorsFile(path string) (*pkg.OperatorsFile, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(root, path)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config := &pkg.OperatorsFile{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
