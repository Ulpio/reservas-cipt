# 🔧 Issue 1: Remover Funcionalidade de Aluguel

**Branch**: `feature/remove-aluguel`  
**Data**: 13/11/2025

---

## 📋 Análise Inicial

### Contextos Identificados:

#### 1. **Tipo de Espaço** (Space.Type)
Valores atuais: `visitantes`, `permissionarios`, `locacao`

**Arquivos afetados:**
- `models/space.go` - Comentário do campo Type
- Documentações (API_ESPACOS_LIST.md, etc)

#### 2. **Role de Usuário** (User.Role)  
Valores: `admin`, `recepcionista`, `locacao`

**Arquivos afetados:**
- `services/reservation_services.go` - Validação de permissão

---

## 🎯 Decisão

### ✅ O QUE REMOVER:
- **Tipo de espaço "locacao"** - Sistema focará apenas em salas de reunião

### ⚠️ O QUE MANTER:
- **Role "locacao"** - Usuários com essa role continuam existindo, mas gerenciam apenas salas de reunião

---

## 📝 Mudanças a Fazer

### 1. Atualizar Modelo de Espaço
```go
// models/space.go
Type string `gorm:"not null"` // visitantes, permissionarios
```

### 2. Adicionar Validação de Tipo
```go
// services/space_services.go
func ValidateSpaceType(spaceType string) error {
    validTypes := []string{"visitantes", "permissionarios"}
    for _, valid := range validTypes {
        if spaceType == valid {
            return nil
        }
    }
    return errors.New("tipo de espaço inválido. Use 'visitantes' ou 'permissionarios'")
}
```

### 3. Atualizar CreateSpace
```go
func CreateSpace(input dto.CreateSpaceDTO) (dto.SpaceOutputDTO, error) {
    // Validar tipo
    if err := ValidateSpaceType(input.Type); err != nil {
        return dto.SpaceOutputDTO{}, err
    }
    // ... resto do código
}
```

### 4. Atualizar UpdateSpace
```go
func UpdateSpace(id uint, input dto.UpdateSpaceDTO) (dto.SpaceOutputDTO, error) {
    // Validar tipo
    if err := ValidateSpaceType(input.Type); err != nil {
        return dto.SpaceOutputDTO{}, err
    }
    // ... resto do código
}
```

### 5. Atualizar Documentações
- API_ESPACOS_LIST.md
- ISSUES_NOVAS.md
- Swagger annotations

---

## ⚠️ Compatibilidade

### Dados Existentes:
Se já existirem espaços com `type = "locacao"` no banco:

**Opção 1 - Migration (Recomendado):**
```sql
-- Converter espaços de locação para visitantes
UPDATE spaces 
SET type = 'visitantes' 
WHERE type = 'locacao';
```

**Opção 2 - Manter:**
- Espaços existentes com "locacao" continuam funcionando
- Novos espaços não podem ter esse tipo

---

## ✅ Checklist de Implementação

- [ ] Atualizar comentário em `models/space.go`
- [ ] Criar função `ValidateSpaceType()` em `services/space_services.go`
- [ ] Adicionar validação em `CreateSpace()`
- [ ] Adicionar validação em `UpdateSpace()`
- [ ] Atualizar documentações
- [ ] Atualizar Swagger annotations
- [ ] (Opcional) Criar migration para converter dados existentes
- [ ] Testar criação de espaço com tipo inválido
- [ ] Testar criação de espaço com tipos válidos
- [ ] Commit e push

---

## 🧪 Testes

### Teste 1: Criar espaço com tipo válido
```json
POST /api/espacos
{
  "name": "Sala A",
  "type": "visitantes",  ✅
  "capacity": 10,
  "status": "ativo"
}
// Esperado: 201 Created
```

### Teste 2: Criar espaço com tipo inválido
```json
POST /api/espacos
{
  "name": "Sala B",
  "type": "locacao",  ❌
  "capacity": 10,
  "status": "ativo"
}
// Esperado: 400 Bad Request
// { "error": "tipo de espaço inválido. Use 'visitantes' ou 'permissionarios'" }
```

---

**Status**: 🟡 Em Progresso

