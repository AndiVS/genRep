package templates

import "text/template"

var ServiceTemplate = template.Must(template.New("").Parse(
	`// Package service applies
package service

import(
	"context"{{ if .TimeFieldExists }}
	"time"{{ end }}

	"{{ .Model.ModelPath }}"

	"github.com/AndiVS/pagination"{{ if .UUIDFieldExists}}
	"github.com/google/uuid"{{ end}}
)

// {{ .Model.Name}}ServiceManager is an interface 
type {{ .Model.Name}}ServiceManager interface { {{ if .CreateMethod }}
	Create(ctx context.Context, obj *model.{{ .Model.Name }}) error {{ end }} {{ if .GetByPrimaryFieldMethod }}
	GetByPrimaryField(ctx context.Context, {{ .PrimaryKeys}}) (*model.{{ .Model.Name }}, error) {{ end }} {{ if .GetAllMethod }}
	GetAll(ctx context.Context) ([]model.{{ .Model.Name }}, error) {{ end }} {{ if .UpdateMethod}}
	Update(ctx context.Context, obj *model.{{ .Model.Name }}) error {{ end }} {{ if .DeleteMethod }}
	Delete(ctx context.Context, {{ .PrimaryKeys}}) error {{ end}} {{ if .GetWithSortAndPagination }}
	GetWithSortAndPagination(ctx context.Context, p *pagination.Pagination[*model.{{ .Model.Name }}]) ([]*model.{{ .Model.Name }}, int, error) {{ end }}
}

// New{{ .Model.Name }}ServiceManager returns a new service manager instance
func New{{ .Model.Name }}ServiceManager(rps repository.{{ .Model.Name}}RepositoryManager) {{ .Model.Name}}ServiceManager {
	return &{{ .ModelNameLower}}ServiceManager{rps: rps}
}

// {{ .ModelNameLower}}ServiceManager is {{ .Model.Name}}ServiceManager implementation
type {{ .ModelNameLower}}ServiceManager struct {
	rps repository.{{ .Model.Name}}RepositoryManager
}

{{ if .CreateMethod}}// Create method 
func (s *{{ .ModelNameLower }}ServiceManager) Create(ctx context.Context, obj *model.{{ .Model.Name}}) error {
	return s.rps.Create(ctx, obj)
}{{ end}}

{{ if .GetByPrimaryFieldMethod}}// GetByPrimaryField method
func (s *{{ .ModelNameLower }}ServiceManager) GetByPrimaryField(ctx context.Context, {{ .PrimaryKeys}}) (*model.{{ .Model.Name}}, error) {
	return s.rps.GetByPrimaryField(ctx, {{ .PrimaryValues}})
}{{ end}}

{{ if .GetAllMethod}}// GetAll method
func (s *{{ .ModelNameLower }}ServiceManager) GetAll(ctx context.Context) ([]model.{{ .Model.Name }}, error) {
	return s.rps.GetAll(ctx)
} {{ end}}

{{ if .UpdateMethod}}// Update method
func (s *{{ .ModelNameLower }}ServiceManager) Update(ctx context.Context, obj *model.{{ .Model.Name }}) error {
	return s.rps.Update(ctx, obj)
} {{ end}}

{{ if .DeleteMethod}}// Delete method
func (s *{{ .ModelNameLower }}ServiceManager) Delete(ctx context.Context, {{ .PrimaryKeys}}) error {
	return s.rps.Delete(ctx, {{ .PrimaryValues}})
} {{ end}}

{{ if .GetWithSortAndPagination}}// GetWithSortAndPagination method
func (s *{{ .ModelNameLower }}ServiceManager) GetWithSortAndPagination(ctx context.Context, 
	p *pagination.Pagination[*model.{{ .Model.Name }}]) ([]*model.{{ .Model.Name }}, int, error) {
	return s.rps.GetWithSortAndPagination(ctx,p)
} {{ end}}
`))
