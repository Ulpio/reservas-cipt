## 📘 Documentação Geral da API CIPT

### Visão Geral
- **Base URL:** `http://{host}:8080/api/v1`
- **Formato:** JSON (`Content-Type: application/json`)
- **Autenticação:** JWT Bearer (`Authorization: Bearer {token}`)
- **Status HTTP Padrão:** `200` (OK), `201` (Criado), `204` (Sem conteúdo), `400` (Requisição inválida), `401` (Não autenticado), `403` (Sem permissão), `404` (Não encontrado), `409` (Conflito), `500` (Erro interno)

> As respostas de erro seguem o formato `{"error": "Mensagem descritiva", ...}`. Alguns endpoints adicionam campos extras como `details` ou `message`.

---

### Autenticação
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/auth/login` | Retorna token JWT (obrigatório para os demais endpoints) |

**Request Body**
```json
{
  "cpf": "00000000000",
  "password": "senha"
}
```

**Resposta 200**
```json
{
  "token": "jwt.token.aqui"
}
```

---

### Clientes (`/clientes`)
> **Permissão:** Token válido. Sem restrição adicional.

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/clientes/buscar-criar` | Busca por CPF e cria cliente se não existir |
| `GET`  | `/clientes/{cpf}` | Busca cliente por CPF |
| `GET`  | `/clientes` | Lista todos os clientes |
| `PATCH`| `/clientes/{id}` | Atualiza dados de um cliente existente |

**DTOs Principais**
- `ClienteInputDTO`: `name`, `cpf`, `email`, `phone`
- `ClienteOutputDTO`: `id`, `name`, `cpf`, `email`, `phone`, `strikes`

**Exemplo de criação ou busca**
```json
POST /clientes/buscar-criar
{
  "name": "João Silva",
  "cpf": "12345678900",
  "email": "joao@email.com",
  "phone": "11987654321"
}
```

---

### Reservas (`/reservas`)
> **Permissão:** Token válido + Papel `recepcionista` ou `admin`.

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/reservas` | Cria nova reserva |
| `GET`  | `/reservas` | Lista todas as reservas |
| `GET`  | `/reservas/{id}` | Detalha reserva por ID |

**DTOs Principais**
- `CreateReservationDTO`: `client_id`, `receptionist_id`, `space_id`, `date` (`ISO`), `start_time` (`ISO`), `duration_hours`
- `ReservationOutputDTO`: `id`, `client_name`, `receptionist_name`, `space_name`, `date`, `start_time`, `duration_hours`, `end_time`

---

### Espaços (`/espacos`)
> **Permissão:** Token válido para leitura/atualizações básicas. Somente `admin` pode criar, atualizar totalmente ou remover.

| Método | Endpoint | Descrição | Papel |
|--------|----------|-----------|-------|
| `GET`  | `/espacos` | Lista espaços | qualquer |
| `GET`  | `/espacos/{id}` | Detalha espaço | qualquer |
| `POST` | `/espacos` | Cria espaço | admin |
| `PUT`  | `/espacos/{id}` | Atualiza espaço | admin |
| `DELETE` | `/espacos/{id}` | Remove espaço | admin |
| `PATCH` | `/espacos/{id}/status` | Atualiza status (`ativo`, `inativo`, `manutencao`, etc.) | qualquer token |
| `PATCH` | `/espacos/{id}/aviso` | Atualiza aviso exibido | qualquer token |

**DTOs Principais**
- `CreateSpaceDTO` / `UpdateSpaceDTO`: `name`, `type` (`visitantes|permissionarios|locacao`), `status`, `notice`, `capacity`
- `SpaceOutputDTO`: `id`, `name`, `type`, `status`, `notice`, `capacity`
- `UpdateStatusDTO`: `status`
- `UpdateNoticeDTO`: `notice`

---

### Usuários (`/users`)
> **Permissão:** Token válido. Algumas operações exigem `admin`.

| Método | Endpoint | Descrição | Papel |
|--------|----------|-----------|-------|
| `GET`  | `/users` | Lista usuários | admin |
| `POST` | `/users` | Cria usuário (admin/recepcionista) | admin |
| `GET`  | `/users/{id}` | Detalha usuário | admin |
| `DELETE` | `/users/{id}` | Remove usuário | admin |
| `GET`  | `/users/me` | Retorna usuário autenticado | qualquer |

**DTOs Principais**
- `UserInputDTO`: `name`, `cpf`, `role`
- `UserOutputDTO`: `user_id`, `name`, `cpf`, `role`

---

### Strikes (`/strikes`)
> **Permissão:** Token válido. Criação/listagem exigem `admin` ou `recepcionista`. Revogação apenas `admin`.

| Método | Endpoint | Descrição | Papel |
|--------|----------|-----------|-------|
| `POST` | `/strikes` | Registra strike | admin/recepcionista |
| `GET`  | `/strikes/client/{id}` | Lista strikes de um cliente | admin/recepcionista |
| `DELETE` | `/strikes/{id}` | Revoga strike | admin |

**DTOs Principais**
- `StrikeInputDTO`: `client_id`, `reason`, `photo`
- `StrikeOutputDTO`: `id`, `client_id`, `reason`, `photo`, `created_at`, `revoked`, `revoked_at`

---

### Resumo de Headers Importantes
```text
Authorization: Bearer {token}
Content-Type: application/json
```

---

### Boas Práticas e Observações
- Sempre normalize CPF/telefone no frontend antes de enviar.
- Datas e horários de reservas devem ser enviados em formato ISO 8601 (`YYYY-MM-DD` e `YYYY-MM-DDTHH:mm:ssZ`).
- Verifique a role do usuário no token antes de habilitar ações sensíveis na interface.
- Em caso de erro, utilize a mensagem retornada pelo servidor e trate detalhes adicionais quando disponíveis.

---

### Referências Úteis
- **Swagger:** `http://{host}:8080/swagger/index.html`
- **README:** Informações sobre setup do projeto e execução local.

