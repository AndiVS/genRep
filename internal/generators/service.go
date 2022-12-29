package generators

import (
	"bytes"
	"fmt"
	"os"

	"github.com/AndiVS/genRep/internal/helper"
	"github.com/AndiVS/genRep/internal/model"
	"github.com/AndiVS/genRep/internal/templates"
	"github.com/AndiVS/genRep/internal/terminal/ubuntu"
)

type serviceTemplateParams struct {
	Model                    *model.Model
	ModelNameLower           string
	PrimaryKeys              string
	PrimaryValues            string
	CreateMethod             bool
	GetByPrimaryFieldMethod  bool
	GetAllMethod             bool
	UpdateMethod             bool
	DeleteMethod             bool
	GetWithSortAndPagination bool
	UUIDFieldExists          bool
	TimeFieldExists          bool
}

func GenerateService(models []*model.Model, outDir string) error {
	workingDir, err := ubuntu.CreateDirectory("service", outDir)
	if err != nil {
		return err
	}
	for _, m := range models {
		sf, err := generateService(m)
		if err != nil {
			return err
		}
		fullPath := fmt.Sprintf("%s/%s_generated.go", workingDir, helper.ToSnakeCase(*m.Name))
		err = os.WriteFile(fullPath, sf.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("service generator: can't write template into the file - %s", err)
		}
	}
	return nil
}

func generateService(m *model.Model) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	params := serviceTemplateParams{
		Model:           m,
		PrimaryKeys:     templates.PrimaryKeysString(m),
		PrimaryValues:   templates.PrimaryKeysValues(m),
		TimeFieldExists: model.CheckTimeFields(m),
		UUIDFieldExists: model.CheckUUIDField(m),
	}

	modelNameLower := helper.LcFirst(*m.Name)
	params.ModelNameLower = modelNameLower

	params.CreateMethod = true
	params.GetAllMethod = true
	params.GetByPrimaryFieldMethod = true
	params.UpdateMethod = true
	params.DeleteMethod = true
	params.GetWithSortAndPagination = true

	err := templates.ServiceTemplate.Execute(&buf, params)
	if err != nil {
		return nil, fmt.Errorf("service generator: can't execute service template - %s", err)
	}
	return &buf, nil
}
