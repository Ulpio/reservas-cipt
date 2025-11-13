# 📊 API - Dashboard Statistics

## Endpoint para Estatísticas do Dashboard

### **GET** `/api/v1/dashboard/stats`

Retorna as estatísticas principais para exibição no dashboard.

---

## Autenticação

Este endpoint requer autenticação via JWT Token.

**Header:**
```
Authorization: Bearer {token}
```

---

## Resposta Esperada

```json
{
  "reservasHoje": 12,
  "espacosAtivos": 8,
  "clientesAtivos": 45,
  "taxaOcupacao": 78
}
```

---

## Campos da Resposta

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `reservasHoje` | `int64` | Total de reservas agendadas para o dia atual |
| `espacosAtivos` | `int64` | Quantidade de espaços com status "ativo" |
| `clientesAtivos` | `int64` | Total de clientes cadastrados no sistema |
| `taxaOcupacao` | `int` | Percentual de ocupação dos espaços (0-100) |

---

## Implementação

### Arquivos Criados/Modificados

1. **`dto/dashboard_dto.go`** - DTO para resposta
2. **`services/dashboard_services.go`** - Lógica de cálculo das estatísticas
3. **`handlers/dashboard_handlers.go`** - Handler do endpoint
4. **`routes/dashboard_routes.go`** - Configuração de rotas
5. **`routes/routes.go`** - Registro da rota no sistema
6. **`tests/handlers/dashboard_handler_test.go`** - Testes do endpoint

### Cálculos Implementados

#### 1. **reservasHoje**
Conta todas as reservas cuja data é igual à data atual.

```go
database.DB.Model(&models.Reservation{}).
    Where("date = ?", today).
    Count(&reservasHoje)
```

#### 2. **espacosAtivos**
Conta todos os espaços com status "ativo".

```go
database.DB.Model(&models.Space{}).
    Where("status = ?", "ativo").
    Count(&espacosAtivos)
```

#### 3. **clientesAtivos**
Conta todos os clientes cadastrados.

```go
database.DB.Model(&models.Client{}).
    Count(&clientesAtivos)
```

#### 4. **taxaOcupacao**
Calcula o percentual de ocupação baseado nas horas reservadas hoje vs. horas disponíveis.

**Fórmula:**
```
Taxa = (Total horas reservadas hoje / Total horas disponíveis hoje) * 100
Onde: Horas disponíveis = Espaços ativos * 24 horas
```

```go
database.DB.Model(&models.Reservation{}).
    Select("COALESCE(SUM(duration_hours), 0) as total").
    Where("date = ?", today).
    Scan(&result)

taxaOcupacao = int((float64(totalHorasReservadas) / float64(horasDisponiveis)) * 100)
```

---

## Exemplo de Uso

### cURL

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/stats" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### JavaScript (Frontend)

```javascript
const { data } = await api.get('/dashboard/stats');
console.log(data);
// {
//   reservasHoje: 12,
//   espacosAtivos: 8,
//   clientesAtivos: 45,
//   taxaOcupacao: 78
// }
```

---

## Status Codes

| Code | Descrição |
|------|-----------|
| `200` | Sucesso - retorna os dados |
| `401` | Não autorizado - token inválido ou ausente |
| `500` | Erro interno do servidor |

---

## Integração com Frontend

O frontend já está configurado para consumir este endpoint automaticamente em:

**`src/pages/Dashboard.jsx`**

```javascript
useEffect(() => {
  async function fetchStats() {
    try {
      const { data } = await api.get('/dashboard/stats');
      setStats(data);
      setLoading(false);
    } catch (error) {
      console.error('Erro ao carregar estatísticas:', error);
      setLoading(false);
    }
  }

  fetchStats();
}, []);
```

---

## Testes

Execute os testes com:

```bash
go test ./tests/handlers/dashboard_handler_test.go -v
```

### Casos de Teste

1. ✅ Retorna estatísticas corretas com dados válidos
2. ✅ Retorna 401 sem token de autenticação
3. ✅ Calcula corretamente a taxa de ocupação
4. ✅ Contabiliza apenas espaços ativos
5. ✅ Contabiliza apenas reservas do dia atual

---

## Performance

### Otimizações Implementadas

1. **Queries eficientes**: Uso de `COUNT()` direto no banco
2. **Single query para soma**: `COALESCE(SUM())` para total de horas
3. **Índices**: As tabelas já possuem índices nas colunas usadas (date, status)

### Sugestões Futuras

1. **Cache**: Implementar cache com TTL de 5-10 minutos
2. **Redis**: Armazenar estatísticas em cache distribuído
3. **WebSocket**: Atualização em tempo real
4. **Agregações**: Pre-calcular estatísticas em background job

---

## Considerações de Segurança

1. ✅ Endpoint protegido por JWT
2. ✅ Validação de token obrigatória
3. ✅ Sem exposição de dados sensíveis
4. ✅ Queries parametrizadas (SQL Injection protection)

---

## Swagger/OpenAPI

A documentação Swagger foi atualizada automaticamente com as anotações do handler:

```go
// @Summary Obter estatísticas do dashboard
// @Description Retorna estatísticas principais para exibição no dashboard
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.DashboardStatsDTO
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /dashboard/stats [get]
```

Acesse: `http://localhost:8080/swagger/index.html`

---

## Logs e Monitoramento

Para adicionar logs detalhados, modifique o service:

```go
log.Printf("Dashboard Stats - Reservas: %d, Espaços: %d, Clientes: %d, Taxa: %d%%",
    stats.ReservasHoje, stats.EspacosAtivos, stats.ClientesAtivos, stats.TaxaOcupacao)
```

---

## Troubleshooting

### Problema: Taxa de ocupação sempre 0

**Solução:** Verifique se há espaços com status "ativo" e reservas para hoje.

### Problema: 401 Unauthorized

**Solução:** Certifique-se de incluir o header `Authorization: Bearer {token}`.

### Problema: Valores incorretos

**Solução:** Verifique o timezone do servidor e do banco de dados.

---

**Desenvolvido por:** Ulpio Paulo de Miranda Netto  
**Data:** 2025  
**Versão API:** 1.0

