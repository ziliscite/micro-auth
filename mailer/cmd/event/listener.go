package main

import (
	"encoding/json"
	"github.com/ziliscite/micro-auth/mailer/internal/domain"
	"log/slog"
)

func (s *service) listen() error {
	ch, err := s.amc.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	mails, err := ch.Consume(
		s.mq.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	congratsMails, err := ch.Consume(
		s.mcq.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// consume til application exits
	forever := make(chan bool)
	go func() {
		for m := range mails {
			var mail domain.Mail

			err = json.Unmarshal(m.Body, &mail)
			if err != nil {
				slog.Error(err.Error())
			}

			go s.sendActivationEmail(mail)
		}
	}()

	go func() {
		for m := range congratsMails {
			var mail domain.Mail

			err = json.Unmarshal(m.Body, &mail)
			if err != nil {
				slog.Error(err.Error())
			}

			go s.sendCongratulationsEmail(mail)
		}
	}()

	slog.Info("Listening for emails...")
	<-forever
	return nil
}
