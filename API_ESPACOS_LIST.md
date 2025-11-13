# 📋 API de Espaços - GET /api/espacos

## ✅ Status: IMPLEMENTADO

Endpoint completo para listar todos os espaços cadastrados no sistema.

---

## 📍 Endpoint

```
GET /api/espacos
```

### Headers:
```http
Authorization: Bearer {token}
```

---

## 📊 Resposta de Sucesso (200 OK)

```json
[
  {
    "id": 1,
    "name": "Sala de Reunião A",
    "type": "visitantes",
    "capacity": 10,
    "status": "ativo",
    "notice": "Equipada com projetor"
  },
  {
    "id": 2,
    "name": "Sala de Reunião B",
    "type": "permissionarios",
    "capacity": 6,
    "status": "ocupado",
    "notice": ""
  },
  {
    "id": 3,
    "name": "Auditório",
    "type": "locacao",
    "capacity": 50,
    "status": "manutencao",
    "notice": "Ar-condicionado em manutenção"
  }
]
```

---

## 📋 Campos Retornados

| Campo | Tipo | Descrição | Valores Possíveis |
|-------|------|-----------|-------------------|
| `id` | `number` | ID único do espaço | - |
| `name` | `string` | Nome do espaço | - |
| `type` | `string` | Tipo de sala de reunião | `visitantes`, `permissionarios` |
| `capacity` | `number` | Capacidade máxima | Número inteiro |
| `status` | `string` | Status atual | `ativo`, `ocupado`, `manutencao` |
| `notice` | `string` | Aviso/observação | Texto ou vazio |

---

## 🎨 Exibição no Frontend

### Select de Espaços (Nova Reserva):

```
Sala de Reunião A — 🟢 Disponível
Sala de Reunião B — 🔴 Ocupado
Auditório — 🟡 Em Manutenção
```

### Mapeamento de Status:

| Status API | Ícone | Label | Pode Reservar? |
|------------|-------|-------|----------------|
| `ativo` ou `disponivel` | 🟢 | Disponível | ✅ Sim |
| `ocupado` | 🔴 | Ocupado | ⚠️ Depende* |
| `manutencao` | 🟡 | Em Manutenção | ❌ Não |

\* _Pode reservar para horário futuro_

### Legenda:
```
🟢 Disponível | 🔴 Ocupado | 🟡 Em Manutenção
```

---

## 💻 Implementação Backend

### Handler (`handlers/space_handlers.go`):

```go
// GetAllSpacesHandler lista todos os espaços cadastrados.
// @Summary Lista espaços
// @Description Retorna todos os espaços registrados.
// @Tags espacos
// @Produce json
// @Success 200 {array} dto.SpaceOutputDTO
// @Failure 500 {object} dto.ErrorResponse
// @Router /espacos [get]
func GetAllSpacesHandler(c *gin.Context) {
    spaces, err := services.GetAllSpaces()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        return
    }
    c.JSON(http.StatusOK, spaces)
}
```

### Service (`services/space_services.go`):

```go
func GetAllSpaces() ([]dto.SpaceOutputDTO, error) {
    var spaces []models.Space
    if err := database.DB.Find(&spaces).Error; err != nil {
        return nil, err
    }
    var output []dto.SpaceOutputDTO
    for _, s := range spaces {
        output = append(output, toSpaceOutput(s))
    }
    return output, nil
}

func toSpaceOutput(space models.Space) dto.SpaceOutputDTO {
    return dto.SpaceOutputDTO{
        ID:       space.ID,
        Name:     space.Name,
        Type:     space.Type,
        Status:   space.Status,
        Notice:   space.Notice,
        Capacity: int(space.Capacity),
    }
}
```

### Route (`routes/space_routes.go`):

```go
func SpaceRoutes(r *gin.RouterGroup) {
    spaceGroup := r.Group("/espacos")
    spaceGroup.Use(middleware.JWTAuthMiddleware())
    
    spaceGroup.GET("", handlers.GetAllSpacesHandler)
    // ... outras rotas
}
```

### DTO (`dto/espaco_dto.go`):

```go
type SpaceOutputDTO struct {
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Type     string `json:"type"`
    Status   string `json:"status"`
    Notice   string `json:"notice"`
    Capacity int    `json:"capacity"`
}
```

---

## 🧪 Exemplo de Uso

### Frontend (JavaScript):

```javascript
// Buscar todos os espaços
const response = await api.get('/espacos');

// Response: 200 OK
const espacos = response.data;

// Filtrar apenas espaços ativos
const disponíveis = espacos.filter(e => e.status === 'ativo');

// Montar select
espacos.forEach(espaco => {
  const icon = getStatusIcon(espaco.status);
  const label = getStatusLabel(espaco.status);
  
  console.log(`${espaco.name} — ${icon} ${label}`);
});

function getStatusIcon(status) {
  switch(status) {
    case 'ativo':
    case 'disponivel':
      return '🟢';
    case 'ocupado':
      return '🔴';
    case 'manutencao':
      return '🟡';
    default:
      return '⚪';
  }
}

function getStatusLabel(status) {
  switch(status) {
    case 'ativo':
    case 'disponivel':
      return 'Disponível';
    case 'ocupado':
      return 'Ocupado';
    case 'manutencao':
      return 'Em Manutenção';
    default:
      return 'Desconhecido';
  }
}
```

---

## 🔐 Autenticação

- **Obrigatória**: Sim
- **Middleware**: `JWTAuthMiddleware()`
- **Roles aceitas**: Todas (admin, recepcionista, locacao)

---

## 📊 Status Codes

| Code | Descrição |
|------|-----------|
| `200` | Sucesso - retorna array de espaços |
| `401` | Não autenticado - token inválido |
| `500` | Erro interno do servidor |

---

## 🔗 Endpoints Relacionados

- `GET /api/espacos/status` - Status detalhado com reservas
- `GET /api/espacos/:id` - Buscar espaço específico
- `POST /api/espacos` - Criar novo espaço (admin)
- `PUT /api/espacos/:id` - Atualizar espaço (admin)
- `PATCH /api/espacos/:id/status` - Atualizar status
- `DELETE /api/espacos/:id` - Deletar espaço (admin)

---

## ⚙️ Diferença entre GET /espacos e GET /espacos/status

| Endpoint | Uso | Retorna |
|----------|-----|---------|
| `GET /espacos` | **Lista simples** para selects | Dados básicos (id, nome, status) |
| `GET /espacos/status` | **Dashboard/Monitor** | Status + reserva atual + próximas reservas |

---

## ✅ Checklist de Implementação

- [x] Handler implementado
- [x] Service implementado
- [x] Rota registrada
- [x] DTO completo
- [x] Middleware de autenticação
- [x] Swagger annotations
- [x] Documentação

---

## 📝 Observações

1. **Status "ocupado"**: Calculado em tempo real pelo endpoint `/espacos/status`
2. **Status "manutencao"**: Definido manualmente pelo admin
3. **Status "ativo"**: Espaço operacional e disponível
4. **Frontend compatível**: Aceita "manutencao" e "manutenção"
5. **Ordenação**: Por ID (padrão do banco)

---

**Versão**: 1.0.0  
**Status**: ✅ Implementado  
**Última atualização**: 13/11/2025

