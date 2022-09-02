package mail

type Service interface {
	Send(email *Email) error
}
