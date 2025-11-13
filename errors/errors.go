package errors

import "net/http"

// APIError representa um erro padronizado da API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Field      string `json:"field,omitempty"`      // Campo específico que causou o erro
	Code       string `json:"code,omitempty"`       // Código do erro para o frontend
	Details    string `json:"details,omitempty"`    // Detalhes adicionais
}

// Error implementa a interface error
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError cria um novo erro da API
func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// WithField adiciona informação sobre o campo que causou o erro
func (e *APIError) WithField(field string) *APIError {
	e.Field = field
	return e
}

// WithCode adiciona um código de erro
func (e *APIError) WithCode(code string) *APIError {
	e.Code = code
	return e
}

// WithDetails adiciona detalhes adicionais
func (e *APIError) WithDetails(details string) *APIError {
	e.Details = details
	return e
}

// Erros comuns - Autenticação (401)
var (
	ErrTokenInvalido = &APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    "Token de autenticação inválido ou expirado",
		Code:       "TOKEN_INVALID",
	}

	ErrTokenNaoFornecido = &APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    "Token de autenticação não fornecido",
		Code:       "TOKEN_MISSING",
	}

	ErrCredenciaisInvalidas = &APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    "CPF ou senha incorretos",
		Code:       "INVALID_CREDENTIALS",
	}
)

// Erros comuns - Autorização (403)
var (
	ErrAcessoNegado = &APIError{
		StatusCode: http.StatusForbidden,
		Message:    "Você não tem permissão para acessar este recurso",
		Code:       "ACCESS_DENIED",
	}

	ErrApenasRecepcionista = &APIError{
		StatusCode: http.StatusForbidden,
		Message:    "Esta ação requer permissão de recepcionista ou administrador",
		Code:       "RECEPTIONIST_ONLY",
	}

	ErrApenasAdmin = &APIError{
		StatusCode: http.StatusForbidden,
		Message:    "Esta ação requer permissão de administrador",
		Code:       "ADMIN_ONLY",
	}
)

// Erros comuns - Não Encontrado (404)
var (
	ErrClienteNaoEncontrado = &APIError{
		StatusCode: http.StatusNotFound,
		Message:    "Cliente não encontrado no sistema",
		Code:       "CLIENT_NOT_FOUND",
	}

	ErrUsuarioNaoEncontrado = &APIError{
		StatusCode: http.StatusNotFound,
		Message:    "Usuário não encontrado no sistema",
		Code:       "USER_NOT_FOUND",
	}

	ErrEspacoNaoEncontrado = &APIError{
		StatusCode: http.StatusNotFound,
		Message:    "Espaço não encontrado no sistema",
		Code:       "SPACE_NOT_FOUND",
	}

	ErrReservaNaoEncontrada = &APIError{
		StatusCode: http.StatusNotFound,
		Message:    "Reserva não encontrada no sistema",
		Code:       "RESERVATION_NOT_FOUND",
	}

	ErrStrikeNaoEncontrado = &APIError{
		StatusCode: http.StatusNotFound,
		Message:    "Strike não encontrado no sistema",
		Code:       "STRIKE_NOT_FOUND",
	}
)

// Erros comuns - Validação (400)
var (
	ErrDadosInvalidos = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Os dados fornecidos são inválidos",
		Code:       "INVALID_DATA",
	}

	ErrCampoObrigatorio = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Campo obrigatório não fornecido",
		Code:       "REQUIRED_FIELD",
	}

	ErrFormatoInvalido = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Formato inválido",
		Code:       "INVALID_FORMAT",
	}

	ErrIDInvalido = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "ID fornecido é inválido",
		Code:       "INVALID_ID",
	}
)

// Erros - CPF
var (
	ErrCPFInvalido = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "CPF inválido. Verifique se possui 11 dígitos e os dígitos verificadores estão corretos",
		Code:       "INVALID_CPF",
		Field:      "cpf",
	}

	ErrCPFJaCadastrado = &APIError{
		StatusCode: http.StatusConflict,
		Message:    "CPF já cadastrado no sistema",
		Code:       "CPF_ALREADY_EXISTS",
		Field:      "cpf",
	}
)

// Erros - Email
var (
	ErrEmailInvalido = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Email inválido. Verifique o formato do email",
		Code:       "INVALID_EMAIL",
		Field:      "email",
	}

	ErrEmailJaCadastrado = &APIError{
		StatusCode: http.StatusConflict,
		Message:    "Email já cadastrado no sistema",
		Code:       "EMAIL_ALREADY_EXISTS",
		Field:      "email",
	}
)

// Erros - Telefone
var (
	ErrTelefoneInvalido = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Telefone inválido. Deve conter entre 10 e 11 dígitos",
		Code:       "INVALID_PHONE",
		Field:      "phone",
	}

	ErrTelefoneJaCadastrado = &APIError{
		StatusCode: http.StatusConflict,
		Message:    "Telefone já cadastrado no sistema",
		Code:       "PHONE_ALREADY_EXISTS",
		Field:      "phone",
	}
)

// Erros - Cliente
var (
	ErrIdadeMinima = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Cliente deve ter pelo menos 18 anos de idade",
		Code:       "MINIMUM_AGE",
		Field:      "birth_date",
	}

	ErrNomeInvalido = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Nome deve ter entre 3 e 255 caracteres",
		Code:       "INVALID_NAME",
		Field:      "name",
	}

	ErrDataNascimentoInvalida = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Data de nascimento inválida. Use o formato YYYY-MM-DD",
		Code:       "INVALID_BIRTH_DATE",
		Field:      "birth_date",
	}
)

// Erros - Reservas
var (
	ErrDataInvalida = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Data inválida. Use o formato YYYY-MM-DD",
		Code:       "INVALID_DATE",
		Field:      "date",
	}

	ErrDataPassado = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "A data da reserva não pode ser no passado",
		Code:       "PAST_DATE",
		Field:      "date",
	}

	ErrHorarioInvalido = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Horário inválido. Use o formato HH:MM ou HH:MM:SS",
		Code:       "INVALID_TIME",
		Field:      "start_time",
	}

	ErrDuracaoInvalida = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Duração deve ser entre 1 e 24 horas",
		Code:       "INVALID_DURATION",
		Field:      "duration_hours",
	}

	ErrEspacoOcupado = &APIError{
		StatusCode: http.StatusConflict,
		Message:    "Espaço já está reservado para este horário",
		Code:       "SPACE_OCCUPIED",
	}

	ErrEspacoIndisponivel = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Espaço está indisponível no momento (em manutenção)",
		Code:       "SPACE_UNAVAILABLE",
	}

	ErrRecepcionistaInvalido = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Usuário especificado não é um recepcionista válido",
		Code:       "INVALID_RECEPTIONIST",
		Field:      "receptionist_id",
	}
)

// Erros - Espaços
var (
	ErrTipoEspacoInvalido = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Tipo de espaço inválido. Use 'visitantes' ou 'permissionarios'",
		Code:       "INVALID_SPACE_TYPE",
		Field:      "type",
	}

	ErrStatusEspacoInvalido = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Status do espaço inválido. Use 'ativo' ou 'manutencao'",
		Code:       "INVALID_SPACE_STATUS",
		Field:      "status",
	}

	ErrCapacidadeInvalida = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Capacidade deve ser um número positivo",
		Code:       "INVALID_CAPACITY",
		Field:      "capacity",
	}
)

// Erros - Strikes
var (
	ErrStrikeJaExiste = &APIError{
		StatusCode: http.StatusConflict,
		Message:    "Cliente já possui um strike ativo para esta reserva",
		Code:       "STRIKE_ALREADY_EXISTS",
	}

	ErrMotivoObrigatorio = &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Motivo do strike é obrigatório",
		Code:       "REQUIRED_REASON",
		Field:      "reason",
	}
)

// Erros - Internos (500)
var (
	ErrInterno = &APIError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Erro interno do servidor. Tente novamente mais tarde",
		Code:       "INTERNAL_ERROR",
	}

	ErrBancoDados = &APIError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Erro ao acessar o banco de dados",
		Code:       "DATABASE_ERROR",
	}
)

