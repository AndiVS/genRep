package generators

import (
	"bytes"
	"fmt"
	"os"

	"github.com/AndiVS/genRep/internal/model"
	"github.com/AndiVS/genRep/internal/templates"
	"github.com/AndiVS/genRep/internal/terminal/ubuntu"
)

type migrationsTemplateParams struct {
	Model           *model.Model
	MultiPrimaryKey string
}

// GenerateSQLMigration generates migrations for each entity from the description and writes it in file
func GenerateSQLMigration(models []*model.Model, outDir string) error {
	fullPath, err := ubuntu.GetFullPath("migrations", outDir)
	if err != nil {
		return fmt.Errorf("migrations generator: - %w", err)
	}

	exists, err := ubuntu.CheckDirectory(fullPath)
	if err != nil {
		return fmt.Errorf("migrations generator: - %w", err)
	}
	if !exists {
		err := ubuntu.CreateDirectory(fullPath)
		if err != nil {
			return fmt.Errorf("migrations generator: - %w", err)
		}
	}

	for i, m := range models {
		fileName := generateFileName(*m.Name, i+1)
		migration, err := generateSQLMigration(m)
		if err != nil {
			return err
		}
		fullPath := fmt.Sprintf("%s/%s", fullPath, fileName)
		err = os.WriteFile(fullPath, migration.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("migrations generator: can't write template into the file - %s", err)
		}
	}
	return nil
}

// generateFileName returns migration file version
func generateFileName(entityName string, migrationVersion int) string {
	fileName := "V1_" + fmt.Sprint(migrationVersion) + "__Init" + entityName + "Table.sql"
	return fileName
}

// generateSQLMigration returns generated migration template
func generateSQLMigration(m *model.Model) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	params := migrationsTemplateParams{
		Model:           m,
		MultiPrimaryKey: templates.PrepareMulticolumnPrimaryKey(m.PrimaryFields),
	}
	err := templates.PostgresCreateTableTemplate.Execute(&buf, params)
	if err != nil {
		return nil, fmt.Errorf("migrations generator: can't execute sql migration template - %s", err)
	}

	return &buf, nil
}
