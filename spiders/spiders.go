package spiders

import "github.com/geziyor/geziyor"

type Spiders interface {
	Name() string
	StartsRequests(g *geziyor.Geziyor)
}
