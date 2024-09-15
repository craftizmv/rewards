package service

type SomeConcreteMailer struct {
}

func NewSomeConcreteMailer() *SomeConcreteMailer {
	return &SomeConcreteMailer{}
}

func (s *SomeConcreteMailer) SendEmail(name string, emailArr string, data string) error {
	// send mail.
	return nil
}
