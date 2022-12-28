package templates

import (
	"strconv"
	"text/template"

	"github.com/AndiVS/genRep/internal/model"
)

var PostgresRepositoryTemplate = template.Must(template.New("").Parse(
	`// Package repository applies for interaction with database
package repository

import(
	"context"{{ if .TimeFieldExists }}
	"time"{{ end }}
	"fmt"{{ if .GetWithSortAndPagination }}
{{ end }}
	{{ if .UUIDFieldExists}}
	"github.com/google/uuid" {{ end}}
	"github.com/jackc/pgx/v4/pgxpool"
)

{{ $tick := "` + "`" + `" }}
// {{ .Model.Name }}RepositoryManager is interface with methods to interact with database
type {{ .Model.Name }}RepositoryManager interface { {{ if .CreateMethod }}
	Create(ctx context.Context, obj *model.{{ .Model.Name }}) error {{ end }} {{ if .GetByIDMethod }}
	GetByID(ctx context.Context, {{ .PrimaryKeys}}) (*model.{{ .Model.Name }}, error) {{ end }} {{ if .GetAllMethod }}
	GetAll(ctx context.Context) ([]model.{{ .Model.Name }}, error) {{ end }} {{ if .UpdateMethod}}
	Update(ctx context.Context, obj *model.{{ .Model.Name }}) error {{ end }} {{ if .DeleteMethod }}
	Delete(ctx context.Context, {{ .PrimaryKeys}}) error {{ end}} {{ if .GetWithSortAndPagination }}
	GetWithSortAndPagination(ctx context.Context, p *pagination.Pagination[*model.{{ .Model.Name }}]) ([]*model.{{ .Model.Name }}, int, error) {{ end }}
}

// {{ .ModelNameLower}}Repository is {{ .Model.Name}}RepositoryManager implementation
type {{ .ModelNameLower}}Repository struct {
	pool *pgxpool.Pool
} 

// New{{ .Model.Name}}RepositoryManager returns {{ .Model.Name }}RepositoryManager instance
func New{{ .Model.Name}}RepositoryManager(pool *pgxpool.Pool) {{ .Model.Name}}RepositoryManager {
	return &{{ .ModelNameLower}}Repository{pool: pool}
} 

{{ if .CreateMethod }}// Create method insert {{ .ModelNameLower}} record into database
func (rps *{{ .ModelNameLower}}Repository) Create(ctx context.Context, obj *model.{{ .Model.Name }}) error {
	_, err := rps.pool.Exec(ctx, 
		{{ $tick }}{{ .SqlCreate }}{{ $tick }}, 
		{{ .CreateValues}})
	if err != nil {
		return fmt.Errorf("repository: can't create {{ .ModelNameLower}} record - %s", err)
	}
	return nil
}{{ end}}

{{ if .GetByIDMethod}}// GetByID method returns {{ .ModelNameLower}} record with selection by id
func (rps *{{ .ModelNameLower}}Repository) GetByID(ctx context.Context, {{ .PrimaryKeys}}) (*model.{{ .Model.Name }}, error) {
	var obj model.{{ .Model.Name }}
	err := rps.pool.QueryRow(ctx,
		{{ $tick }}{{ .SqlGetByID}}{{ $tick }}, {{ .PrimaryValues}}).Scan( {{range .Model.PrimaryFields}}
		&obj.{{ .Name}}, {{ end }} {{ range .Model.Fields}}
		&obj.{{ .Name}}, {{ end }}
	)
	if err != nil {
		return nil, fmt.Errorf("repository: can't get {{ .ModelNameLower }} record - %s", err)
	}
	return &obj, nil
}{{ end }}

{{ if .GetAllMethod}}// GetAll method return all {{ .ModelNameLower }} records from the database
func (rps *{{ .ModelNameLower}}Repository) GetAll(ctx context.Context) ([]model.{{ .Model.Name }}, error) {
	var objs []model.{{ .Model.Name }}
	rows, err := rps.pool.Query(ctx,
		{{ $tick }}{{ .SqlGetAll}}{{ $tick }})
	if err != nil {
		return nil, fmt.Errorf("repository: can't get all {{ .ModelNameLower }} records - %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var obj model.{{ .Model.Name }}
		err = rows.Scan({{range .Model.PrimaryFields}}
			&obj.{{ .Name }},{{ end }} {{ range .Model.Fields}}
			&obj.{{ .Name}}, {{ end }}
		)
		if err != nil {
			return nil, fmt.Errorf("repository: can't parse {{ .ModelNameLower }} record - %s", err)
		}
		objs = append(objs, obj)
	}
	return objs, nil
}{{ end }}

{{ if .UpdateMethod}}// Update method updates {{ .ModelNameLower}} record
func (rps *{{ .ModelNameLower}}Repository) Update(ctx context.Context, obj *model.{{ .Model.Name }}) error {
	_, err := rps.pool.Exec(ctx,
		{{ $tick }}{{ .SqlUpdate}}{{ $tick }},
		{{ .UpdateValues}})
	if err != nil {
		return fmt.Errorf("repository: can't update {{ .ModelNameLower }} record - %s", err)
	}
	return nil
}{{ end }}

{{ if .DeleteMethod}}// Delete method removes {{ .ModelNameLower}} record
func (rps *{{ .ModelNameLower}}Repository) Delete(ctx context.Context, {{ .PrimaryKeys}}) error {
	_, err := rps.pool.Exec(ctx, 
		{{ $tick }}{{ .SqlDelete }}{{ $tick }},
		{{ .PrimaryValues}})
	if err != nil {
		return fmt.Errorf("repository: can't delete {{ .ModelNameLower }} record - %s", err)
	}
	return nil
}{{ end }}

{{ if .GetWithSortAndPagination}}// GetWithSortAndPagination method return all {{ .ModelNameLower }} records from the database
func (rps *{{ .ModelNameLower}}Repository) GetWithSortAndPagination(ctx context.Context, p *pagination.Pagination[*model.{{ .Model.Name }}]) ([]*model.{{ .Model.Name }}, int, error) {
	var objs []*model.{{ .Model.Name }}
	query := {{ $tick }}{{ .SqlGetAll}}{{ $tick }}

	// batch count query
	countClause, countArgs := p.Filter.ToSQL()
	countQuery := {{ $tick }}{{ .SqlGetCount}}{{ $tick }} + countClause
	count := 0 
	err := rps.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: can't get {{ .ModelNameLower }} record - %s", err)
	}

	// batch selection query
	selectionClause, args := p.ToSQL()
	selectionQuery := fmt.Sprintf("%s %s", query, selectionClause)
	rows, err := rps.pool.Query(ctx, selectionQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: can't get all {{ .ModelNameLower }} records - %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var obj *model.{{ .Model.Name }}
		err = rows.Scan({{range .Model.PrimaryFields}}
			&obj.{{ .Name }},{{ end }} {{ range .Model.Fields}}
			&obj.{{ .Name}}, {{ end }}
		)
		if err != nil {
			return nil, 0, fmt.Errorf("repository: can't parse {{ .ModelNameLower }} record - %s", err)
		}
		objs = append(objs, obj)
	}
	return objs, count, nil
}{{ end }}
`))

func SqlCreate(m *model.Model) string {
	sqlRequest := "INSERT INTO \"" + *m.TableName + "\" \n\t\t("
	for i := 0; i < len(m.PrimaryFields); i++ {
		sqlRequest += "\"" + *m.PrimaryFields[i].SqlName + "\", \n\t\t"
	}
	for i := 0; i < len(m.Fields)-1; i++ {
		sqlRequest += "\"" + *m.Fields[i].SqlName + "\", \n\t\t"
	}
	sqlRequest += "\"" + *m.Fields[len(m.Fields)-1].SqlName + "\") \n\t\t"
	sqlRequest += "VALUES ("
	for i := 1; i < len(m.Fields)+len(m.PrimaryFields); i++ {
		sqlRequest += "$" + strconv.Itoa(i) + ","
	}
	sqlRequest += "$" + strconv.Itoa(len(m.Fields)+len(m.PrimaryFields)) + ")"
	return sqlRequest
}

func CreateValues(m *model.Model) string {
	var values string
	for i := 0; i < len(m.PrimaryFields); i++ {
		values += "obj." + *m.PrimaryFields[i].Name + ", "
	}
	for i := 0; i < len(m.Fields)-1; i++ {
		values += "obj." + *m.Fields[i].Name + ", "
	}
	values += "obj." + *m.Fields[len(m.Fields)-1].Name
	return values
}

// SqlGetByID returns script for get request
func SqlGetByID(m *model.Model) string {
	sqlRequest := "SELECT \n\t\t"
	for i := 0; i < len(m.PrimaryFields); i++ {
		sqlRequest += "\"" + *m.PrimaryFields[i].SqlName + "\", \n\t\t"
	}
	for i := 0; i < len(m.Fields)-1; i++ {
		sqlRequest += "\"" + *m.Fields[i].SqlName + "\", \n\t\t"
	}
	sqlRequest += "\"" + *m.Fields[len(m.Fields)-1].SqlName + "\"\n\t\t"
	sqlRequest += "FROM \"" + *m.TableName + "\" WHERE "
	for i := 0; i < len(m.PrimaryFields)-1; i++ {
		sqlRequest += "\"" + *m.PrimaryFields[i].SqlName + "\" = $" + strconv.Itoa(i+1) + ", "
	}
	sqlRequest += "\"" + *m.PrimaryFields[len(m.PrimaryFields)-1].SqlName + "\" = $" + strconv.Itoa(len(m.PrimaryFields))
	return sqlRequest
}

// SqlGetAll returns sql script for get all request
func SqlGetAll(m *model.Model) string {
	sqlRequest := "SELECT \n\t\t"
	for i := 0; i < len(m.PrimaryFields); i++ {
		sqlRequest += "\"" + *m.PrimaryFields[i].SqlName + "\", \n\t\t"
	}
	for i := 0; i < len(m.Fields)-1; i++ {
		sqlRequest += "\"" + *m.Fields[i].SqlName + "\", \n\t\t"
	}
	sqlRequest += "\"" + *m.Fields[len(m.Fields)-1].SqlName + "\" \n\t\tFROM \"" + *m.TableName + "\""
	return sqlRequest
}

// SqlGetCount returns sql script for get count of element request
func SqlGetCount(m *model.Model) string {
	sqlRequest := "SELECT COUNT(1) FROM \"" + *m.TableName + "\" "
	return sqlRequest
}

// SqlUpdate returns sql script for update request
func SqlUpdate(m *model.Model) string {
	sqlRequest := "UPDATE \"" + *m.TableName + "\" SET \n\t\t"
	for i := 0; i < len(m.Fields)-1; i++ {
		sqlRequest += "\"" + *m.Fields[i].SqlName + "\" = $" + strconv.Itoa(i+1) + ", \n\t\t"
	}
	sqlRequest += "\"" + *m.Fields[len(m.Fields)-1].SqlName + "\" = $" + strconv.Itoa(len(m.Fields)) + " \n\t\t"
	sqlRequest += "WHERE "
	for i := len(m.Fields); i < len(m.Fields)+len(m.PrimaryFields)-1; i++ {
		sqlRequest += "\"" + *m.PrimaryFields[i].SqlName + "\" = $" + strconv.Itoa(i+1) + ", "
	}
	sqlRequest += "\"" + *m.PrimaryFields[len(m.PrimaryFields)-1].SqlName + "\" = $" + strconv.Itoa(len(m.Fields)+len(m.PrimaryFields))
	return sqlRequest
}

// UpdateValues returns sql
func UpdateValues(m *model.Model) string {
	var values string
	for i := 0; i < len(m.Fields); i++ {
		values += "obj." + *m.Fields[i].Name + ", "
	}
	for i := 0; i < len(m.PrimaryFields)-1; i++ {
		values += "obj." + *m.PrimaryFields[i].Name + ", "
	}
	values += "obj." + *m.PrimaryFields[len(m.PrimaryFields)-1].Name
	return values
}

// SqlDelete method returns sql script for delete request
func SqlDelete(m *model.Model) string {
	sqlRequest := "DELETE FROM \"" + *m.TableName + "\" WHERE \n\t\t"
	for i := 0; i < len(m.PrimaryFields)-1; i++ {
		sqlRequest += "\"" + *m.PrimaryFields[i].SqlName + "\" = $" + strconv.Itoa(i+1) + ", "
	}
	sqlRequest += "\"" + *m.PrimaryFields[len(m.PrimaryFields)-1].SqlName + "\" = $" + strconv.Itoa(len(m.PrimaryFields))
	return sqlRequest
}

func PrimaryKeysString(m *model.Model) string {
	var primaryKeys string
	for i := 0; i < len(m.PrimaryFields)-1; i++ {
		primaryKeys += *m.PrimaryFields[i].SqlName + " " + *m.PrimaryFields[i].Type + ", "
	}
	primaryKeys += *m.PrimaryFields[len(m.PrimaryFields)-1].SqlName + " " + *m.PrimaryFields[len(m.PrimaryFields)-1].Type

	return primaryKeys
}

func PrimaryKeysValues(m *model.Model) string {
	var values string
	for i := 0; i < len(m.PrimaryFields)-1; i++ {
		values += *m.PrimaryFields[i].SqlName + ", "
	}
	values += *m.PrimaryFields[len(m.PrimaryFields)-1].SqlName
	return values
}
