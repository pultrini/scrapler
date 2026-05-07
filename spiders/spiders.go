package spiders

import (
	"github.com/geziyor/geziyor"
	"github.com/pultrini/scrapler/db"
)

type Spiders interface {
	Name() string
	StartsRequests(g *geziyor.Geziyor)
}

type BaseSpider struct {
	Storage db.Storage
}
