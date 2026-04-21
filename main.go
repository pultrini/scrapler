package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pultrini/scrapler/core"
	"github.com/pultrini/scrapler/db"
	"github.com/pultrini/scrapler/spiders"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	connStr := os.Getenv("SUPABASE_URL")
	db.Connect(connStr)

	core.Run(&spiders.VunespSpider{})
	core.Run(&spiders.FGVSpider{})
}
