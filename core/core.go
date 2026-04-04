package core

import (
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/export"
	"github.com/pultrini/scrapler/spiders"
)

func Run(spider spiders.Spiders) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: spider.StartsRequests,
		Exporters: []export.Exporter{
			&export.JSON{},
		},
	}).Start()
}
