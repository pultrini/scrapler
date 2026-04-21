package spiders

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/pultrini/scrapler/models"
)

type VunespSpider struct {
	existentes map[string]bool
}

func (v *VunespSpider) loadExistentes() {
	v.existentes = make(map[string]bool)

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
		if c.Origem == v.Name() {
			v.existentes[c.Titulo] = true
		}
	}
}

func (v *VunespSpider) Name() string {
	return "vunesp"
}
func (v *VunespSpider) StartsRequests(g *geziyor.Geziyor) {
	v.loadExistentes()
	g.Get("https://www.vunesp.com.br/busca/concurso/inscricoes%20abertas", v.parse)
}

func (v *VunespSpider) parse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("article.concurso").Each(func(i int, s *goquery.Selection) {
		var concursos models.Concurso

		concursos.Titulo = strings.TrimSpace(s.Find(".titulo").Text())
		concursos.Escolaridade = strings.TrimSpace(s.Find(".escolaridade").Text())
		concursos.Origem = v.Name()

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
			c := concursos
			v.fetchEdital(g, c)
		}
	})
}

type VunespDocumento struct {
	Titulo string `json:"titulo"`
	URL    string `json:"url"`
}

func (v *VunespSpider) fetchEdital(g *geziyor.Geziyor, c models.Concurso) {

	if v.existentes[c.Titulo] {
		return
	}

	parts := strings.Split(strings.TrimRight(c.Link, "/"), "/")
	sigla := parts[len(parts)-1]

	apiURL := fmt.Sprintf("https://documento.vunesp.com.br/projeto/%s/documento/", sigla)

	g.Get(apiURL, func(g *geziyor.Geziyor, r *client.Response) {
		r.HTMLDoc.Find("a[href*='documento/stream']").Each(func(i int, s *goquery.Selection) {
			title, _ := s.Attr("title")
			if strings.Contains(strings.ToLower(title), "edital de abertura") {
				if href, ok := s.Attr("href"); ok {
					c.EditalLink = href
				}
			}
		})

		if c.Titulo != "" {
			g.Exports <- c
		}
	})
}
