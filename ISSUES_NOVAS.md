# 📋 Issues Novas do Projeto - CIPT

**Data**: 13/11/2025

---

## 1️⃣ Issue 1: Remover Funcionalidade de Aluguel de Espaços

### Objetivo:
Sistema deve focar apenas em **reservas de salas de reunião**, removendo a parte de aluguel/locação.

### Mudanças Necessárias:

#### 1.1 Modelo de Espaços
- ❌ Remover tipo `"locacao"` 
- ✅ Manter apenas: `"visitantes"` e `"permissionarios"`

#### 1.2 Validações
```go
// Antes: aceita 3 tipos
types := []string{"visitantes", "permissionarios", "locacao"}

// Depois: aceita 2 tipos
types := []string{"visitantes", "permissionarios"}
```

#### 1.3 Frontend
- Remover opção "Locação" dos formulários
- Atualizar labels para "Salas de Reunião"

#### 1.4 Documentação
- Atualizar todas as referências de "aluguel" para "reserva"
- Remover menções a "locação"

### Impacto:
- **Baixo**: Apenas remove uma opção de tipo
- **Breaking Change**: ❌ Não (compatível com dados existentes)

---

## 2️⃣ Issue 2: Envio de Email Automático ao Cliente

### Objetivo:
Enviar emails automáticos aos clientes quando:
1. ✅ Reserva criada (confirmação)
2. ✅ Reserva cancelada
3. ✅ Lembrete 24h antes (opcional)
4. ✅ Strike aplicado (opcional)

### Implementação Sugerida:

#### 2.1 Configuração SMTP
```go
// config/email.go
type EmailConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string
}
```

#### 2.2 Templates de Email

**Email de Confirmação:**
```
Assunto: Reserva Confirmada - Sala {nome_sala}

Olá {nome_cliente},

Sua reserva foi confirmada com sucesso!

📅 Data: {data}
🕐 Horário: {hora_inicio} às {hora_fim}
🏢 Sala: {nome_sala}
👤 Capacidade: {capacidade} pessoas

Obrigado por utilizar o CIPT.
```

**Email de Cancelamento:**
```
Assunto: Reserva Cancelada - Sala {nome_sala}

Olá {nome_cliente},

Sua reserva foi cancelada.

📅 Data: {data}
🕐 Horário: {hora_inicio} às {hora_fim}
🏢 Sala: {nome_sala}

Se você não solicitou este cancelamento, entre em contato conosco.
```

#### 2.3 Service de Email
```go
// services/email_service.go
func SendReservationConfirmation(reservation Reservation, client Client) error
func SendReservationCancellation(reservation Reservation, client Client) error
func SendReminder24h(reservation Reservation, client Client) error
```

#### 2.4 Integração
```go
// Em CreateReservation
reservation := services.CreateReservation(input)
go services.SendReservationConfirmation(reservation, client)  // Async
```

### Bibliotecas Recomendadas:
1. **gomail** - github.com/go-gomail/gomail
2. **email** - github.com/jordan-wright/email
3. **smtp nativo** - net/smtp

### Variáveis de Ambiente:
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=seu-email@gmail.com
SMTP_PASSWORD=sua-senha-app
EMAIL_FROM=noreply@cipt.com
```

### Impacto:
- **Médio**: Requer configuração SMTP
- **Breaking Change**: ❌ Não (funcionalidade nova)

---

## 3️⃣ Issue 3: Melhorar Mensagens de Erro

### Objetivo:
Padronizar e tornar mensagens de erro mais claras e descritivas para o frontend.

### Problemas Atuais:

#### Antes (❌):
```json
{
  "error": "Erro ao criar reserva"
}
```

#### Depois (✅):
```json
{
  "error": "Não foi possível criar a reserva",
  "message": "O espaço já está reservado para o horário selecionado (14:00 - 16:00)",
  "code": "SPACE_ALREADY_BOOKED",
  "statusCode": 409
}
```

### Estrutura Padronizada:

```go
type ErrorResponse struct {
    Error      string `json:"error"`       // Mensagem curta
    Message    string `json:"message"`     // Mensagem detalhada
    Code       string `json:"code"`        // Código do erro
    StatusCode int    `json:"statusCode"`  // HTTP status
    Details    string `json:"details,omitempty"` // Detalhes técnicos (dev only)
}
```

### Códigos de Erro Padronizados:

| Código | Descrição |
|--------|-----------|
| `VALIDATION_ERROR` | Erro de validação de dados |
| `NOT_FOUND` | Recurso não encontrado |
| `ALREADY_EXISTS` | Recurso já existe (duplicado) |
| `SPACE_ALREADY_BOOKED` | Espaço já reservado |
| `INVALID_CPF` | CPF inválido |
| `UNDERAGE_CLIENT` | Cliente menor de idade |
| `SPACE_UNAVAILABLE` | Espaço em manutenção |
| `UNAUTHORIZED` | Não autenticado |
| `FORBIDDEN` | Sem permissão |
| `INTERNAL_ERROR` | Erro interno |

### Exemplos de Mensagens Melhoradas:

#### CPF Inválido:
```json
{
  "error": "CPF inválido",
  "message": "O CPF informado (111.111.111-11) não é válido. Verifique os dígitos verificadores.",
  "code": "INVALID_CPF",
  "statusCode": 400
}
```

#### Cliente Menor de Idade:
```json
{
  "error": "Cliente não atende aos requisitos",
  "message": "O cliente deve ter pelo menos 18 anos. Data de nascimento informada: 2010-05-15 (15 anos).",
  "code": "UNDERAGE_CLIENT",
  "statusCode": 400
}
```

#### Conflito de Horário:
```json
{
  "error": "Conflito de horário",
  "message": "A sala 'Terreo 3' já possui uma reserva para o dia 13/11/2025 das 14:00 às 16:00.",
  "code": "SPACE_ALREADY_BOOKED",
  "statusCode": 409
}
```

### Implementação:

```go
// utils/errors.go
package utils

type AppError struct {
    Error      string
    Message    string
    Code       string
    StatusCode int
}

func NewValidationError(message string) AppError {
    return AppError{
        Error:      "Dados inválidos",
        Message:    message,
        Code:       "VALIDATION_ERROR",
        StatusCode: 400,
    }
}

func NewNotFoundError(resource string) AppError {
    return AppError{
        Error:      fmt.Sprintf("%s não encontrado", resource),
        Message:    fmt.Sprintf("O %s solicitado não foi encontrado no sistema.", resource),
        Code:       "NOT_FOUND",
        StatusCode: 404,
    }
}
```

### Impacto:
- **Médio**: Requer refatoração de handlers
- **Breaking Change**: ⚠️ Sim (estrutura de resposta muda)

---

## 📊 Priorização Sugerida

| Issue | Prioridade | Esforço | Impacto |
|-------|------------|---------|---------|
| **1. Remover Aluguel** | 🔴 Alta | Baixo (2h) | Baixo |
| **3. Melhorar Erros** | 🟡 Média | Médio (4h) | Alto |
| **2. Email Automático** | 🟢 Baixa | Alto (8h) | Médio |

### Ordem Recomendada:
1. ✅ **Issue 1** - Remover aluguel (rápido e simples)
2. ✅ **Issue 3** - Melhorar erros (melhora UX)
3. ✅ **Issue 2** - Email automático (mais complexo)

---

## 🚀 Plano de Ação

### Sprint 1 (Hoje):
- [x] Documentar issues
- [ ] Implementar Issue 1 (Remover aluguel)
- [ ] Implementar Issue 3 (Melhorar erros)

### Sprint 2 (Próxima):
- [ ] Configurar SMTP
- [ ] Implementar Issue 2 (Email automático)
- [ ] Testar integração completa

---

**Qual issue você quer que eu comece a implementar primeiro?**

1️⃣ Remover funcionalidade de aluguel (mais rápido)  
2️⃣ Envio de email automático (mais complexo)  
3️⃣ Melhorar mensagens de erro (melhora UX)

