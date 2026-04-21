package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pultrini/scrapler/models"
)

var DB *sql.DB

func Connect(connStr string) {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("erro ao conectar:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Erro ao pingar banco:", err)
	}
	runMigrations()
	log.Println("Banco conectado e migrations aplicadas!")
}

func runMigrations() {
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		log.Fatal("Erro ao criar driver de migration:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("Erro ao criar migrate:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Erro ao rodar migrations:", err)
	}
}

func InsertConcurso(c models.Concurso) error {
	query := `
        INSERT INTO concursos (titulo, faixa_inicial, faixa_final, escolaridade, resumo_vaga, link, origem, edital_link)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (titulo) DO NOTHING
    `
	_, err := DB.Exec(query,
		c.Titulo,
		c.FaixaInicial,
		c.FaixaFinal,
		c.Escolaridade,
		c.ResumoVaga,
		c.Link,
		c.Origem,
		c.EditalLink,
	)
	return err
}
