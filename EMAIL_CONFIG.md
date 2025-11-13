# 📧 Configuração de Email Automático

## 🔧 Variáveis de Ambiente

Adicione estas variáveis no seu arquivo `.env`:

```env
# Configuração de Email (SMTP)
EMAIL_ENABLED=false                    # true para habilitar, false para desabilitar
SMTP_HOST=smtp.gmail.com              # Host do servidor SMTP
SMTP_PORT=587                          # Porta SMTP (587 para TLS, 465 para SSL)
SMTP_USER=seu-email@gmail.com         # Usuário SMTP
SMTP_PASSWORD=sua-senha-app           # Senha do app (não a senha normal)
EMAIL_FROM=noreply@cipt.com           # Email remetente
```

---

## 📮 Configuração por Provedor

### Gmail

1. Acesse: https://myaccount.google.com/apppasswords
2. Crie uma "Senha de app"
3. Use essa senha no `SMTP_PASSWORD`

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=seu-email@gmail.com
SMTP_PASSWORD=xxxx-xxxx-xxxx-xxxx
```

### Outlook/Hotmail

```env
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USER=seu-email@outlook.com
SMTP_PASSWORD=sua-senha
```

### SendGrid

```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=SG.xxxxxxxxxxxxxx
```

---

## 🚀 Como Funciona

### Modo Desabilitado (EMAIL_ENABLED=false)
- Emails NÃO são enviados
- Sistema apenas loga mensagens no console:
  ```
  📧 Email desabilitado. Simulando envio de confirmação para: cliente@example.com
  ```

### Modo Habilitado (EMAIL_ENABLED=true)
- Emails são enviados via SMTP
- Sistema loga sucesso ou erro:
  ```
  ✅ Email enviado com sucesso para: cliente@example.com
  ❌ Erro ao enviar email para cliente@example.com: ...
  ```

---

## 📬 Emails Enviados Automaticamente

### 1. Confirmação de Reserva
**Quando**: Ao criar uma nova reserva  
**Para**: Email do cliente (se cadastrado)  
**Assunto**: `Reserva Confirmada - {nome_da_sala}`

**Conteúdo**:
- Nome do cliente
- Data e horário da reserva
- Nome da sala
- Capacidade
- Instruções importantes

### 2. Cancelamento de Reserva
**Quando**: Ao cancelar uma reserva  
**Para**: Email do cliente (se cadastrado)  
**Assunto**: `Reserva Cancelada - {nome_da_sala}`

**Conteúdo**:
- Nome do cliente
- Data e horário da reserva
- Nome da sala
- Alerta de segurança

---

## 🧪 Testando

### Teste em Desenvolvimento:

1. **Sem Email Real** (Recomendado):
```env
EMAIL_ENABLED=false
```
- Sistema simula envio
- Nenhum email real é enviado
- Perfeito para desenvolvimento

2. **Com Email Real**:
```env
EMAIL_ENABLED=true
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=seu-email@gmail.com
SMTP_PASSWORD=sua-senha-app
```

3. **Criar uma Reserva**:
```bash
POST /api/reservas
{
  "client_id": 1,
  "receptionist_id": 2,
  "space_id": 3,
  "date": "2025-11-15",
  "start_time": "14:00",
  "duration_hours": 2
}
```

4. **Verificar Logs**:
```
📧 Email desabilitado. Simulando envio...
# ou
✅ Email enviado com sucesso para: cliente@example.com
```

---

## ⚠️ Troubleshooting

### Erro: "535 Authentication failed"
**Problema**: Senha incorreta ou 2FA ativado  
**Solução**: Use senha de app (não sua senha normal)

### Erro: "Connection timeout"
**Problema**: Porta ou host incorretos  
**Solução**: Verifique SMTP_HOST e SMTP_PORT

### Email não chega
**Problema**: Email pode estar no spam  
**Solução**: 
1. Verifique pasta de spam
2. Adicione EMAIL_FROM nos contatos
3. Configure SPF/DKIM (produção)

### Cliente sem email
**Log**: `⚠️  Cliente não tem email cadastrado`  
**Solução**: Normal. Email opcional para clientes.

---

## 🔒 Segurança

### ⚠️ NUNCA commite credenciais!
```bash
# Adicione ao .gitignore
.env
.env.local
.env.production
```

### ✅ Use variáveis de ambiente
```bash
# Desenvolvimento
export SMTP_PASSWORD="senha-segura"

# Produção (Docker/Kubernetes)
kubectl create secret generic smtp-credentials \
  --from-literal=password='senha-segura'
```

---

## 📊 Logs

O sistema registra todas as tentativas de envio:

```
📧 Email desabilitado. Simulando envio de confirmação para: joao@example.com
✅ Email enviado com sucesso para: maria@example.com
⚠️  Cliente não tem email cadastrado
❌ Erro ao enviar email para erro@example.com: authentication failed
```

---

## 🎨 Templates

Os templates de email são HTML responsivos com:
- ✅ Design moderno e profissional
- ✅ Gradientes coloridos
- ✅ Informações bem organizadas
- ✅ Compatibilidade com clients de email
- ✅ Encodings UTF-8 (acentos funcionam)

---

## 🔄 Próximas Features (Opcional)

- [ ] Email de lembrete 24h antes
- [ ] Email ao receber strike
- [ ] Email de boas-vindas ao cadastrar cliente
- [ ] Filas de email (Redis/RabbitMQ)
- [ ] Retry automático em caso de falha
- [ ] Dashboard de emails enviados

---

**Status**: ✅ Implementado e Pronto para Uso

