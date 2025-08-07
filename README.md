# 🏢 Sistema de Reservas - CIPT Jaraguá

Backend desenvolvido em Go para gerenciamento interno de reservas de espaços físicos no Centro de Inovação do Polo Tecnológico (CIPT - Jaraguá, Maceió). O sistema será utilizado exclusivamente pela equipe administrativa (recepcionistas e administradores).

---

## ⚙️ Tecnologias Utilizadas

- **Golang** (com arquitetura Clean)
- **Gin** (framework HTTP)
- **GORM** (ORM para PostgreSQL)
- **PostgreSQL + pgAdmin** via Docker
- **JWT** (futuro)
- **godotenv** para variáveis de ambiente
- **Makefile** (futuro)

---

## 📁 Estrutura do Projeto

```
cmd/                # Entry point da aplicação
internal/
  config/           # Leitura do .env e conexão com banco
  domain/           # Entidades e interfaces do domínio
  dto/              # Structs de entrada e saída da API
  handler/          # Controllers HTTP com Gin
  middleware/       # Autenticação, logs (em breve)
  repository/       # Implementações de persistência (GORM)
  server/           # Setup de rotas e servidor
  service/          # Serviços auxiliares (JWT, hash)
  usecase/          # Regras de negócio (application layer)
pkg/                # Pacotes reutilizáveis
test/               # Testes
.env                # Variáveis de ambiente
docker-compose.yml  # PostgreSQL + pgAdmin
```

---

## 🧪 Funcionalidades atuais

- ✅ Cadastro de usuários (`POST /users`)
- ✅ Listagem de usuários (`GET /users`)
- ✅ Buscar usuário por ID (`GET /users/:id`)
- ✅ Remoção de usuário (`DELETE /users/:id`)
- ⚠️ Sem autenticação por enquanto
- 🚧 Estrutura pronta para expansão (espaços, reservas, etc)

---

## 🐳 Como rodar o projeto

### Pré-requisitos

- [Go 1.22+](https://go.dev/)
- [Docker](https://www.docker.com/)

### 1. Clone o projeto

```bash
git clone https://github.com/seuprojeto/reservas-cipt.git
cd reservas-cipt
```

### 2. Configure o `.env`

Crie um `.env` na raiz com:

```env
POSTGRES_DB=cipt
POSTGRES_USER=ciptuser
POSTGRES_PASSWORD=ciptpass
POSTGRES_PORT=5432
```

### 3. Suba os containers

```bash
docker-compose up -d
```

### 4. Rode a aplicação

```bash
go run cmd/app/main.go
```

---

## 📬 Endpoints disponíveis

| Método | Rota           | Descrição                     |
|--------|----------------|-------------------------------|
| POST   | `/users`       | Cadastrar usuário             |
| GET    | `/users`       | Listar todos os usuários      |
| GET    | `/users/:id`   | Buscar usuário por ID         |
| DELETE | `/users/:id`   | Remover usuário               |

---

## 🛠️ Em desenvolvimento

- [ ] Middleware de autenticação (JWT)
- [ ] Recursos de reserva de espaços
- [ ] Controle de strikes e permissões
- [ ] Exportação de relatórios (PDF/CSV)
- [ ] Testes automatizados

---

## 📌 Licença

Este projeto é parte de um sistema interno e não possui licença pública ainda.

---

## ✍️ Autor

Desenvolvido por [Ulpio](https://github.com/Ulpio)
