package message

// Standard essage that will be sent to client and sent from client
type JSON struct {
	Type string `json:"type"`
	Body any    `json:"body"`
}
