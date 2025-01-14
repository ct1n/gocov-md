package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/axw/gocov"
)

var repo = flag.String("repo", "", "Repo to strip from package names")

func main() {
	flag.Parse()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read coverage file: %s\n", err)
		os.Exit(1)
	}
	report, err := unmarshalJson(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal coverage data: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("| Package | Lines | Coverage |")
	fmt.Println("|---------|-------|----------|")
	printCov := func(name string, reached, statements int) {
		p := percentage(reached, statements)
		fmt.Printf("| %s | %d/%d | %.2f%% |\n", name, reached, statements, p)
	}

	var totalStatements, totalReached int

	for _, pkg := range report.Packages {
		reached, statements := 0, 0
		for _, fn := range pkg.Functions {
			for _, stmt := range fn.Statements {
				statements++
				if stmt.Reached > 0 {
					reached++
				}
			}
		}
		totalStatements += statements
		totalReached += reached

		name := pkg.Name
		if *repo != "" {
			name = strings.TrimPrefix(name, *repo)
			name = strings.TrimPrefix(name, "/")
		}
		printCov(name, reached, statements)
	}
	printCov("Total", totalReached, totalStatements)
}

func unmarshalJson(data []byte) (*report, error) {
	r := &report{}
	err := json.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}
	return r, err
}

func percentage(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return 100 * float64(a) / float64(b)
}

type report struct {
	Packages []*gocov.Package
}

func (r *report) addPackage(p *gocov.Package) {
	i := sort.Search(len(r.Packages), func(i int) bool {
		return r.Packages[i].Name >= p.Name
	})
	if i < len(r.Packages) && r.Packages[i].Name == p.Name {
		r.Packages[i].Accumulate(p)
	} else {
		r.Packages = slices.Insert(r.Packages, i, p)
	}
}
