package templates

import "text/template"

var PaginationTemplate = template.Must(template.New("").Parse(
	`// Package pagination applies logic of pagination for database request
package pagination

import "fmt"

// Filter interface represents data which can be converted to WHERE sql clause
type Filter interface {
	ToSQL() (string, []any)
}

// Sort contains sorting data
type Sort struct {
	OrderBy string
	Asc     bool
}

// ToSQL returns ORDER BY SQL clause if sorting is present
func (s Sort) ToSQL() string {
	if s.OrderBy == "" {
		return ""
	}
	if s.Asc {
		return fmt.Sprintf("ORDER BY %s ASC", s.OrderBy)
	}
	return fmt.Sprintf("ORDER BY %s DESC", s.OrderBy)
}

// Pagination helps to build filter, pagination and sorting SQL clause
type Pagination[T Filter] struct {
	Sort   Sort
	Limit  int
	Offset int
	Filter T
}

// ToSQL parses and build filter, pagination and sorting SQL clause
func (p Pagination[T]) ToSQL() (string, []any) {
	q, args := p.Filter.ToSQL()

	q = fmt.Sprintf("%s %s", q, p.Sort.ToSQL())

	if p.Limit > 0 {
		q = fmt.Sprintf("%s LIMIT %d", q, p.Limit)
	}

	if p.Offset > 0 {
		q = fmt.Sprintf("%s OFFSET %d", q, p.Offset)
	}

	return q, args
}`))
