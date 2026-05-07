package spiders

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/pultrini/scrapler/models"
)

type VunespSpider struct {
	BaseSpider
	EditalBaseURL string
}

func (v *VunespSpider) Name() string {
	return "vunesp"
}

func (v *VunespSpider) editalURL(sigla string) string {
	base := "https://documento.vunesp.com.br"
	if v.EditalBaseURL != "" {
		base = v.EditalBaseURL
	}
	return fmt.Sprintf("%s/projeto/%s/documento/", base, sigla)
}

func (v *VunespSpider) StartsRequests(g *geziyor.Geziyor) {
	g.Get("https://www.vunesp.com.br/busca/concurso/inscricoes%20abertas", v.Parse)
}

func ParseConcursosVunesp(doc *goquery.Document, baseURL string) []models.Concurso {
	var result []models.Concurso

	doc.Find("article.concurso").Each(func(i int, s *goquery.Selection) {
		var c models.Concurso

		c.Titulo = strings.TrimSpace(s.Find(".titulo").Text())
		c.Escolaridade = strings.TrimSpace(s.Find(".escolaridade").Text())
		c.Origem = "vunesp"

		s.Find(".course-informations .negrito").Each(func(i int, n *goquery.Selection) {
			text := strings.TrimSpace(n.Text())
			if strings.Contains(text, ",") {
				if c.FaixaInicial == "" {
					c.FaixaInicial = text
				} else {
					c.FaixaFinal = text
				}
			}
		})
		c.ResumoVaga = strings.TrimSpace(s.Find(".course-description p").Text())
		if link, ok := s.Find(".read-more-box a").Attr("href"); ok {
			u := mustParseURL(baseURL)
			resolved, _ := u.Parse(link)
			c.Link = resolved.String()
		}
		result = append(result, c)
	})
	return result
}

func (v *VunespSpider) Parse(g *geziyor.Geziyor, r *client.Response) {
	concursos := ParseConcursosVunesp(r.HTMLDoc, r.Request.URL.String())

	for _, c := range concursos {
		c := c
		v.fetchEdital(g, c)
	}
}

func (v *VunespSpider) fetchEdital(g *geziyor.Geziyor, c models.Concurso) {
	parts := strings.Split(strings.TrimRight(c.Link, "/"), "/")
	sigla := parts[len(parts)-1]
	apiURL := v.editalURL(sigla)

	g.Get(apiURL, func(g *geziyor.Geziyor, r *client.Response) {
		r.HTMLDoc.Find("a[href*='documento/stream']").Each(func(i int, s *goquery.Selection) {
			title, _ := s.Attr("title")
			if strings.Contains(strings.ToLower(title), "edital de abertura") {
				if href, ok := s.Attr("href"); ok {
					c.EditalLink = href
				}
			}
		})

		if err := v.Storage.InsertConcurso(c); err != nil {
			log.Println("Erro ao salvar vunesp:", err)
		}
	})
}
