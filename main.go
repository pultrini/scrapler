package main

import (
	"github.com/pultrini/scrapler/core"
	"github.com/pultrini/scrapler/spiders"
)

func main() {
	core.Run(&spiders.VunespSpider{})
	core.Run(&spiders.FGVSpider{})
}
