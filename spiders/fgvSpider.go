package spiders

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/pultrini/scrapler/db"
	"github.com/pultrini/scrapler/models"
)

type FGVSpider struct{}

func (f *FGVSpider) Name() string {
	return "fgv"
}

func (f *FGVSpider) StartsRequests(g *geziyor.Geziyor) {
	g.Get("https://conhecimento.fgv.br/concursos", f.parse)
}

func (f *FGVSpider) parse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find(".views-row").Each(func(i int, s *goquery.Selection) {
		var c models.Concurso

		c.Titulo = strings.TrimSpace(s.Find(".views-field-title").Text())
		c.Origem = f.Name()

		if link, ok := s.Find("a").First().Attr("href"); ok {
			u, _ := r.Request.URL.Parse(link)
			c.Link = u.String()

			g.Get(c.Link, func(g *geziyor.Geziyor, rInternal *client.Response) {
				rInternal.HTMLDoc.Find("a").Each(func(j int, sel *goquery.Selection) {
					texto := strings.ToLower(sel.Text())
					href, _ := sel.Attr("href")

					if strings.Contains(texto, "edital") && strings.HasSuffix(strings.ToLower(href), ".pdf") {
						editalUrl, _ := rInternal.Request.URL.Parse(href)
						c.EditalLink = editalUrl.String()
					}
				})

				if err := db.InsertConcurso(c); err != nil {
					log.Println("Erro ao salvar fgv:", err)
				}
			})
		}
	})
}
