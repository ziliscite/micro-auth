package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ziliscite/micro-auth/mailer/internal"
	"github.com/ziliscite/micro-auth/mailer/internal/domain"
	"log/slog"
)

type service struct {
	amc   *amqp.Connection
	queue amqp.Queue

	mr *internal.Mailer
}

func newService(amc *amqp.Connection, mr *internal.Mailer) (*service, error) {
	ch, err := amc.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"send_mail",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &service{
		amc:   amc,
		queue: q,
		mr:    mr,
	}, nil
}

func (s *service) sendActivationEmail(mail domain.Mail) {
	err := s.mr.Send(mail.Email, "user_welcome.tmpl", map[string]interface{}{
		"username":        mail.Username,
		"userID":          mail.ID,
		"activationToken": mail.Token,
		"expiresAt":       mail.ExpiresAt.Format("Monday, January 2, 2006 at 3:04 PM"),
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to send an email for %s", mail.Email), "error", err.Error())
	}
}
