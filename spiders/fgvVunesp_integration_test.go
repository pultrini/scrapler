package spiders_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/geziyor/geziyor"
	"github.com/pultrini/scrapler/models"
	"github.com/pultrini/scrapler/spiders"
)

// mockStorage implementa db.Storage para os testes
type mockStorage struct {
	saved []models.Concurso
}

func (m *mockStorage) InsertConcurso(c models.Concurso) error {
	m.saved = append(m.saved, c)
	return nil
}

func TestVunespSpider_Integration(t *testing.T) {
	// Servidor fake para a página de editais
	editalServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
			<html><body>
				<a href="/documento/stream/123" title="Edital de Abertura">Edital de Abertura</a>
				<a href="/documento/stream/456" title="Gabarito">Gabarito</a>
			</body></html>
		`)
	}))
	defer editalServer.Close()

	// Servidor fake para a listagem de concursos
	// Aponta o link do concurso para o editalServer
	listServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<html><body>
				<article class="concurso">
					<div class="titulo">Câmara Municipal de Teste</div>
					<div class="escolaridade">Superior</div>
					<div class="course-description"><p>Contador</p></div>
					<div class="course-informations">
						<span class="negrito">R$ 4.000,00</span>
						<span class="negrito">R$ 8.000,00</span>
					</div>
					<div class="read-more-box">
						<a href="%s/concurso/CMT1">Ver mais</a>
					</div>
				</article>
			</body></html>
		`, editalServer.URL)
	}))
	defer listServer.Close()

	store := &mockStorage{}
	spider := &spiders.VunespSpider{
		BaseSpider:    spiders.BaseSpider{Storage: store},
		EditalBaseURL: editalServer.URL,
	}

	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.Get(listServer.URL, spider.Parse)
		},
		ConcurrentRequests: 1,
		RequestDelay:       0,
	}).Start()

	time.Sleep(500 * time.Millisecond)

	if len(store.saved) == 0 {
		t.Fatal("nenhum concurso foi salvo")
	}

	c := store.saved[0]
	if c.Titulo != "Câmara Municipal de Teste" {
		t.Errorf("título errado: %q", c.Titulo)
	}
	if !strings.Contains(c.EditalLink, "stream/123") {
		t.Errorf("edital link errado: %q", c.EditalLink)
	}
}
