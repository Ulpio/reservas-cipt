# 🔧 Correção: Status em Tempo Real de Espaços

## ❌ Problema Identificado

O endpoint `GET /api/espacos` estava retornando o status **estático do banco de dados**, não verificando se havia reservas ativas no momento.

### Exemplo do problema:
- Sala "Terreo 3" tinha reserva ativa das 14:00 às 16:00
- Hora atual: 15:00 (dentro da reserva)
- Endpoint retornava: `"status": "ativo"` ❌
- Deveria retornar: `"status": "ocupado"` ✅

---

## ✅ Solução Implementada

Modificamos o serviço `GetAllSpaces()` para **calcular o status em tempo real** baseado nas reservas ativas.

### Lógica Implementada:

```
1. Se espaço.status == "manutencao" → retorna "manutencao"
2. Busca reservas de hoje para o espaço
3. Para cada reserva:
   - Calcula end_time = start_time + duration_hours
   - Se NOW está entre start_time e end_time → retorna "ocupado"
4. Se não está ocupado → retorna "ativo"
```

---

## 📝 Alterações no Código

### Arquivo: `services/space_services.go`

#### Antes:
```go
func GetAllSpaces() ([]dto.SpaceOutputDTO, error) {
    var spaces []models.Space
    if err := database.DB.Find(&spaces).Error; err != nil {
        return nil, err
    }
    var output []dto.SpaceOutputDTO
    for _, s := range spaces {
        output = append(output, toSpaceOutput(s))  // ❌ Status estático
    }
    return output, nil
}
```

#### Depois:
```go
func GetAllSpaces() ([]dto.SpaceOutputDTO, error) {
    var spaces []models.Space
    if err := database.DB.Find(&spaces).Error; err != nil {
        return nil, err
    }
    var output []dto.SpaceOutputDTO
    for _, s := range spaces {
        output = append(output, toSpaceOutputWithRealTimeStatus(s))  // ✅ Status dinâmico
    }
    return output, nil
}
```

#### Nova Função:
```go
func toSpaceOutputWithRealTimeStatus(space models.Space) dto.SpaceOutputDTO {
    // 1. Se está em manutenção, manter esse status
    if space.Status == "manutencao" || space.Status == "manutenção" {
        return dto.SpaceOutputDTO{
            ID:       space.ID,
            Name:     space.Name,
            Type:     space.Type,
            Status:   "manutencao",
            Notice:   space.Notice,
            Capacity: int(space.Capacity),
        }
    }

    // 2. Buscar reservas de hoje
    now := time.Now()
    today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    endOfDay := today.Add(24 * time.Hour)

    var activeReservations []models.Reservation
    database.DB.Where("space_id = ? AND date >= ? AND date < ?", space.ID, today, endOfDay).Find(&activeReservations)

    // 3. Verificar se alguma reserva está ativa AGORA
    for _, reservation := range activeReservations {
        endTime := reservation.StartTime.Add(time.Duration(reservation.DurationHours) * time.Hour)
        
        if now.After(reservation.StartTime) && now.Before(endTime) {
            return dto.SpaceOutputDTO{
                ID:       space.ID,
                Name:     space.Name,
                Type:     space.Type,
                Status:   "ocupado",  // ✅ Ocupado em tempo real
                Notice:   space.Notice,
                Capacity: int(space.Capacity),
            }
        }
    }

    // 4. Se não está ocupado, está disponível
    return dto.SpaceOutputDTO{
        ID:       space.ID,
        Name:     space.Name,
        Type:     space.Type,
        Status:   "ativo",
        Notice:   space.Notice,
        Capacity: int(space.Capacity),
    }
}
```

---

## 🎯 Comportamento Atual

### Cenário 1: Sala com reserva ativa
```
Reserva: 14:00 - 16:00
Hora atual: 15:00

GET /api/espacos
{
  "id": 3,
  "name": "Terreo 3",
  "status": "ocupado",  ✅
  ...
}
```

### Cenário 2: Sala com reserva futura
```
Reserva: 18:00 - 20:00
Hora atual: 15:00

GET /api/espacos
{
  "id": 3,
  "name": "Terreo 3",
  "status": "ativo",  ✅ Disponível agora
  ...
}
```

### Cenário 3: Sala em manutenção
```
Status no banco: "manutencao"

GET /api/espacos
{
  "id": 5,
  "name": "Auditório",
  "status": "manutencao",  ✅
  ...
}
```

### Cenário 4: Sala sem reservas
```
Nenhuma reserva hoje

GET /api/espacos
{
  "id": 1,
  "name": "Sala A",
  "status": "ativo",  ✅
  ...
}
```

---

## 📊 Tabela de Estados

| Condição | Status Retornado |
|----------|------------------|
| Status banco = "manutencao" | `"manutencao"` |
| Reserva ativa AGORA | `"ocupado"` |
| Reserva futura (não agora) | `"ativo"` |
| Sem reservas | `"ativo"` |

---

## 🔍 Verificação de Reserva Ativa

```go
now := time.Now()                    // Ex: 15:30
reservation.StartTime                // Ex: 14:00
endTime = startTime + duration       // Ex: 16:00

if now.After(14:00) && now.Before(16:00) {
    // 15:30 está entre 14:00 e 16:00
    return "ocupado"  ✅
}
```

---

## ⚡ Performance

### Query executada por espaço:
```sql
SELECT * FROM reservations 
WHERE space_id = ? 
  AND date >= '2025-11-13 00:00:00' 
  AND date < '2025-11-14 00:00:00'
```

- **Índice recomendado**: `(space_id, date)` para otimização
- **Custo**: Uma query por espaço
- **Tempo estimado**: ~1-2ms por espaço

### Otimização futura (se necessário):
```sql
-- Buscar todas as reservas de uma vez
SELECT * FROM reservations 
WHERE date >= '2025-11-13 00:00:00' 
  AND date < '2025-11-14 00:00:00'
```

---

## 🧪 Como Testar

### 1. Criar uma reserva para AGORA:
```json
POST /api/reservas
{
  "space_id": 3,
  "date": "2025-11-13",
  "start_time": "14:00",
  "duration_hours": 2
}
```

### 2. Verificar status (entre 14:00 e 16:00):
```json
GET /api/espacos
[
  {
    "id": 3,
    "name": "Terreo 3",
    "status": "ocupado"  ✅
  }
]
```

### 3. Verificar status (fora do horário):
```json
GET /api/espacos (antes das 14:00 ou após 16:00)
[
  {
    "id": 3,
    "name": "Terreo 3",
    "status": "ativo"  ✅
  }
]
```

---

## 🔧 Timezone

O sistema usa **timezone local** (configurado no servidor).

```go
now := time.Now()  // Hora local do servidor
```

Certifique-se de que o servidor está configurado com o timezone correto (UTC-3 para Brasil).

---

## 📝 Observações Importantes

1. **Status "manutencao"**: Tem prioridade sobre reservas. Mesmo com reserva ativa, se o status for "manutencao", retorna "manutencao".

2. **Status "ativo" vs "disponivel"**: O endpoint retorna "ativo" para espaços disponíveis. O frontend pode interpretar como "disponível".

3. **Atualização automática**: O status é recalculado a cada requisição, sempre refletindo o estado atual.

4. **Endpoint /espacos/status**: Continue usando este endpoint para visão detalhada com lista de reservas.

---

## ✅ Status da Correção

- [x] Problema identificado
- [x] Solução implementada
- [x] Código compilado sem erros
- [x] Lógica de verificação de tempo real
- [x] Prioridade de status (manutencao > ocupado > ativo)
- [x] Documentação atualizada

---

**Data**: 13/11/2025  
**Versão**: 1.1.0  
**Status**: ✅ Corrigido e Testado

