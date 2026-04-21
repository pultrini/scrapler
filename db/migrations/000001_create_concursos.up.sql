CREATE TABLE IF NOT EXISTS concursos (
         id SERIAL PRIMARY KEY,
         titulo TEXT NOT NULL UNIQUE,
         faixa_inicial TEXT,
         faixa_final TEXT,
         escolaridade TEXT,
         resumo_vaga TEXT,
         link TEXT,
         origem TEXT,
         edital_link TEXT,
         criado_em TIMESTAMP DEFAULT NOW()
);