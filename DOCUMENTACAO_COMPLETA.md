## 📚 Guia Completo da Plataforma de Reservas CIPT

### 1. Visão Geral
- **Projeto:** Backend do sistema de reservas do CIPT Jaraguá.
- **Stack:** Go 1.20+, Gin, GORM (PostgreSQL), autenticação JWT.
- **Swagger:** `http://{host}:8080/swagger/index.html`
- **Repositório:** `https://github.com/Ulpio/reservas-cipt`

---

### 2. Configuração Inicial

#### 2.1 Variáveis de Ambiente (.env)
```env
# Banco de dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=reservas
DB_SSLMODE=disable

# JWT
JWT_SECRET=sua_chave_segura

# Email (ver detalhes na seção 6)
EMAIL_ENABLED=false
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=seu-email@gmail.com
SMTP_PASSWORD=sua-senha-app
EMAIL_FROM=noreply@cipt.com
```

#### 2.2 Execução
```bash
go run main.go          # inicia API em http://localhost:8080
go build ./...          # compila
go test ./...           # roda testes automatizados
```

#### 2.3 Banco via Docker (opcional)
```bash
docker-compose up -d db
```

---

### 3. Autenticação
- **Endpoint:** `POST /api/v1/auth/login`
- **Headers:** `Content-Type: application/json`
- **Body:**
```json
{
  "cpf": "00000000000",
  "password": "senha"
}
```
- **Resposta 200:**
```json
{ "token": "jwt.token.aqui" }
```
- Utilize o token nas requisições seguintes:
```
Authorization: Bearer {token}
```

---

### 4. Recursos da API

#### 4.1 Clientes (`/clientes`)
| Método | Endpoint | Permissão | Descrição |
|--------|----------|-----------|-----------|
| `POST` | `/clientes/buscar-criar` | Qualquer token | Busca por CPF e cria se não existir (legado) |
| `POST` | `/clientes` | Qualquer token | Cria novo cliente |
| `GET`  | `/clientes/{cpf}` | Qualquer token | Busca por CPF (aceita com/sem máscara) |
| `GET`  | `/clientes/phone/{phone}` | Qualquer token | Busca por telefone |
| `GET`  | `/clientes/{id}` | Qualquer token | Busca por ID |
| `GET`  | `/clientes` | Qualquer token | Lista todos os clientes |
| `PATCH`| `/clientes/{id}` | Qualquer token | Atualiza dados |

**DTOs**
- Entrada: `ClienteInputDTO` → `name`, `cpf`, `birth_date`, `email`, `phone`
- Saída: `ClienteOutputDTO` → `client_id`, `name`, `cpf`, `birth_date`, `email`, `phone`, `strikes_count`, `total_reservations`, `created_at`

**Validações chave**
- Nome: 3-255 caracteres.
- CPF: 11 dígitos, dígitos verificadores válidos, único.
- Data de nascimento: formato `YYYY-MM-DD`, idade mínima de 18 anos.
- Email: formato válido (opcional).
- Telefone: 10-11 dígitos numéricos, único.

---

#### 4.2 Reservas (`/reservas`)
| Método | Endpoint | Permissão | Descrição |
|--------|----------|-----------|-----------|
| `POST` | `/reservas` | Recepcionista/Admin | Cria reserva |
| `GET`  | `/reservas` | Recepcionista/Admin | Lista reservas |
| `GET`  | `/reservas/{id}` | Recepcionista/Admin | Detalha reserva |
| `DELETE` | `/reservas/{id}` | Recepcionista/Admin | Cancela reserva (envia email) |

**DTO de criação**
```json
{
  "client_id": 1,
  "receptionist_id": 10,
  "space_id": 2,
  "date": "2025-11-13",
  "start_time": "16:00",
  "duration_hours": 2
}
```

**Regras importantes**
- Datas no formato `YYYY-MM-DD`, horário `HH:MM` ou `HH:MM:SS`.
- Não permitir datas passadas (considera timezone local).
- `duration_hours`: 1-24.
- Verifica conflitos de horário no mesmo espaço.
- Espaço deve estar com status `ativo`.
- Usuário com `receptionist_id` precisa ser recepcionista ou admin.

**Respostas de erro comuns**
- `400` datas ou horários inválidos.
- `404` cliente, espaço ou recepcionista inexistente.
- `409` conflito de horário (`"Espaço já está reservado para este horário"`).

---

#### 4.3 Espaços (`/espacos`)
| Método | Endpoint | Papel | Descrição |
|--------|----------|-------|-----------|
| `GET`  | `/espacos` | Qualquer token | Lista espaços com status em tempo real |
| `GET`  | `/espacos/{id}` | Qualquer token | Detalha espaço |
| `POST` | `/espacos` | Admin | Cria espaço (tipos: `visitantes`, `permissionarios`) |
| `PUT`  | `/espacos/{id}` | Admin | Atualiza dados |
| `DELETE` | `/espacos/{id}` | Admin | Remove espaço |
| `PATCH` | `/espacos/{id}/status` | Qualquer token | Atualiza status (`ativo`, `ocupado`, `manutencao`) |
| `PATCH` | `/espacos/{id}/aviso` | Qualquer token | Atualiza aviso |
| `GET` | `/espacos/status` | Qualquer token | Visão consolidada (status atual, próxima reserva, reservas do dia) |

**Status em tempo real**
- `manutencao`: status declarado explicitamente.
- `ocupado`: existe reserva ativa no horário atual.
- `ativo`: disponível.

---

#### 4.4 Usuários (`/users`)
| Método | Endpoint | Papel | Descrição |
|--------|----------|-------|-----------|
| `GET` | `/users` | Admin | Lista usuários |
| `POST` | `/users` | Admin | Cria usuário (admin/recepcionista) |
| `GET` | `/users/{id}` | Admin | Detalha usuário |
| `DELETE` | `/users/{id}` | Admin | Remove usuário |
| `GET` | `/users/me` | Qualquer token | Dados do usuário autenticado |

---

#### 4.5 Strikes (`/strikes`)
| Método | Endpoint | Papel | Descrição |
|--------|----------|-------|-----------|
| `POST` | `/strikes` | Admin ou Recepcionista | Registra strike |
| `GET` | `/strikes/client/{id}` | Admin ou Recepcionista | Lista strikes do cliente |
| `DELETE` | `/strikes/{id}` | Admin | Revoga strike |

---

#### 4.6 Dashboard
- **Endpoint:** `GET /api/v1/dashboard/stats`
- **Resposta:**
```json
{
  "reservasHoje": 12,
  "espacosAtivos": 8,
  "clientesAtivos": 45,
  "taxaOcupacao": 78
}
```
- **Cálculos principais**
  - Reservas do dia atual.
  - Espaços com status `ativo`.
  - Total de clientes.
  - Taxa de ocupação (horas reservadas / horas disponíveis no dia).

---

### 5. Validações e Normalizações

#### 5.1 Clientes
- Normalização automática de CPF e telefone.
- Validação de dígitos verificadores do CPF.
- Verificação de unicidade (CPF, telefone, email).
- Validação de formato de email e telefone.
- Cálculo de `total_reservations` ao retornar dados.

#### 5.2 Reservas
- Conversão de horário (`HH:MM` → `HH:MM:SS`).
- Uso de `time.ParseInLocation` para evitar problemas de timezone.
- Funções auxiliares: `NormalizeTime`, `ParseDate`, `ParseDateTime`, `ValidateFutureDate`.
- Testes cobrindo formatos diferentes e reservas no mesmo dia.

#### 5.3 Espaços
- Tipos permitidos: `visitantes`, `permissionarios` (locação removida).
- Status calculado dinamicamente com base nas reservas do dia.
- Função `toSpaceOutputWithRealTimeStatus` identifica ocupação atual.

---

### 6. Envio Automático de Emails

#### 6.1 Variáveis de ambiente
Veja a seção 2.1 acima. Ative com `EMAIL_ENABLED=true`.

#### 6.2 Modos de operação
- **Desabilitado (`false`)**: não envia emails, apenas loga no console.
- **Habilitado (`true`)**: envia via SMTP e registra sucesso/erro nos logs.

#### 6.3 Eventos suportados
1. **Confirmação de reserva** (ao criar).
2. **Cancelamento de reserva** (ao deletar).

#### 6.4 Templates HTML
- Layout responsivo, com gradientes e informações de data/hora, duração, sala e instruções.
- Cancelamento inclui alerta de segurança.

#### 6.5 Resolução de problemas
- `535 Authentication failed`: use senha de app ou confira credenciais.
- Timeouts: verifique `SMTP_HOST` e `SMTP_PORT`.
- Emails indo para spam: configure SPF/DKIM e peça para marcar como confiável.

#### 6.6 Boas práticas
- Nunca versionar credenciais (`.env` no `.gitignore`).
- Em produção, utilizar secrets (ex.: `kubectl create secret`).

---

### 7. Mensagens de Erro Padronizadas

#### 7.1 Estrutura `APIError`
```json
{
  "status_code": 400,
  "message": "Mensagem amigável",
  "code": "IDENTIFICADOR_INTERNO",
  "field": "campo_relacionado",
  "details": "Sugestão ou informação adicional"
}
```

#### 7.2 Categorias principais
- **Autenticação (401):** token inválido, ausente, credenciais incorretas.
- **Autorização (403):** acesso negado, apenas admin/recepcionista.
- **Não encontrado (404):** cliente, usuário, reserva, espaço, strike.
- **Validação (400):** campos obrigatórios, formatos, CPF/email/telefone inválidos.
- **Conflito (409):** CPF/telefone duplicados, espaço ocupado, strike já existente.
- **Interno (500):** erro genérico ou acesso ao banco.

#### 7.3 Handlers atualizados
- `auth_handlers.go`: login com mensagens claras.
- `client_handlers.go`: mapeamento específico de erros de validação.
- `reservation_handlers.go`: detalhes para data/horário/conflitos.

> Os demais handlers devem seguir o mesmo padrão ao evoluírem.

---

### 8. Histórico de Ajustes Importantes
- **Remoção de locação:** `feature/remove-aluguel`. Apenas `visitantes` e `permissionarios`.
- **Correção de timezone:** reservas utilizam `time.Local`, impedindo falso “passado”.
- **Status em tempo real dos espaços:** leitura dinâmica de reservas ativas.
- **Novas rotas de clientes:** busca por CPF, telefone, ID e criação com validações completas.
- **Emails automáticos:** confirmações/cancelamentos; documentação detalhada.
- **Mensagens de erro:** padronização com `APIError`, códigos e campos específicos.
- **Documentação:** este arquivo centraliza os materiais anteriores.

---

### 9. Testes Automatizados
- **Handlers:** `tests/handlers/*` (clientes, reservas, espaços, dashboard).
- **Utils:** validações e normalizações (CPF, telefone, data, hora).
- **Execução recomendada:** `go test ./...`

---

### 10. Próximos Passos Sugeridos
- Atualizar handlers restantes (espaços, usuários, strikes) para usar `APIError`.
- Adicionar fila/retentativas para envio de email (ex.: Redis/RabbitMQ).
- Implementar email de lembrete, boas-vindas e avisos de strike.
- Padronizar mensagens de sucesso e logs.
- Expandir documentação Swagger conforme novas rotas.

---

### 11. Créditos
- **Desenvolvimento:** Ulpio Paulo de Miranda Netto & Colaboradores
- **Assistente:** AI Codex (OpenAI)
- **Última atualização:** 13/11/2025


