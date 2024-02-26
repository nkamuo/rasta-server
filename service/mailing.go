package service

import (
	"sync"
	// "github.com/sendgrid/sendgrid-go"
	// "github.com/sendgrid/sendgrid-go/helpers/mail"
)

var mailingService MailingService
var mailingRepoMutext *sync.Mutex = &sync.Mutex{}

func GetMailingService() MailingService {
	mailingRepoMutext.Lock()
	if mailingService == nil {
		mailingService = &mailingServiceImpl{}
	}
	mailingRepoMutext.Unlock()
	return mailingService
}

type MailingService interface {
	// NewEmail(name string, address string,) *mail.SGMailV3

}

type mailingServiceImpl struct {
}
