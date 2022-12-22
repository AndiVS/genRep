package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/AndiVS/genRep/internal/helper"
	"github.com/AndiVS/genRep/internal/model"
)

func ParseGoStructToModel(targetFile, targetStruct string) (*model.Model, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, targetFile, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("error while open parse file : %w", err)
	}

	mod := &model.Model{
		Name: &targetStruct,
	}
	for _, v := range file.Decls {
		if v.(*ast.GenDecl).Tok == token.TYPE {
			k := v.(*ast.GenDecl).Specs[0].(*ast.TypeSpec)
			if k.Name.Name == targetStruct {
				for _, j := range k.Type.(*ast.StructType).Fields.List {
					f := &model.Field{
						Name: &j.Names[0].Name,
						Type: GetType(j.Type),
					}

					tgMap := GetTags(j.Tag)
					if sqlName, ok := tgMap["sqlName"]; ok {
						f.SqlName = sqlName
					} else {
						buf := helper.ToSnakeCase(j.Names[0].Name)
						f.SqlName = &buf
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

	return mod, err
}

func GetType(exp ast.Expr) *string {
	var typeName *string
	switch exp.(type) {
	case *ast.SelectorExpr:
		tempStr := exp.(*ast.SelectorExpr).X.(*ast.Ident).Name + "." + exp.(*ast.SelectorExpr).Sel.Name
		typeName = &tempStr
	case *ast.Ident:
		typeName = &exp.(*ast.Ident).Name
	}

	return typeName
}

func GetTags(tags *ast.BasicLit) map[string]*string {
	tagMap := map[string]*string{}
	if tags != nil {
		str := tags.Value
		str = strings.Replace(str, "\"", " ", -1)
		str = strings.Replace(str, "`", "", -1)
		str = strings.Replace(str, "  ", " ", -1)
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

//func ToSQLType(str string) *string {
//	rez := ""
//	switch str {
//	case "int":
//	case "int32":
//	case "int64":
//		rez = "int"
//		break
//	case "float32":
//	case "float64":
//		rez = "float"
//		break
//	case "uuid":
//		rez = "uuid"
//		break
//	case "string":
//		rez = "varchar(60)"
//		break
//	case "bool":
//		rez = "boolean"
//		break
//	}
//
//	return &rez
//}