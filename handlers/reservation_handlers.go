package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Ulpio/reservas-cipt/dto"
	apierrors "github.com/Ulpio/reservas-cipt/errors"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

// CreateReservationHandler registra uma nova reserva.
// @Summary Cria reserva
// @Description Cria uma nova reserva para um espaço com validações completas.
// @Tags reservas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dto.CreateReservationInputDTO true "Dados da reserva"
// @Success 201 {object} dto.ReservationOutputDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /reservas [post]
func CreateReservationHandler(c *gin.Context) {
	var input dto.CreateReservationInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		apierrors.SendError(c, apierrors.ErrDadosInvalidos.WithDetails("Verifique os campos obrigatórios: client_id, receptionist_id, space_id, date, start_time, duration_hours"))
		return
	}

	reservation, err := services.CreateReservationFromInput(input)
	if err != nil {
		handleReservationServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

// GetReservationByIDHandler retorna uma reserva pelo ID.
// @Summary Busca reserva por ID
// @Description Retorna os dados de uma reserva específica.
// @Tags reservas
// @Produce json
// @Param id path int true "ID da reserva"
// @Success 200 {object} dto.ReservationOutputDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /reservas/{id} [get]
func GetReservationByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		apierrors.SendError(c, apierrors.ErrIDInvalido.WithDetails("O ID da reserva deve ser um número válido"))
		return
	}

	reservation, err := services.GetReservationByID(uint(id))
	if err != nil {
		apierrors.SendError(c, apierrors.ErrReservaNaoEncontrada)
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// GetAllReservationsHandler lista todas as reservas.
// @Summary Lista reservas
// @Description Lista todas as reservas cadastradas.
// @Tags reservas
// @Produce json
// @Success 200 {array} dto.ReservationOutputDTO
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservas [get]
func GetAllReservationsHandler(c *gin.Context) {
	reservations, err := services.GetAllReservations()
	if err != nil {
		apierrors.SendError(c, apierrors.ErrBancoDados.WithDetails("Erro ao listar reservas"))
		return
	}
	c.JSON(http.StatusOK, reservations)
}

// handleReservationServiceError trata erros específicos dos serviços de reserva
func handleReservationServiceError(c *gin.Context, err error) {
	errMsg := err.Error()

	// Erros de validação de data
	if strings.Contains(errMsg, "data inválida") || strings.Contains(errMsg, "formato de data") {
		apierrors.SendError(c, apierrors.ErrDataInvalida)
		return
	}
	if strings.Contains(errMsg, "não pode ser no passado") || strings.Contains(errMsg, "data passada") {
		apierrors.SendError(c, apierrors.ErrDataPassado)
		return
	}

	// Erros de validação de horário
	if strings.Contains(errMsg, "horário inválido") || strings.Contains(errMsg, "formato de horário") {
		apierrors.SendError(c, apierrors.ErrHorarioInvalido)
		return
	}

	// Erros de validação de duração
	if strings.Contains(errMsg, "duração") {
		apierrors.SendError(c, apierrors.ErrDuracaoInvalida)
		return
	}

	// Erros de entidades não encontradas
	if strings.Contains(errMsg, "cliente não encontrado") {
		apierrors.SendError(c, apierrors.ErrClienteNaoEncontrado.WithDetails("Verifique se o client_id está correto"))
		return
	}
	if strings.Contains(errMsg, "espaço não encontrado") {
		apierrors.SendError(c, apierrors.ErrEspacoNaoEncontrado.WithDetails("Verifique se o space_id está correto"))
		return
	}
	if strings.Contains(errMsg, "usuário não encontrado") || strings.Contains(errMsg, "recepcionista não encontrado") {
		apierrors.SendError(c, apierrors.ErrUsuarioNaoEncontrado.WithDetails("Verifique se o receptionist_id está correto"))
		return
	}

	// Erros de validação de recepcionista
	if strings.Contains(errMsg, "não é um recepcionista") || strings.Contains(errMsg, "role de recepcionista") {
		apierrors.SendError(c, apierrors.ErrRecepcionistaInvalido)
		return
	}

	// Erros de disponibilidade
	if strings.Contains(errMsg, "já reservado") || strings.Contains(errMsg, "conflito") {
		apierrors.SendError(c, apierrors.ErrEspacoOcupado.WithDetails("Escolha outro horário ou espaço"))
		return
	}
	if strings.Contains(errMsg, "manutenção") || strings.Contains(errMsg, "indisponível") {
		apierrors.SendError(c, apierrors.ErrEspacoIndisponivel.WithDetails("Escolha outro espaço"))
		return
	}

	// Erro genérico
	apierrors.SendError(c, apierrors.ErrInterno.WithDetails(errMsg))
}
