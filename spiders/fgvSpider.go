package spiders

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/pultrini/scrapler/models"
)

type FGVSpider struct {
	BaseSpider
}

func (f *FGVSpider) Name() string {
	return "fgv"
}

func (f *FGVSpider) StartsRequests(g *geziyor.Geziyor) {
	g.Get("https://conhecimento.fgv.br/concursos", f.parse)
}

func ParseConcursosFGV(doc *goquery.Document, baseURL string) []models.Concurso {
	var result []models.Concurso

	doc.Find(".views-row").Each(func(i int, s *goquery.Selection) {
		var c models.Concurso

		c.Titulo = strings.TrimSpace(s.Find(".views-field-title").Text())
		c.Origem = "fgv"

		if link, ok := s.Find("a").First().Attr("href"); ok {
			base := mustParseURL(baseURL)
			u, err := base.Parse(link)
			if err == nil {
				c.Link = u.String()
			}
		}

		result = append(result, c)
	})

	return result
}

func ParseEditalFGV(doc *goquery.Document, baseURL string) string {
	var editalLink string

	doc.Find("a").Each(func(j int, sel *goquery.Selection) {
		texto := strings.ToLower(sel.Text())
		href, _ := sel.Attr("href")

		if strings.Contains(texto, "edital") && strings.HasSuffix(strings.ToLower(href), ".pdf") {
			base := mustParseURL(baseURL)
			u, err := base.Parse(href)
			if err == nil {
				editalLink = u.String()
			}
		}
	})

	return editalLink
}

func (f *FGVSpider) parse(g *geziyor.Geziyor, r *client.Response) {
	concursos := ParseConcursosFGV(r.HTMLDoc, r.Request.URL.String())

	for _, c := range concursos {
		c := c
		g.Get(c.Link, func(g *geziyor.Geziyor, rInternal *client.Response) {
			c.EditalLink = ParseEditalFGV(rInternal.HTMLDoc, rInternal.Request.URL.String())

			if err := f.Storage.InsertConcurso(c); err != nil {
				log.Println("Erro ao salvar fgv:", err)
			}
		})
	}
}
