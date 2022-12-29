// Package model represents models which are needed during code generation
package model

// Model struct which contains all necessary information about model
type Model struct {
	Name          *string
	TableName     *string
	Schema        *string
	PrimaryFields []*Field
	Fields        []*Field
	ModelPath     *string
}

// Field struct which contains all necessary information about model field
type Field struct {
	Name    *string
	SQLName *string
	Type    *string
	SQLType *string
}

// CheckUUIDField function returns true if model contain field of type uuid
func CheckUUIDField(m *Model) bool {
	const uuidType = "uuid.UUID"
	for _, field := range m.Fields {
		if *field.Type == uuidType {
			return true
		}
	}
	for _, field := range m.PrimaryFields {
		if *field.Type == uuidType {
			return true
		}
	}
	return false
}

// CheckTimeFields function returns true if model contain field of type time
func CheckTimeFields(m *Model) bool {
	const timeType = "time.Time"
	for _, field := range m.Fields {
		if *field.Type == timeType {
			return true
		}
	}
	for _, field := range m.PrimaryFields {
		if *field.Type == timeType {
			return true
		}
	}
	return false
}
