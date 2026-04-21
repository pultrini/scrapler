package spiders

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/pultrini/scrapler/models"
)

type FGVSpider struct {
	existentes map[string]bool
}

func (f *FGVSpider) loadExistentes() {
	f.existentes = make(map[string]bool)

	data, err := os.ReadFile("out.json")
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "[" || line == "]" {
			continue
		}
		line = strings.TrimSuffix(line, ",")

		var c models.Concurso
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			continue
		}
		if c.Origem == f.Name() {
			f.existentes[c.Titulo] = true
		}
	}
}

func (f *FGVSpider) Name() string {
	return "fgv"
}

func (f *FGVSpider) StartsRequests(g *geziyor.Geziyor) {
	f.loadExistentes()
	g.Get("https://conhecimento.fgv.br/concursos", f.parse)
}

func (f *FGVSpider) parse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find(".views-row").Each(func(i int, s *goquery.Selection) {
		var concurso models.Concurso

		concurso.Titulo = strings.TrimSpace(s.Find(".views-field-title").Text())
		concurso.Origem = f.Name()

		if link, ok := s.Find("a").First().Attr("href"); ok {
			u, _ := r.Request.URL.Parse(link)
			concurso.Link = u.String()

			if f.existentes[concurso.Titulo] {
				return
			}

			g.Get(concurso.Link, func(g *geziyor.Geziyor, rInternal *client.Response) {
				rInternal.HTMLDoc.Find("a").Each(func(j int, sel *goquery.Selection) {
					texto := strings.ToLower(sel.Text())
					href, _ := sel.Attr("href")

					if strings.Contains(texto, "edital") && strings.HasSuffix(strings.ToLower(href), ".pdf") {
						editalUrl, _ := rInternal.Request.URL.Parse(href)
						concurso.EditalLink = editalUrl.String()
					}
				})

				if concurso.Titulo != "" {
					g.Exports <- concurso
				}
			})
		}
	})
}
