package db

import "github.com/pultrini/scrapler/models"

type Storage interface {
	InsertConcursos(c models.Concurso) error
}
