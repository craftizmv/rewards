package contracts

type Mailer interface {
	SendEmail(name string, emailArr string, data string) error
}
