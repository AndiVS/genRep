package validator

import (
	"fmt"

	"gitlab.effective-soft.com/gogenerator/repositoriBuilder/internal/model"
)

// Validate method
func Validate(m *model.Model) error {
	if err := modelValidation(m); err != nil {
		return fmt.Errorf("validator: %s", err)
	}
	return nil
}

func modelValidation(m *model.Model) error {
	if m.Name == nil {
		return fmt.Errorf("model error - name is missing")
	}
	if m.TableName == nil {
		return fmt.Errorf("%s model error - table name is missing", *m.Name)
	}
	if m.Schema == nil {
		return fmt.Errorf("%s model error - table schema is missing", *m.Name)
	}
	for _, f := range m.Fields {
		if err := fieldValidation(f); err != nil {
			return fmt.Errorf("%s model error - %s", *m.Name, err)
		}
	}
	for _, pf := range m.PrimaryFields {
		if err := fieldValidation(pf); err != nil {
			return fmt.Errorf("%s model error - %s", *m.Name, err)
		}
	}

	return nil
}

func fieldValidation(f *model.Field) error {
	if f.Name == nil {
		return fmt.Errorf("field name is missing")
	}
	if f.SqlName == nil {
		return fmt.Errorf("field sql name is missing")
	}
	if f.Type == nil {
		return fmt.Errorf("field type is missing")
	}
	return nil
}
