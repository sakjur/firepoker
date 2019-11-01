package sms

// Sender delivers a SMS to the telephony network.
type Sender interface {
	Send(message Message) error
}
