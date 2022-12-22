package generators

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/AndiVS/genRep/internal/model"
	"github.com/AndiVS/genRep/internal/templates"
	"github.com/AndiVS/genRep/internal/terminal/ubuntu"
)

func GeneratePagination(m *model.Model, outDir string) error {
	workingDir, err := ubuntu.CreateDirectory("pagination", outDir)
	if err != nil {
		return err
	}
	sf, err := generatepag()
	if err != nil {
		return err
	}
	fullPath := fmt.Sprintf("%s/%s.go", workingDir, strings.ToLower(*m.Name))
	err = os.WriteFile(fullPath, sf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("service generator: can't write template into the file - %s", err)
	}
	return nil
}

func generatepag() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := templates.PaginationTemplate.Execute(&buf, nil)
	if err != nil {
		return nil, fmt.Errorf("service generator: can't execute service template - %s", err)
	}
	return &buf, nil
}
