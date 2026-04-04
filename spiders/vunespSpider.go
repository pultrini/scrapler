package spiders

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/pultrini/scrapler/models"
)

type VunespSpider struct {
}

func (v *VunespSpider) Name() string {
	return "vunesp"
}
func (v *VunespSpider) StartsRequests(g *geziyor.Geziyor) {
	g.Get("https://www.vunesp.com.br/busca/concurso/inscricoes%20abertas", v.parse)
}

func (v *VunespSpider) parse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("article.concurso").Each(func(i int, s *goquery.Selection) {
		var concursos models.Concurso

		concursos.Titulo = strings.TrimSpace(s.Find(".titulo").Text())
		concursos.Escolaridade = strings.TrimSpace(s.Find(".escolaridade").Text())

		s.Find(".course-informations .negrito").Each(func(i int, n *goquery.Selection) {
			text := strings.TrimSpace(n.Text())

			if strings.Contains(text, ",") {
				if concursos.FaixaInicial == "" {
					concursos.FaixaInicial = text
				} else {
					concursos.FaixaFinal = text
				}
			}
		})

		concursos.ResumoVaga = strings.TrimSpace(
			s.Find(".course-description p").Text(),
		)

		if link, ok := s.Find(".read-more-box a").Attr("href"); ok {
			u, _ := r.Request.URL.Parse(link)
			concursos.Link = u.String()
		}
		concursos.Origem = v.Name()

		if concursos.Titulo != "" {
			g.Exports <- concursos
		}
	})
}
