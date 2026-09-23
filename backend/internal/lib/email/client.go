package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"inktype-backend/internal/config"

	"github.com/pkg/errors"
	"github.com/resend/resend-go/v2"
	"github.com/rs/zerolog"
)

//go:embed templates/emails/*.html
var emailTemplates embed.FS

type Client struct {
	client    *resend.Client
	logger    *zerolog.Logger
	templates *template.Template
}

func NewClient(cfg *config.Config, logger *zerolog.Logger) *Client {
	// Pre-parse all templates
	tmpl, err := template.ParseFS(emailTemplates, "templates/emails/*.html")
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse email templates")
	}

	return &Client{
		client:    resend.NewClient(cfg.Integration.ResendAPIKey),
		logger:    logger,
		templates: tmpl,
	}
}

func (c *Client) SendEmail(to, subject string, templateName Template, data map[string]string) error {
	if c.templates == nil {
		return errors.New("email templates not initialized")
	}

	var body bytes.Buffer
	tmplName := fmt.Sprintf("%s.html", templateName)
	if err := c.templates.ExecuteTemplate(&body, tmplName, data); err != nil {
		return errors.Wrapf(err, "failed to execute email template %s", templateName)
	}

	params := &resend.SendEmailRequest{
		From:    "InkType <hello@inktype.app>",
		To:      []string{to},
		Subject: subject,
		Html:    body.String(),
	}

	_, err := c.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
