// Package validator is used to validate model for generation
package validator

import (
	"fmt"

	"github.com/AndiVS/genRep/internal/model"
)

// Validate method
func Validate(m []*model.Model) error {
	for _, v := range m {
		if err := modelValidation(v); err != nil {
			return fmt.Errorf("validator: %s", err)
		}
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
	if len(m.PrimaryFields) == 0 {
		return fmt.Errorf("%s model error - at least one field must be primary;"+
			" set it by adding the `primary` tag to the structure field ", *m.Name)
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
	if f.SQLName == nil {
		return fmt.Errorf("field sql name is missing")
	}
	if f.Type == nil {
		return fmt.Errorf("field type is missing")
	}
	return nil
}
