//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"sort"

	"go.yaml.in/yaml/v4"
)

type designator struct {
	Abbr []string `yaml:"abbr"`
	Lead string   `yaml:"lead"`
}

var baseSuffixes = []string{
	"private limited company",
	"public limited company",
	"proprietary limited",
	"limited liability company",
	"limited liability partnership",
	"limited partnership",
	"limited company",
	"no liability",
	"not for profit",
	"incorporated",
	"corporation",
	"unlimited",
	"cooperative",
	"foundation",
	"holdings",
	"holding",
	"partners",
	"limited",
	"company",
	"société",
	"unltd",
	"corp",
	"inc",
	"llc",
	"llp",
	"lp",
	"plc",
	"ltd",
	"co",
	"société en commandite",
	"private company",
	"public company",
	"pvt ltd",
	"pvt",
	"pte ltd",
	"pte",
	"sdn",
}

func main() {
	data, err := os.ReadFile("company_designator.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read YAML: %v\n", err)
		os.Exit(1)
	}

	var designators map[string]designator
	if err := yaml.Unmarshal(data, &designators); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse YAML: %v\n", err)
		os.Exit(1)
	}

	seen := make(map[string]bool)
	var suffixes []string

	for _, s := range baseSuffixes {
		if !seen[s] {
			seen[s] = true
			suffixes = append(suffixes, s)
		}
	}

	for name, d := range designators {
		if !seen[name] {
			seen[name] = true
			suffixes = append(suffixes, name)
		}
		for _, abbr := range d.Abbr {
			if !seen[abbr] {
				seen[abbr] = true
				suffixes = append(suffixes, abbr)
			}
		}
	}
	sort.Strings(suffixes)

	var buf bytes.Buffer
	buf.WriteString("package coname\n\n")
	buf.WriteString("var generatedSuffixes = []string{\n")
	for _, s := range suffixes {
		fmt.Fprintf(&buf, "\t%q,\n", s)
	}
	buf.WriteString("}\n")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to format: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("designators_gen.go", formatted, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated designators_gen.go with %d suffixes\n", len(suffixes))
}
