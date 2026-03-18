package codegen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Entry struct {
	Name     string
	Tags     map[string]string
	Messages []LocaleMessage
}

type LocaleMessage struct {
	Key   string
	Value string
}

type ParserConfig struct {
	Dir          string
	TypeName     string
	ExcludeFiles map[string]bool
	TagParsers   []TagParser
}

type TagParser struct {
	Pattern *regexp.Regexp
	Handler func(entry *Entry, matches []string)
}

func DiscoverFiles(dir string, excludeFiles map[string]bool) []string {
	pattern := filepath.Join(dir, "*.go")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("glob %s: %v", pattern, err)
	}

	var files []string
	for _, f := range matches {
		base := filepath.Base(f)
		if !excludeFiles[base] {
			files = append(files, f)
		}
	}
	sort.Strings(files)
	return files
}

func ParseFile(path string, cfg ParserConfig) []Entry {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("parse %s: %v", path, err)
	}

	var entries []Entry

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		for _, spec := range genDecl.Specs {
			vspec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			ident, ok := vspec.Type.(*ast.Ident)
			if !ok || ident.Name != cfg.TypeName {
				continue
			}

			if len(vspec.Names) != 1 {
				continue
			}
			name := vspec.Names[0].Name

			doc := vspec.Doc
			if doc == nil {
				continue
			}

			entry := Entry{
				Name: name,
				Tags: map[string]string{},
			}
			msgs := map[string]string{}

			for _, c := range doc.List {
				text := strings.TrimSpace(c.Text)

				for _, tp := range cfg.TagParsers {
					if m := tp.Pattern.FindStringSubmatch(text); m != nil {
						tp.Handler(&entry, m)
						break
					}
				}

				if m := reLocale.FindStringSubmatch(text); m != nil {
					msgs[m[1]] = m[2]
				}
			}

			if len(entry.Tags) == 0 && len(msgs) == 0 {
				continue
			}

			keys := make([]string, 0, len(msgs))
			for k := range msgs {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				entry.Messages = append(entry.Messages, LocaleMessage{Key: k, Value: msgs[k]})
			}

			entries = append(entries, entry)
		}
	}

	return entries
}

var reLocale = regexp.MustCompile(`^//\s*@Locale\s+(\w+)\s+"([^"]*)"`)

func ParseFiles(files []string, cfg ParserConfig) []Entry {
	var entries []Entry
	seen := map[string]string{}

	for _, f := range files {
		parsed := ParseFile(f, cfg)
		for _, e := range parsed {
			if prev, dup := seen[e.Name]; dup {
				log.Fatalf("duplicate %s %q: first in %s, again in %s", cfg.TypeName, e.Name, prev, f)
			}
			seen[e.Name] = f
		}
		entries = append(entries, parsed...)
	}

	return entries
}
