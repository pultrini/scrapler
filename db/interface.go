package db

import "github.com/pultrini/scrapler/models"

type Storage interface {
	InsertConcurso(c models.Concurso) error
}
