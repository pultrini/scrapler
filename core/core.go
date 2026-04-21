package core

import (
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/export"
	"github.com/pultrini/scrapler/spiders"
)

func Run(spider spiders.Spiders) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc:  spider.StartsRequests,
		ConcurrentRequests: 10,
		RequestDelay:       1,
		Exporters: []export.Exporter{
			&export.JSON{},
		},
	}).Start()
}
