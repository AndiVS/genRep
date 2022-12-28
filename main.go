package main

import (
	"flag"
	"fmt"
	"github.com/AndiVS/genRep/internal/generators"
	"github.com/AndiVS/genRep/internal/helper"
	"github.com/AndiVS/genRep/internal/model"
	"github.com/AndiVS/genRep/internal/parser"
	validator "github.com/AndiVS/genRep/internal/validator"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
)

var (
	typeNames  = flag.String("type", "", "comma-separated list of type names; must be set")
	tableNames = flag.String("table", "type name in snake case", "comma-separated list of table names")
	schemes    = flag.String("schema", "public", "comma-separated list of schema")
	output     = flag.String("output", ".", "output path;")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of repBuilder:\n")
	fmt.Fprintf(os.Stderr, "\trepBuilder -type=TypeName -table=TableName -schema=dbSchema -output=outputDir\n")
	fmt.Fprintf(os.Stderr, "\tonli type mandatory\n")
	fmt.Fprintf(os.Stderr, "\tif table is specified, the number of tables must be equal to the number of types\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("repBuilder: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	types := strings.Split(*typeNames, ",")
	tables := strings.Split(*tableNames, ",")
	schemes := strings.Split(*schemes, ",")

	if tables[0] != "type name in snake case" && len(tables) != len(types) {
		flag.Usage()
		os.Exit(2)
	}
	if schemes[0] != "public" && len(schemes) != len(types) {
		flag.Usage()
		os.Exit(2)
	}

	models := make([]*model.Model, len(types))
	for i := range types {
		models[i] = &model.Model{
			Name: &types[i],
		}
		if len(tables) == 1 {
			buf := helper.ToSnakeCase(*models[i].Name)
			models[i].TableName = &buf
		} else {
			models[i].TableName = &tables[i]
		}
		if len(schemes) == 1 {
			models[i].Schema = &schemes[0]
		} else {
			models[i].TableName = &schemes[i]
		}
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	files := parser.ParsePackage(args)

	mod := parser.ParseGoStructToModel(files, models)
	if mod[0].Fields == nil {
		logrus.Fatal("zero fields")
	}

	err := validator.Validate(mod)
	if err != nil {
		logrus.Fatal(err)
	}

	err = generators.Generate(mod, *output)
	if err != nil {
		logrus.Fatal(err)
	}

	err = generators.GeneratePagination(*output)
	if err != nil {
		logrus.Fatal(err)
	}
}
