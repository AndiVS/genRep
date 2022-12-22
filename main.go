package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/AndiVS/genRep/internal/generators"
	"github.com/AndiVS/genRep/internal/parser"
	"github.com/AndiVS/genRep/internal/validator"
)

var (
	typeName  = flag.String("type", "", "type name; must be set")
	tableName = flag.String("table", "typeName", "table name")
	schema    = flag.String("schema", "public", "schema of database")
	output    = flag.String("output", "../repository", "output path;")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of repBuilder:\n")
	fmt.Fprintf(os.Stderr, "\trepBuilder -type=TypeName -table=TableName -schema=dbSchema -output=outputDir\n")
	fmt.Fprintf(os.Stderr, "\tonli type mandatory\n")
	flag.PrintDefaults()
}

func main() {

	log.SetFlags(0)
	log.SetPrefix("repBuilder: ")
	flag.Usage = Usage
	flag.Parse()
	if *typeName == "" {
		flag.Usage()
		os.Exit(2)
	}
	if *tableName == "typeName" {
		tableName = typeName
	}

	mod, err := parser.ParseGoStructToModel("../test/test_str.go", *typeName)
	if err != nil {
		log.Fatal(err)
	}

	mod.TableName = tableName
	mod.Schema = schema

	err = validator.Validate(mod)
	if err != nil {
		logrus.Fatal(err)
	}

	err = generators.Generate(mod, "../testRepo")
	if err != nil {
		logrus.Fatal(err)
	}

	err = generators.GeneratePagination(mod, "../testRepo")
	if err != nil {
		logrus.Fatal(err)
	}
}
