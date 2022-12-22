package generators

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gitlab.effective-soft.com/gogenerator/repositoriBuilder/internal/helper"
	"gitlab.effective-soft.com/gogenerator/repositoriBuilder/internal/model"
	"gitlab.effective-soft.com/gogenerator/repositoriBuilder/internal/templates"
	"gitlab.effective-soft.com/gogenerator/repositoriBuilder/internal/terminal/ubuntu"
)

type repositoryTemplateParams struct {
	Model                    *model.Model
	ModelNameLower           string
	CreateMethod             bool
	CreateMethodTransaction  bool
	GetByIDMethod            bool
	GetAllMethod             bool
	UpdateMethod             bool
	UpdateMethodTransaction  bool
	DeleteMethod             bool
	DeleteMethodTransaction  bool
	GetWithSortAndPagination bool
	PrimaryKeys              string
	PrimaryValues            string
	SqlCreate                string
	CreateValues             string
	SqlGetByID               string
	SqlGetAll                string
	SqlGetCount              string
	SqlUpdate                string
	UpdateValues             string
	SqlDelete                string
	UUIDFieldExists          bool
	TimeFieldExists          bool
}

func Generate(m *model.Model, outDir string) error {
	workingDir, err := ubuntu.CreateDirectory("repository", outDir)
	if err != nil {
		return err
	}
	rf, err := generateRepository(m)
	if err != nil {
		return err
	}
	fullPath := fmt.Sprintf("%s/%s_generated.go", workingDir, helper.ToSnakeCase(*m.Name))
	err = os.WriteFile(fullPath, rf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("repository generator: can't write template into the file - %s", err)
	}
	return nil
}

func generateRepository(m *model.Model) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	params, err := preparePostgresParams(m)
	if err != nil {
		return nil, err
	}
	err = templates.PostgresRepositoryTemplate.Execute(&buf, params)
	if err != nil {
		return nil, fmt.Errorf("repository generator: can't execute postgressql repository template - %s", err)
	}
	return &buf, nil
}

func preparePostgresParams(m *model.Model) (*repositoryTemplateParams, error) {
	params := &repositoryTemplateParams{
		Model:           m,
		PrimaryValues:   templates.PrimaryKeysValues(m),
		PrimaryKeys:     templates.PrimaryKeysString(m),
		TimeFieldExists: model.CheckTimeFields(m),
		UUIDFieldExists: model.CheckUUIDField(m),
		ModelNameLower:  strings.ToLower(*m.Name),
	}

	params.CreateMethod = true
	params.SqlCreate = templates.SqlCreate(m)
	params.CreateValues = templates.CreateValues(m)

	params.GetByIDMethod = true
	params.SqlGetByID = templates.SqlGetByID(m)

	params.GetAllMethod = true
	params.SqlGetAll = templates.SqlGetAll(m)

	params.UpdateMethod = true
	params.SqlUpdate = templates.SqlUpdate(m)
	params.UpdateValues = templates.UpdateValues(m)

	params.DeleteMethod = true
	params.SqlDelete = templates.SqlDelete(m)

	params.DeleteMethodTransaction = true
	params.GetWithSortAndPagination = true
	params.SqlGetCount = templates.SqlGetCount(m)

	return params, nil
}
