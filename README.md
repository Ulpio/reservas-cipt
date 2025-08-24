
# Reservas-CIPT

Backend do sistema de reservas do CIPT Jaraguá.

## Requisitos

- [Go](https://go.dev/) 1.20+
- [Docker](https://www.docker.com/) e Docker Compose (opcional, para o banco de dados)

## Configuração

1. Crie um arquivo `.env` na raiz do projeto com as variáveis abaixo:

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=reservas
   DB_SSLMODE=disable
   JWT_SECRET=sua_chave_segura
   ```

2. (Opcional) Inicie um banco PostgreSQL via Docker:

   ```bash
   docker-compose up -d db
   ```

3. Execute a API:

   ```bash
   go run main.go
   ```

   O serviço estará disponível em `http://localhost:8080`.


4. Documentação (Swagger):

   A documentação está disponível em `http://localhost:8080/swagger/index.html`.
   
## Testes

Para rodar os testes dos services e handlers:

```bash
go test ./tests/services ./tests/handlers
```

Os testes usam um banco SQLite em memória e não dependem do PostgreSQL.
