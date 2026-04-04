package models

type Concurso struct {
	Titulo       string `json:"titulo"`
	FaixaInicial string `json:"faixa_inicial"`
	FaixaFinal   string `json:"faixa_final"`
	Escolaridade string `json:"escolaridade"`
	ResumoVaga   string `json:"resumo_vaga"`
	Link         string `json:"link"`
	Origem       string `json:"origem"`
}
