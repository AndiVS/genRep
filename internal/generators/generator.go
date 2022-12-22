package generators

import "github.com/AndiVS/genRep/internal/model"

type Generator interface {
	Generate(microservice *model.Model, outDir string) error
}
