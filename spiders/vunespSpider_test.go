package spiders_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/pultrini/scrapler/spiders"
)

func newDoc(t *testing.T, html string) *goquery.Document {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestParseConcursosVunesp_Titulo(t *testing.T) {
	doc := newDoc(t, `
		<article class="concurso">
			<div class="titulo">Prefeitura de São Paulo</div>
			<div class="escolaridade">Superior</div>
			<div class="course-description"><p>Analista de TI</p></div>
			<div class="read-more-box"><a href="/concurso/PREF-SP">Ver mais</a></div>
		</article>
	`)

	result := spiders.ParseConcursosVunesp(doc, "https://www.vunesp.com.br")

	if len(result) != 1 {
		t.Fatalf("esperava 1 concurso, got %d", len(result))
	}
	if result[0].Titulo != "Prefeitura de São Paulo" {
		t.Errorf("título errado: %q", result[0].Titulo)
	}
	if result[0].Escolaridade != "Superior" {
		t.Errorf("escolaridade errada: %q", result[0].Escolaridade)
	}
	if result[0].ResumoVaga != "Analista de TI" {
		t.Errorf("resumo errado: %q", result[0].ResumoVaga)
	}
	if result[0].Origem != "vunesp" {
		t.Errorf("origem errada: %q", result[0].Origem)
	}
}

func TestParseConcursosVunesp_FaixaSalarial(t *testing.T) {
	doc := newDoc(t, `
		<article class="concurso">
			<div class="titulo">TJ-SP</div>
			<div class="course-informations">
				<span class="negrito">R$ 3.500,00</span>
				<span class="negrito">R$ 7.000,00</span>
			</div>
		</article>
	`)

	result := spiders.ParseConcursosVunesp(doc, "https://www.vunesp.com.br")

	if result[0].FaixaInicial != "R$ 3.500,00" {
		t.Errorf("faixa inicial errada: %q", result[0].FaixaInicial)
	}
	if result[0].FaixaFinal != "R$ 7.000,00" {
		t.Errorf("faixa final errada: %q", result[0].FaixaFinal)
	}
}

func TestParseConcursosVunesp_Link(t *testing.T) {
	doc := newDoc(t, `
		<article class="concurso">
			<div class="titulo">UFABC</div>
			<div class="read-more-box"><a href="/concurso/UFABC1">Ver mais</a></div>
		</article>
	`)

	result := spiders.ParseConcursosVunesp(doc, "https://www.vunesp.com.br")

	expected := "https://www.vunesp.com.br/concurso/UFABC1"
	if result[0].Link != expected {
		t.Errorf("link errado: got %q, want %q", result[0].Link, expected)
	}
}

func TestParseConcursosVunesp_Vazio(t *testing.T) {
	doc := newDoc(t, `<html><body><p>Nenhum concurso</p></body></html>`)

	result := spiders.ParseConcursosVunesp(doc, "https://www.vunesp.com.br")

	if len(result) != 0 {
		t.Errorf("esperava slice vazia, got %d itens", len(result))
	}
}

func TestParseConcursosVunesp_Multiplos(t *testing.T) {
	doc := newDoc(t, `
		<article class="concurso"><div class="titulo">Concurso A</div></article>
		<article class="concurso"><div class="titulo">Concurso B</div></article>
		<article class="concurso"><div class="titulo">Concurso C</div></article>
	`)

	result := spiders.ParseConcursosVunesp(doc, "https://www.vunesp.com.br")

	if len(result) != 3 {
		t.Fatalf("esperava 3 concursos, got %d", len(result))
	}
}
