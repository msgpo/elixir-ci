package internal

import (
	"gopkg.in/gomail.v2"
	"log"
)

type Specification struct {
	AppName             string   `envconfig:"APP_NAME"`
	ConclusionsToIgnore []string `envconfig:"IGNORED_CONCLUSIONS" default:"success,neutral"`
	EventPath           string   `envconfig:"GITHUB_EVENT_PATH" required:"true"`
	GitHubToken         string   `envconfig:"GITHUB_TOKEN"`
	MailHost            string   `envconfig:"MAIL_HOST" required:"true"`
	MailPort            int      `envconfig:"MAIL_PORT" default:"587"`
	MailFrom            string   `envconfig:"MAIL_FROM" required:"true"`
	MailUsername        string   `envconfig:"MAIL_USERNAME" required:"true"`
	MailPassword        string   `envconfig:"MAIL_PASSWORD" required:"true"`
}

func SendNotification(spec Specification) {
	log.Printf("Parsing %s...", spec.EventPath)
	event, commit, err := Parse(spec.EventPath)
	if err != nil {
		log.Fatalf("Failed to parse event! %s", err)
	}

	if *event.CheckSuite.App.Name != spec.AppName {
		log.Printf("No need to send email for %s app!", *event.CheckSuite.App.Name)
		NeutralExit()
	}

	if contains(spec.ConclusionsToIgnore, *event.CheckSuite.Conclusion) {
		log.Printf("No need to send email for check suite with %s conclusion!", *event.CheckSuite.Conclusion)
		NeutralExit()
	}

	log.Printf("Dialing %v:%v...\n", spec.MailHost, spec.MailPort)
	var dialer *gomail.Dialer
	if spec.MailUsername == "" && spec.MailPassword == "" {
		dialer = &gomail.Dialer{Host: spec.MailHost, Port: spec.MailPort}
	} else {
		dialer = gomail.NewDialer(spec.MailHost, spec.MailPort, spec.MailUsername, spec.MailPassword)
	}
	message, err := generateEmail(spec, event, commit, *commit.Author.Email, *commit.Author.Name, "the author of this commit")
	if err != nil {
		log.Fatalf("Failed to generate email! %s", err)
	}
	err = dialer.DialAndSend(message)
	if err != nil {
		log.Fatalf("Failed to send email! %s", err)
	}

	if *commit.Author.Email != *commit.Committer.Email && *commit.Committer.Email != "noreply@github.com" {
		message, err := generateEmail(spec, event, commit, *commit.Committer.Email, *commit.Committer.Name, "the committer of this commit")
		if err != nil {
			log.Fatalf("Failed to generate email! %s", err)
		}
		err = dialer.DialAndSend(message)
		if err != nil {
			log.Fatalf("Failed to send email! %s", err)
		}
	}
}
