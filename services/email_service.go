package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	"github.com/Ulpio/reservas-cipt/config"
	"github.com/Ulpio/reservas-cipt/models"
)

// EmailService gerencia o envio de emails
type EmailService struct {
	config config.EmailConfig
}

// NewEmailService cria uma nova instância do serviço de email
func NewEmailService() *EmailService {
	return &EmailService{
		config: config.GetEmailConfig(),
	}
}

// EmailData contém os dados para preencher os templates
type EmailData struct {
	ClientName   string
	SpaceName    string
	Date         string
	StartTime    string
	EndTime      string
	DurationHours int
	Capacity     int
}

// SendReservationConfirmation envia email de confirmação de reserva
func (s *EmailService) SendReservationConfirmation(client models.Client, space models.Space, reservation models.Reservation, startTime, endTime string) error {
	if !s.config.Enabled {
		log.Println("📧 Email desabilitado. Simulando envio de confirmação para:", client.Email)
		return nil
	}

	if client.Email == "" {
		log.Println("⚠️  Cliente não tem email cadastrado")
		return nil
	}

	data := EmailData{
		ClientName:    client.Name,
		SpaceName:     space.Name,
		Date:          reservation.Date.Format("02/01/2006"),
		StartTime:     startTime,
		EndTime:       endTime,
		DurationHours: reservation.DurationHours,
		Capacity:      int(space.Capacity),
	}

	subject := fmt.Sprintf("Reserva Confirmada - %s", space.Name)
	body, err := s.renderTemplate("confirmation", data)
	if err != nil {
		return err
	}

	return s.send(client.Email, subject, body)
}

// SendReservationCancellation envia email de cancelamento de reserva
func (s *EmailService) SendReservationCancellation(client models.Client, space models.Space, reservation models.Reservation, startTime, endTime string) error {
	if !s.config.Enabled {
		log.Println("📧 Email desabilitado. Simulando envio de cancelamento para:", client.Email)
		return nil
	}

	if client.Email == "" {
		log.Println("⚠️  Cliente não tem email cadastrado")
		return nil
	}

	data := EmailData{
		ClientName:    client.Name,
		SpaceName:     space.Name,
		Date:          reservation.Date.Format("02/01/2006"),
		StartTime:     startTime,
		EndTime:       endTime,
		DurationHours: reservation.DurationHours,
	}

	subject := fmt.Sprintf("Reserva Cancelada - %s", space.Name)
	body, err := s.renderTemplate("cancellation", data)
	if err != nil {
		return err
	}

	return s.send(client.Email, subject, body)
}

// send envia o email usando SMTP
func (s *EmailService) send(to, subject, body string) error {
	if !s.config.Enabled {
		return nil
	}

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
		"\r\n"+
		"%s\r\n", s.config.From, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	err := smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("❌ Erro ao enviar email para %s: %v\n", to, err)
		return err
	}

	log.Printf("✅ Email enviado com sucesso para: %s\n", to)
	return nil
}

// renderTemplate renderiza um template de email
func (s *EmailService) renderTemplate(templateName string, data EmailData) (string, error) {
	templates := map[string]string{
		"confirmation": confirmationTemplate,
		"cancellation": cancellationTemplate,
	}

	tmplStr, exists := templates[templateName]
	if !exists {
		return "", fmt.Errorf("template %s não encontrado", templateName)
	}

	tmpl, err := template.New(templateName).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Templates de email em HTML
const confirmationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .info-box { background: white; padding: 20px; margin: 20px 0; border-left: 4px solid #667eea; border-radius: 5px; }
        .info-item { margin: 10px 0; }
        .info-label { font-weight: bold; color: #667eea; }
        .footer { text-align: center; margin-top: 30px; color: #999; font-size: 12px; }
        .button { display: inline-block; padding: 12px 30px; background: #667eea; color: white; text-decoration: none; border-radius: 5px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Reserva Confirmada!</h1>
        </div>
        <div class="content">
            <p>Olá <strong>{{.ClientName}}</strong>,</p>
            <p>Sua reserva foi confirmada com sucesso no <strong>CIPT - Centro de Inovação Público de Tecnologia</strong>!</p>
            
            <div class="info-box">
                <div class="info-item">
                    <span class="info-label">📅 Data:</span> {{.Date}}
                </div>
                <div class="info-item">
                    <span class="info-label">🕐 Horário:</span> {{.StartTime}} às {{.EndTime}}
                </div>
                <div class="info-item">
                    <span class="info-label">⏱️  Duração:</span> {{.DurationHours}} hora(s)
                </div>
                <div class="info-item">
                    <span class="info-label">🏢 Sala:</span> {{.SpaceName}}
                </div>
                <div class="info-item">
                    <span class="info-label">👥 Capacidade:</span> {{.Capacity}} pessoas
                </div>
            </div>
            
            <p><strong>Importante:</strong></p>
            <ul>
                <li>Chegue com 5-10 minutos de antecedência</li>
                <li>Apresente um documento de identificação na recepção</li>
                <li>Em caso de atraso ou cancelamento, entre em contato conosco</li>
            </ul>
            
            <p>Obrigado por utilizar o CIPT!</p>
        </div>
        <div class="footer">
            <p>CIPT - Centro de Inovação Público de Tecnologia</p>
            <p>Este é um email automático, por favor não responda.</p>
        </div>
    </div>
</body>
</html>
`

const cancellationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .info-box { background: white; padding: 20px; margin: 20px 0; border-left: 4px solid #f5576c; border-radius: 5px; }
        .info-item { margin: 10px 0; }
        .info-label { font-weight: bold; color: #f5576c; }
        .footer { text-align: center; margin-top: 30px; color: #999; font-size: 12px; }
        .alert { background: #fff3cd; border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>❌ Reserva Cancelada</h1>
        </div>
        <div class="content">
            <p>Olá <strong>{{.ClientName}}</strong>,</p>
            <p>Sua reserva foi <strong>cancelada</strong>.</p>
            
            <div class="info-box">
                <div class="info-item">
                    <span class="info-label">📅 Data:</span> {{.Date}}
                </div>
                <div class="info-item">
                    <span class="info-label">🕐 Horário:</span> {{.StartTime}} às {{.EndTime}}
                </div>
                <div class="info-item">
                    <span class="info-label">🏢 Sala:</span> {{.SpaceName}}
                </div>
            </div>
            
            <div class="alert">
                <strong>⚠️  Atenção:</strong> Se você não solicitou este cancelamento, entre em contato conosco imediatamente.
            </div>
            
            <p>Se desejar fazer uma nova reserva, entre em contato com a recepção.</p>
            
            <p>Atenciosamente,<br>Equipe CIPT</p>
        </div>
        <div class="footer">
            <p>CIPT - Centro de Inovação Público de Tecnologia</p>
            <p>Este é um email automático, por favor não responda.</p>
        </div>
    </div>
</body>
</html>
`

