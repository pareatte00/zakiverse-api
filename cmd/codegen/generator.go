package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type GeneratorConfig struct {
	OutputFile string
	Template   *template.Template
	Data       any
}

func Generate(cfg GeneratorConfig) {
	var buf bytes.Buffer
	if err := cfg.Template.Execute(&buf, cfg.Data); err != nil {
		log.Fatalf("template: %v", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("format generated code: %v\n%s", err, buf.String())
	}

	if err := os.WriteFile(cfg.OutputFile, formatted, 0644); err != nil {
		log.Fatalf("write %s: %v", cfg.OutputFile, err)
	}

	fmt.Printf("codegen: generated %s\n", cfg.OutputFile)
}

func ReadModuleName() string {
	for _, candidate := range []string{"go.mod", "../go.mod", "../../go.mod", "../../../go.mod"} {
		data, err := os.ReadFile(candidate)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			if after, ok := strings.CutPrefix(line, "module "); ok {
				return strings.TrimSpace(after)
			}
		}
	}
	log.Fatal("cannot find go.mod to determine module name")
	return ""
}

func ResolveDir(subpath string) string {
	if matches, _ := filepath.Glob("*.go"); len(matches) > 0 {
		return "."
	}
	if matches, _ := filepath.Glob(subpath + "/*.go"); len(matches) > 0 {
		return subpath
	}
	log.Fatalf("cannot find .go files; run from project root or %s/ directory", subpath)
	return ""
}
