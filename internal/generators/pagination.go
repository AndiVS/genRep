package generators

import (
	"bytes"
	"fmt"
	"os"

	"github.com/AndiVS/genRep/internal/templates"
	"github.com/AndiVS/genRep/internal/terminal/ubuntu"
)

// GeneratePagination used to generate pagination package
func GeneratePagination(outDir string) error {
	workingDir, err := ubuntu.CreateDirectory("pagination", outDir)
	if err != nil {
		return err
	}
	sf, err := generatePagination()
	if err != nil {
		return err
	}
	fullPath := fmt.Sprintf("%s/%s.go", workingDir, "paginator")
	err = os.WriteFile(fullPath, sf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("pagination generator: can't write template into the file - %s", err)
	}
	return nil
}

func generatePagination() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := templates.PaginationTemplate.Execute(&buf, nil)
	if err != nil {
		return nil, fmt.Errorf("pagination generator: can't execute service template - %s", err)
	}
	return &buf, nil
}
