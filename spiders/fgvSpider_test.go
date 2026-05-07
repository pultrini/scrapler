package spiders_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/pultrini/scrapler/spiders"
)

func TestParseConcursosFGV(t *testing.T) {
	html := `
		<div class="views-row">
			<div class="views-field-title"> Concurso Teste Alpha </div>
			<a href="/concursos/alpha">Ver edital</a>
		</div>
		<div class="views-row">
			<div class="views-field-title"> Concurso Beta </div>
			<a href="/concursos/beta">Ver edital</a>
		</div>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}

	result := spiders.ParseConcursosFGV(doc, "https://conhecimento.fgv.br")

	if len(result) != 2 {
		t.Fatalf("esperava 2 concursos, got %d", len(result))
	}
	if result[0].Titulo != "Concurso Teste Alpha" {
		t.Errorf("título errado: %q", result[0].Titulo)
	}
	if result[0].Link != "https://conhecimento.fgv.br/concursos/alpha" {
		t.Errorf("link errado: %q", result[0].Link)
	}
	if result[0].Origem != "fgv" {
		t.Errorf("origem errada: %q", result[0].Origem)
	}
}

func TestParseEditalFGV(t *testing.T) {
	html := `
		<a href="/docs/edital-abertura.pdf">Edital de Abertura</a>
		<a href="/docs/gabarito.pdf">Gabarito</a>`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	link := spiders.ParseEditalFGV(doc, "https://conhecimento.fgv.br")

	if link != "https://conhecimento.fgv.br/docs/edital-abertura.pdf" {
		t.Errorf("edital link errado: %q", link)
	}
}

func TestParseEditalFGV_SemEdital(t *testing.T) {
	html := `<a href="/docs/gabarito.pdf">Gabarito</a>`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	link := spiders.ParseEditalFGV(doc, "https://conhecimento.fgv.br")

	if link != "" {
		t.Errorf("esperava link vazio, got %q", link)
	}
}
