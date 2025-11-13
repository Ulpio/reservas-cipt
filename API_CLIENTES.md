# 📋 API - Busca de Clientes

## Endpoints Implementados

Todos os endpoints requerem autenticação via JWT Token.

**Header obrigatório:**
```
Authorization: Bearer {token}
```

---

## 1. Buscar Cliente por CPF

### **GET** `/api/v1/clientes/cpf/:cpf`

Busca um cliente pelo número do CPF (aceita formato com ou sem máscara).

#### Parâmetros:
- `cpf` (URL param): CPF do cliente
  - Aceita: `12345678900` ou `123.456.789-00`

#### Resposta de Sucesso (200):
```json
{
  "id": 1,
  "name": "João Silva",
  "cpf": "12345678900",
  "phone": "11987654321",
  "email": "joao@email.com",
  "strikes": 0
}
```

#### Resposta de Erro (404):
```json
{
  "error": "Cliente não encontrado"
}
```

#### Exemplo:
```bash
curl -X GET "http://localhost:8080/api/v1/clientes/cpf/123.456.789-00" \
  -H "Authorization: Bearer {token}"
```

---

## 2. Buscar Cliente por Telefone

### **GET** `/api/v1/clientes/phone/:phone`

Busca um cliente pelo número de telefone (aceita formato com ou sem máscara).

#### Parâmetros:
- `phone` (URL param): Telefone do cliente
  - Aceita: `11987654321` ou `(11) 98765-4321`

#### Resposta de Sucesso (200):
```json
{
  "id": 1,
  "name": "João Silva",
  "cpf": "12345678900",
  "phone": "11987654321",
  "email": "joao@email.com",
  "strikes": 0
}
```

#### Resposta de Erro (404):
```json
{
  "error": "Cliente não encontrado"
}
```

#### Exemplo:
```bash
curl -X GET "http://localhost:8080/api/v1/clientes/phone/(11) 98765-4321" \
  -H "Authorization: Bearer {token}"
```

---

## 3. Criar Novo Cliente

### **POST** `/api/v1/clientes`

Cadastra um novo cliente no sistema.

#### Body (JSON):
```json
{
  "name": "João Silva",
  "cpf": "123.456.789-00",
  "phone": "(11) 98765-4321",
  "email": "joao@email.com"
}
```

#### Campos:
| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `name` | string | Sim | Nome completo do cliente |
| `cpf` | string | Sim | CPF do cliente (único) |
| `phone` | string | Sim | Telefone do cliente (único) |
| `email` | string | Não | Email do cliente |

#### Resposta de Sucesso (201):
```json
{
  "id": 1,
  "name": "João Silva",
  "cpf": "12345678900",
  "phone": "11987654321",
  "email": "joao@email.com",
  "strikes": 0
}
```

#### Respostas de Erro:

**400 - CPF duplicado:**
```json
{
  "error": "CPF já cadastrado"
}
```

**400 - Telefone duplicado:**
```json
{
  "error": "Telefone já cadastrado"
}
```

**400 - Dados inválidos:**
```json
{
  "error": "Dados inválidos"
}
```

#### Exemplo:
```bash
curl -X POST "http://localhost:8080/api/v1/clientes" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "cpf": "123.456.789-00",
    "phone": "(11) 98765-4321",
    "email": "joao@email.com"
  }'
```

---

## 4. Buscar Cliente por ID

### **GET** `/api/v1/clientes/:id`

Busca um cliente pelo ID.

#### Parâmetros:
- `id` (URL param): ID do cliente

#### Resposta de Sucesso (200):
```json
{
  "id": 1,
  "name": "João Silva",
  "cpf": "12345678900",
  "phone": "11987654321",
  "email": "joao@email.com",
  "strikes": 0
}
```

#### Resposta de Erro (404):
```json
{
  "error": "Cliente não encontrado"
}
```

#### Exemplo:
```bash
curl -X GET "http://localhost:8080/api/v1/clientes/1" \
  -H "Authorization: Bearer {token}"
```

---

## 5. Listar Todos os Clientes

### **GET** `/api/v1/clientes`

Retorna todos os clientes cadastrados.

#### Resposta de Sucesso (200):
```json
[
  {
    "id": 1,
    "name": "João Silva",
    "cpf": "12345678900",
    "phone": "11987654321",
    "email": "joao@email.com",
    "strikes": 0
  },
  {
    "id": 2,
    "name": "Maria Santos",
    "cpf": "98765432100",
    "phone": "21987654321",
    "email": "maria@email.com",
    "strikes": 1
  }
]
```

---

## 6. Buscar ou Criar Cliente (Legado)

### **POST** `/api/v1/clientes/buscar-criar`

Busca um cliente pelo CPF e cria caso não exista.

**Nota:** Este endpoint está mantido para compatibilidade. Recomenda-se usar os endpoints específicos de busca e criação.

#### Body (JSON):
```json
{
  "cpf": "123.456.789-00",
  "name": "João Silva",
  "phone": "(11) 98765-4321",
  "email": "joao@email.com"
}
```

---

## Validações e Normalizações

### CPF
- ✅ Aceita com ou sem formatação: `123.456.789-00` ou `12345678900`
- ✅ Normalizado automaticamente (remove pontos e traços)
- ✅ Deve ser único no sistema

### Telefone
- ✅ Aceita com ou sem formatação: `(11) 98765-4321` ou `11987654321`
- ✅ Normalizado automaticamente (remove caracteres especiais)
- ✅ Deve ser único no sistema

### Email
- ⚠️ Opcional
- ✅ Validação de formato

---

## Status Codes

| Code | Descrição |
|------|-----------|
| `200` | Sucesso - retorna os dados |
| `201` | Criado - cliente cadastrado com sucesso |
| `400` | Dados inválidos ou duplicados |
| `401` | Não autenticado - token inválido ou ausente |
| `404` | Cliente não encontrado |
| `500` | Erro interno do servidor |

---

## Exemplos de Uso no Frontend

### JavaScript/Axios

```javascript
import api from './services/api';

// Buscar por CPF
const buscarPorCPF = async (cpf) => {
  try {
    const { data } = await api.get(`/clientes/cpf/${cpf}`);
    console.log('Cliente encontrado:', data);
  } catch (error) {
    if (error.response?.status === 404) {
      console.log('Cliente não encontrado');
    }
  }
};

// Buscar por telefone
const buscarPorTelefone = async (phone) => {
  const { data } = await api.get(`/clientes/phone/${phone}`);
  return data;
};

// Criar novo cliente
const criarCliente = async (clienteData) => {
  try {
    const { data } = await api.post('/clientes', {
      name: 'João Silva',
      cpf: '123.456.789-00',
      phone: '(11) 98765-4321',
      email: 'joao@email.com'
    });
    console.log('Cliente criado:', data);
  } catch (error) {
    if (error.response?.data?.error === 'CPF já cadastrado') {
      alert('Este CPF já está cadastrado!');
    }
  }
};

// Buscar por ID
const buscarPorID = async (id) => {
  const { data } = await api.get(`/clientes/${id}`);
  return data;
};
```

---

## Testes

Execute os testes com:

```bash
go test ./tests/handlers/client_handler_test.go -v
```

### Casos de Teste Implementados:

1. ✅ Criar cliente com sucesso
2. ✅ Validar CPF duplicado
3. ✅ Validar telefone duplicado
4. ✅ Buscar por CPF formatado e não formatado
5. ✅ Buscar por telefone formatado e não formatado
6. ✅ Buscar por ID
7. ✅ Retornar 404 para cliente não encontrado
8. ✅ Retornar 401 sem autenticação

---

## Arquivos Criados/Modificados

### Novos Arquivos:
1. **`utils/normalize.go`** - Funções de normalização de CPF e telefone
2. **`tests/handlers/client_handler_test.go`** - Testes dos endpoints

### Arquivos Modificados:
1. **`services/client_services.go`** - Adicionadas funções:
   - `GetClientByPhone()`
   - `GetClientByID()`
   - `CreateClient()`
   - Normalização em `BuscarOuCriarCliente()`

2. **`handlers/client_handlers.go`** - Adicionados handlers:
   - `BuscarClientePorTelefone()`
   - `BuscarClientePorID()`
   - `CreateClientHandler()`

3. **`routes/client_routes.go`** - Adicionadas rotas:
   - `GET /clientes/cpf/:cpf`
   - `GET /clientes/phone/:phone`
   - `GET /clientes/:id`
   - `POST /clientes`

---

## Segurança

✅ **Autenticação JWT obrigatória** em todos os endpoints  
✅ **Validação de unicidade** para CPF e telefone  
✅ **Normalização automática** de dados  
✅ **Queries parametrizadas** (SQL Injection protection)  
✅ **Validação de entrada** em todos os endpoints

---

## Swagger/OpenAPI

A documentação Swagger foi atualizada automaticamente. Acesse:

```
http://localhost:8080/swagger/index.html
```

---

## Performance

### Otimizações Implementadas:
1. ✅ Índices no banco de dados (CPF e telefone)
2. ✅ Queries otimizadas com GORM
3. ✅ Normalização antes da busca
4. ✅ Validações eficientes

### Recomendações Futuras:
- Cache de clientes frequentemente buscados
- Paginação na listagem de todos os clientes
- Busca fuzzy/parcial por nome

---

**Desenvolvido por:** Ulpio Paulo de Miranda Netto  
**Data:** 2025  
**Versão API:** 1.0

