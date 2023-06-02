package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
)

type Module struct {
	XMLName xml.Name `xml:"module"`
	Id      string   `xml:"id,attr"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name,attr"`
	Fat     int      `xml:"fat,attr"`
	Size    int      `xml:"size,attr"`
}

type Dependency struct {
	XMLName xml.Name `xml:"dependency"`
	From    string   `xml:"from,attr"`
	To      string   `xml:"to,attr"`
	Type    string   `xml:"type,attr"`
}

type Data struct {
	XMLName      xml.Name     `xml:"data"`
	Flavor       string       `xml:"flavor,attr"`
	Version      string       `xml:"version,attr"`
	Site         string       `xml:"origin,attr"`
	Modules      []Module     `xml:"modules>module"`
	Dependencies []Dependency `xml:"dependencies>dependency"`
}

const mode packages.LoadMode = packages.NeedName |
	packages.NeedImports

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintln(out, "usage: goflavor [options] <module dir>\n")
		fmt.Fprintln(out, "Options:")
		flag.PrintDefaults()
	}

	pattern := flag.String("pattern", "./...", "Go package pattern")
	output := flag.String("output", "go-flavor-output.xml", "Output file")
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("Expecting a single argument: directory of module")
	}

	var fset = token.NewFileSet()
	cfg := &packages.Config{Fset: fset, Mode: mode, Dir: flag.Args()[0]}
	pkgs, err := packages.Load(cfg, *pattern)
	if err != nil {
		log.Fatal(err)
	}

	var modules []Module
	var dependencies []Dependency

	for _, pkg := range pkgs {
		modules = append(modules, Module{Id: pkg.ID, Type: "package", Name: pkg.ID, Fat: 0, Size: 0})
		for pkgid := range pkg.Imports {
			dependencies = append(dependencies, Dependency{From: pkg.ID, To: pkgid, Type: "imports"})
		}
	}

	data := &Data{Flavor: "com.earldata.s101.golangflavor", Version: "0.0.1", Site: "http://github.com/earldata/go-flavor", Modules: modules, Dependencies: dependencies}
	out, _ := xml.MarshalIndent(data, "", "  ")
	os.WriteFile(*output, []byte(xml.Header+string(out)), 0644)
}
