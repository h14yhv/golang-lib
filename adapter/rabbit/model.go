package rabbit

type Message struct {
	Body        []byte
	ContentType string
	Mode        uint8
	Priority    uint8
}
