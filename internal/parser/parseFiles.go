// Package parser provides functions for parsing input files
package parser

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/AndiVS/genRep/internal/helper"
	"github.com/AndiVS/genRep/internal/model"
)

// ParseGoStructToModel function thant pars ast files to model for generating
func ParseGoStructToModel(files []*ast.File, models []*model.Model) []*model.Model {
	for _, file := range files {
		for _, v := range file.Decls {
			v, ok := v.(*ast.GenDecl)
			if ok && v.Tok == token.TYPE {
				k := v.Specs[0].(*ast.TypeSpec)
				for _, mod := range models {
					if k.Name.Name == *mod.Name {
						for _, j := range k.Type.(*ast.StructType).Fields.List {
							f := &model.Field{
								Name: &j.Names[0].Name,
								Type: getType(j.Type),
							}

							f.SQLType = toSQLType(*f.Type)

							tgMap := getTags(j.Tag)
							if sqlName, ok := tgMap["sqlName"]; ok {
								f.SQLName = sqlName
							} else {
								buf := helper.ToSnakeCase(j.Names[0].Name)
								f.SQLName = &buf
							}

							if _, ok := tgMap["primary"]; ok {
								mod.PrimaryFields = append(mod.PrimaryFields, f)
							} else {
								mod.Fields = append(mod.Fields, f)
							}
						}
					}
				}
			}
		}
	}

	return models
}

func getType(exp ast.Expr) *string {
	var typeName *string
	switch exp := exp.(type) {
	case *ast.SelectorExpr:
		tempStr := exp.X.(*ast.Ident).Name + "." + exp.Sel.Name
		typeName = &tempStr
	case *ast.Ident:
		typeName = &exp.Name
	}

	return typeName
}

func getTags(tags *ast.BasicLit) map[string]*string {
	tagMap := map[string]*string{}
	if tags != nil {
		str := tags.Value
		str = strings.ReplaceAll(str, "\"", " ")
		str = strings.ReplaceAll(str, "`", "")
		str = strings.ReplaceAll(str, "  ", " ")
		str = strings.TrimSpace(str)
		arr := strings.Split(str, " ")

		for i := 0; i < len(arr); i++ {
			buf := strings.Split(arr[i], ":")
			if len(buf) > 1 {
				tagMap[buf[0]] = &buf[1]
			} else {
				tagMap[buf[0]] = nil
			}
		}
	}
	return tagMap
}

func toSQLType(str string) *string {
	rez := ""
	switch str {
	case "int":
	case "int32":
	case "int64":
		rez = "int"
	case "float32":
	case "float64":
		rez = "float"
	case "uuid.UUID":
		rez = "uuid"
	case "string":
		rez = "varchar(60)"
	case "bool":
		rez = "boolean"
	}

	return &rez
}
