package templates

import (
	"fmt"
	"text/template"

	"github.com/AndiVS/genRep/internal/model"
)

var PostgresCreateTableTemplate = template.Must(template.New("").Parse(
	`CREATE TABLE IF NOT EXISTS {{ .Model.TableName}}( {{range .Model.PrimaryFields}} 
		{{ .SQLName}} {{ .SQLType}}, {{end}} {{range .Model.Fields}}
		{{ .SQLName}} {{ .SQLType}}, {{end}}
	    PRIMARY KEY ({{ .MultiPrimaryKey}})
	);
`))

func PrepareMulticolumnPrimaryKey(primaryKeys []*model.Field) string {
	var primaryKey string
	for i := 0; i < len(primaryKeys)-1; i++ {
		primaryKey += fmt.Sprint(*primaryKeys[i].SQLName + ", ")
	}
	primaryKey += fmt.Sprint(*primaryKeys[len(primaryKeys)-1].SQLName)
	return primaryKey
}
