package parser

import (
	"go/ast"
	"log"

	"golang.org/x/tools/go/packages"
)

// ParsePackage function that parse package to files
func ParsePackage(patterns []string) ([]*ast.File, string) {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, patterns...)

	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	return pkgs[0].Syntax, pkgs[0].PkgPath
}
