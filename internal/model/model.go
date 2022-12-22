// Package model represents models which are needed during code generation
package model

// Model struct which contains all necessary information about model
type Model struct {
	Name          *string
	TableName     *string
	Schema        *string
	PrimaryFields []*Field
	Fields        []*Field
}

// Field struct which contains all necessary information about model field
type Field struct {
	Name    *string
	SqlName *string
	Type    *string
}

func CheckUUIDField(m *Model) bool {
	for _, field := range m.Fields {
		if *field.Type == "uuid.UUID" {
			return true
		}
	}
	for _, field := range m.PrimaryFields {
		if *field.Type == "uuid.UUID" {
			return true
		}
	}
	return false
}

func CheckTimeFields(m *Model) bool {
	for _, field := range m.Fields {
		if *field.Type == "time.Time" {
			return true
		}
	}
	for _, field := range m.PrimaryFields {
		if *field.Type == "time.Time" {
			return true
		}
	}
	return false
}
