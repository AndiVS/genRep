// Package generators provide functions for generating using templates
package generators

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/AndiVS/genRep/internal/helper"
	"github.com/AndiVS/genRep/internal/model"
	"github.com/AndiVS/genRep/internal/templates"
	"github.com/AndiVS/genRep/internal/terminal/ubuntu"
)

type repositoryTemplateParams struct {
	Model                    *model.Model
	ModelNameLower           string
	CreateMethod             bool
	CreateMethodTransaction  bool
	GetByPrimaryFieldMethod  bool
	GetAllMethod             bool
	UpdateMethod             bool
	UpdateMethodTransaction  bool
	DeleteMethod             bool
	DeleteMethodTransaction  bool
	GetWithSortAndPagination bool
	PrimaryKeys              string
	PrimaryValues            string
	SQLCreate                string
	CreateValues             string
	SQLGetByID               string
	SQLGetAll                string
	SQLGetCount              string
	SQLUpdate                string
	UpdateValues             string
	SQLDelete                string
	UUIDFieldExists          bool
	TimeFieldExists          bool
}

// GenerateRepository used to generate repository package
func GenerateRepository(models []*model.Model, outDir string) error {
	workingDir, err := ubuntu.CreateDirectory("repository", outDir)
	if err != nil {
		return err
	}

	for _, m := range models {
		rf, err := generateRepository(m)
		if err != nil {
			return err
		}
		fullPath := fmt.Sprintf("%s/%s_generated.go", workingDir, helper.ToSnakeCase(*m.Name))
		err = os.WriteFile(fullPath, rf.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("repository generator: can't write template into the file - %s", err)
		}
	}
	return nil
}

func generateRepository(m *model.Model) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	params := preparePostgresParams(m)

	err := templates.PostgresRepositoryTemplate.Execute(&buf, params)
	if err != nil {
		return nil, fmt.Errorf("repository generator: can't execute postgressql repository template - %s", err)
	}
	return &buf, nil
}

func preparePostgresParams(m *model.Model) *repositoryTemplateParams {
	params := &repositoryTemplateParams{
		Model:           m,
		PrimaryValues:   templates.PrimaryKeysValues(m),
		PrimaryKeys:     templates.PrimaryKeysString(m),
		TimeFieldExists: model.CheckTimeFields(m),
		UUIDFieldExists: model.CheckUUIDField(m),
		ModelNameLower:  strings.ToLower(*m.Name),
	}

	params.CreateMethod = true
	params.SQLCreate = templates.SQLCreate(m)
	params.CreateValues = templates.CreateValues(m)

	params.GetByPrimaryFieldMethod = true
	params.SQLGetByID = templates.SQLGetByID(m)

	params.GetAllMethod = true
	params.SQLGetAll = templates.SQLGetAll(m)

	params.UpdateMethod = true
	params.SQLUpdate = templates.SQLUpdate(m)
	params.UpdateValues = templates.UpdateValues(m)

	params.DeleteMethod = true
	params.SQLDelete = templates.SQLDelete(m)

	params.DeleteMethodTransaction = true
	params.GetWithSortAndPagination = true
	params.SQLGetCount = templates.SQLGetCount(m)

	return params
}
