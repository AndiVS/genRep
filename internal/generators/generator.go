package generators

import "gitlab.effective-soft.com/gogenerator/repositoriBuilder/internal/model"

type Generator interface {
	Generate(microservice *model.Model, outDir string) error
}
