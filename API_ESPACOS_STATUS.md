# 🏢 API - Visualização Geral de Espaços

## Endpoint: `/api/v1/espacos/status`

### Método: `GET`

### Descrição:
Este endpoint fornece uma visão geral do status de todos os espaços, incluindo informações sobre ocupação atual, próximas reservas e histórico de reservas do dia.

---

## Autenticação

**Obrigatório:** JWT Token

```
Authorization: Bearer {token}
```

---

## Resposta de Sucesso (200)

```json
[
  {
    "id": 1,
    "nome": "Sala de Reunião A",
    "status": "disponivel",
    "reservaAtual": null,
    "proximaReserva": {
      "id": 124,
      "cliente": "Maria Santos",
      "horaInicio": "2024-01-15T17:00:00Z",
      "horaFim": "2024-01-15T19:00:00Z"
    },
    "reservasHoje": [
      {
        "id": 122,
        "cliente": "Pedro Costa",
        "horaInicio": "2024-01-15T09:00:00Z",
        "horaFim": "2024-01-15T11:00:00Z"
      },
      {
        "id": 124,
        "cliente": "Maria Santos",
        "horaInicio": "2024-01-15T17:00:00Z",
        "horaFim": "2024-01-15T19:00:00Z"
      }
    ]
  },
  {
    "id": 2,
    "nome": "Sala de Reunião B",
    "status": "ocupado",
    "reservaAtual": {
      "id": 125,
      "cliente": "Ana Paula",
      "horaInicio": "2024-01-15T15:00:00Z",
      "horaFim": "2024-01-15T18:00:00Z"
    },
    "proximaReserva": null,
    "reservasHoje": [
      {
        "id": 125,
        "cliente": "Ana Paula",
        "horaInicio": "2024-01-15T15:00:00Z",
        "horaFim": "2024-01-15T18:00:00Z"
      }
    ]
  },
  {
    "id": 3,
    "nome": "Auditório",
    "status": "manutencao",
    "reservaAtual": null,
    "proximaReserva": null,
    "reservasHoje": []
  }
]
```

---

## Campos da Resposta

### Espaço (SpaceStatusDTO)

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | number | ID único do espaço |
| `nome` | string | Nome do espaço |
| `status` | string | Status atual: "disponivel", "ocupado" ou "manutencao" |
| `reservaAtual` | object\|null | Dados da reserva em andamento (apenas se ocupado) |
| `proximaReserva` | object\|null | Próxima reserva do dia (se houver) |
| `reservasHoje` | array | Lista de todas as reservas do dia atual |

### Reserva (ReservaStatusDTO)

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | number | ID da reserva |
| `cliente` | string | Nome do cliente que fez a reserva |
| `horaInicio` | string | Data/hora de início (ISO 8601) |
| `horaFim` | string | Data/hora de término (ISO 8601) |

---

## Lógica de Negócio para Determinar Status

### 1. **"ocupado"**
- Existe uma reserva ativa no momento da consulta
- Hora atual está entre `horaInicio` e `horaFim`
- `reservaAtual` contém os dados da reserva em andamento

### 2. **"manutencao"**
- O espaço está marcado como indisponível no sistema
- Não permite reservas até ser reativado
- Não retorna informações de reservas

### 3. **"disponivel"**
- Não há reserva ativa no momento
- O espaço está disponível para reservas
- Se houver reservas futuras no dia, `proximaReserva` contém a próxima

---

## Exemplos de Casos de Uso

### Caso 1: Sala Disponível sem Reservas

```json
{
  "id": 5,
  "nome": "Sala Privativa 1",
  "status": "disponivel",
  "reservaAtual": null,
  "proximaReserva": null,
  "reservasHoje": []
}
```

### Caso 2: Sala Disponível com Próxima Reserva

```json
{
  "id": 6,
  "nome": "Sala Privativa 2",
  "status": "disponivel",
  "reservaAtual": null,
  "proximaReserva": {
    "id": 150,
    "cliente": "Carlos Eduardo",
    "horaInicio": "2024-01-15T20:00:00Z",
    "horaFim": "2024-01-15T22:00:00Z"
  },
  "reservasHoje": [
    {
      "id": 150,
      "cliente": "Carlos Eduardo",
      "horaInicio": "2024-01-15T20:00:00Z",
      "horaFim": "2024-01-15T22:00:00Z"
    }
  ]
}
```

### Caso 3: Sala Ocupada

```json
{
  "id": 7,
  "nome": "Sala de Treinamento",
  "status": "ocupado",
  "reservaAtual": {
    "id": 160,
    "cliente": "Fernanda Lima",
    "horaInicio": "2024-01-15T13:00:00Z",
    "horaFim": "2024-01-15T17:00:00Z"
  },
  "proximaReserva": {
    "id": 161,
    "cliente": "Roberto Alves",
    "horaInicio": "2024-01-15T18:00:00Z",
    "horaFim": "2024-01-15T20:00:00Z"
  },
  "reservasHoje": [
    {
      "id": 160,
      "cliente": "Fernanda Lima",
      "horaInicio": "2024-01-15T13:00:00Z",
      "horaFim": "2024-01-15T17:00:00Z"
    },
    {
      "id": 161,
      "cliente": "Roberto Alves",
      "horaInicio": "2024-01-15T18:00:00Z",
      "horaFim": "2024-01-15T20:00:00Z"
    }
  ]
}
```

---

## Exemplo de Uso

### cURL

```bash
curl -X GET "http://localhost:8080/api/v1/espacos/status" \
  -H "Authorization: Bearer {token}"
```

### JavaScript/Frontend

```javascript
import api from './services/api';

// Buscar status de todos os espaços
const fetchSpacesStatus = async () => {
  try {
    const { data } = await api.get('/espacos/status');
    console.log('Status dos espaços:', data);
    
    // Filtrar espaços disponíveis
    const disponíveis = data.filter(s => s.status === 'disponivel');
    
    // Filtrar espaços ocupados
    const ocupados = data.filter(s => s.status === 'ocupado');
    
    return data;
  } catch (error) {
    console.error('Erro ao buscar status:', error);
  }
};

// Atualizar automaticamente a cada 30 segundos
setInterval(fetchSpacesStatus, 30000);
```

---

## Status Codes

| Code | Descrição |
|------|-----------|
| `200` | Sucesso - retorna array de espaços |
| `401` | Não autenticado - token inválido ou ausente |
| `500` | Erro interno do servidor |

---

## Implementação

### Arquivos Criados:

1. **`dto/space_status_dto.go`** - DTOs para resposta
   - `ReservaStatusDTO`
   - `SpaceStatusDTO`

2. **`services/space_status_services.go`** - Lógica de negócio
   - `GetAllSpacesStatus()`
   - `determineSpaceStatus()`
   - `calculateEndTime()`

3. **`tests/handlers/space_status_handler_test.go`** - Testes

### Arquivos Modificados:

1. **`handlers/space_handlers.go`**
   - Adicionado `GetAllSpacesStatusHandler()`

2. **`routes/space_routes.go`**
   - Adicionada rota `GET /espacos/status`

---

## Testes

Execute os testes com:

```bash
go test ./tests/handlers/space_status_handler_test.go -v
```

### Casos de Teste Implementados:

1. ✅ Espaço disponível sem reservas
2. ✅ Espaço disponível com próxima reserva
3. ✅ Espaço ocupado com reserva atual e próxima
4. ✅ Espaço em manutenção
5. ✅ Ordenação de reservas por horário
6. ✅ Cálculo correto de hora de término
7. ✅ Retorna 401 sem autenticação

---

## Otimizações Implementadas

### 1. **Query Eficiente com JOIN**
```go
database.DB.Table("reservations").
    Select("reservations.id, clients.name, reservations.start_time, reservations.duration_hours").
    Joins("JOIN clients ON clients.id = reservations.client_id").
    Where("reservations.space_id = ? AND reservations.date >= ? AND reservations.date < ?", spaceID, startOfDay, endOfDay).
    Order("reservations.start_time ASC").
    Scan(&results)
```

### 2. **Cálculo de Hoje com Timezone**
- Usa timezone do servidor
- Considera início e fim do dia corretamente

### 3. **Estrutura Eficiente**
- Uma única query por espaço
- JOIN para evitar N+1 queries
- Ordenação no banco de dados

---

## Observações Importantes

### Timezone
- A API considera "hoje" baseado no timezone do servidor
- Formato ISO 8601 para datas (`2024-01-15T14:00:00Z`)

### Ordenação
- `reservasHoje` sempre ordenado por `horaInicio` (ASC)

### Cálculo de Status
- Status é calculado em tempo real a cada requisição
- Considera hora exata da consulta

### Performance
- Recomenda-se atualização a cada 30 segundos no frontend
- Considerar cache para sistemas com muitos espaços

---

## Sugestões de Melhorias Futuras

1. **Cache**: Implementar cache com TTL de 30 segundos
2. **WebSocket**: Atualização em tempo real via WebSocket
3. **Filtros**: Permitir filtrar por status ou tipo de espaço
4. **Paginação**: Para sistemas com muitos espaços
5. **Histórico**: Incluir estatísticas de ocupação

---

## Integração com Frontend

### Atualização Automática

```javascript
// No componente React
useEffect(() => {
  const fetchStatus = async () => {
    const { data } = await api.get('/espacos/status');
    setSpaces(data);
  };

  fetchStatus(); // Buscar imediatamente
  const interval = setInterval(fetchStatus, 30000); // A cada 30s

  return () => clearInterval(interval); // Cleanup
}, []);
```

### Renderização

```jsx
{spaces.map(space => (
  <SpaceCard
    key={space.id}
    name={space.nome}
    status={space.status}
    currentReservation={space.reservaAtual}
    nextReservation={space.proximaReserva}
    todayReservations={space.reservasHoje}
  />
))}
```

---

## Swagger/OpenAPI

Documentação Swagger disponível em:
```
http://localhost:8080/swagger/index.html
```

---

**Desenvolvido por:** Ulpio Paulo de Miranda Netto  
**Data:** 2025  
**Versão API:** 1.0

