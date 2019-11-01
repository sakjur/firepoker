package sms

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Phonenumber represents a phonenumber in the E.164 format.
type Phonenumber string

// Valid performs a basic E.164 lookup.
func (n Phonenumber) Valid() error {
	if len(n) <= 3 {
		return fmt.Errorf("the phone number must be longer than 3 characters, got %d characters", len(n))
	} else if len(n) > 15 {
		return fmt.Errorf("the phone number must be at most 15 characters, got %d characters", len(n))
	}

	if n[0] != '+' {
		return fmt.Errorf("expected the first character in the phone number to be '+', found character %c", n[0])
	}

	for _, c := range n[1:] {
		if !unicode.IsNumber(c) {
			return fmt.Errorf("phone numbers must contain only numbers after the initial '+', found character %c", c)
		}
	}

	return nil
}

// Message represents an SMS, possibly with multiple parts.
type Message struct {
	Content string
	Target  Phonenumber
	Inbound bool
}

const gsm7 = "@£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà"

// GSM7 checks if the characters in the Message are adhering to the 7-bit SMS
// character set. When a character outside of this range is used, the SMS is
// forced to use the 16 bit UTF-16 structure instead reducing the number of
// characters which can be used in a SMS from ~153 to 67.
func (m Message) GSM7() bool {
	table := make(map[rune]struct{}, 0)
	for _, c := range gsm7 {
		table[c] = struct{}{}
	}

	for _, c := range m.Content {
		if _, exists := table[c]; !exists {
			return false
		}
	}
	return true
}

// Parts returns the estimated number of parts in which the SMS is sent for
// multipart SMSes.
func (m Message) Parts() int {
	characters := utf8.RuneCountInString(m.Content)

	singleMax := 72
	partLen := 67

	if m.GSM7() {
		singleMax = 160
		partLen = 153
	}

	if characters > singleMax {
		return characters/partLen + 1
	}
	return 1
}

// Valid checks if it is reasonable to send the message.
func (m Message) Valid() error {
	if err := m.Target.Valid(); err != nil {
		return fmt.Errorf("recipient phonenumber is not valid: %w", err)
	}

	if m.Parts() > 10 {
		return fmt.Errorf("number of parts of a message is limited to 10")
	}
	return nil
}
