package sms

import (
	"strings"
	"testing"
)

func TestMessage_Parts_emptyMessage_onePart(t *testing.T) {
	empty := Message{
		Content: "",
	}

	if empty.Parts() != 1 {
		t.Errorf("expected an empty message to be 1 part, got %d", empty.Parts())
	}
}
func TestMessage_simpleMessage_gsm7_onePart(t *testing.T) {
	m := Message{
		Content: "Good morning, upper east side. XOXO Gossip Girl",
	}

	if !m.GSM7() {
		t.Errorf("expected the message to be encodable with GSM7, but it isn't")
	}

	if m.Parts() != 1 {
		t.Errorf("expected the message to have 1 part, got %d", m.Parts())
	}
}

func TestMessage_complicatedMessage_singlePart(t *testing.T) {
	m := Message{
		// not a single GSM7 character
		Content: "こんにちは世界",
	}

	if m.GSM7() {
		t.Errorf("the message shouldn't be encodable with GSM7, the GSM7 function claims it is")
	}

	if m.Parts() != 1 {
		t.Errorf("expected the message to have 1 part, got %d", m.Parts())
	}
}

func TestMessage_Parts_nonGSM7(t *testing.T) {
	m := Message{
		Content: "🚀 An emoji forces the message to be re-encoded to a unicode charset, meaning that this message has to be split",
	}

	if m.GSM7() {
		t.Errorf("the message shouldn't be encodable with GSM7, the GSM7 function claims it is")
	}

	if m.Parts() != 2 {
		t.Errorf("expected the message to have 2 parts, got %d", m.Parts())
	}
}

func TestMessage_Parts_GSM7_single_v_multi_count(t *testing.T) {
	msg := strings.Repeat("A", 160)

	x1 := Message{
		Content: msg,
	}

	if x1.Parts() != 1 {
		t.Errorf("expected the 160 character message to have 1 parts, got %d", x1.Parts())
	}

	x2 := Message{
		Content: strings.Repeat(msg, 2),
	}

	if x2.Parts() != 3 {
		t.Errorf("expected the 320 character message to have 3 parts, got %d", x1.Parts())
	}
}

func TestMessage_Parts_Unicode_single_v_multi_count(t *testing.T) {
	msg := strings.Repeat("🚀", 72)

	x1 := Message{
		Content: msg,
	}

	if x1.Parts() != 1 {
		t.Errorf("expected the 72 character message to have 1 parts, got %d", x1.Parts())
	}

	x2 := Message{
		Content: strings.Repeat(msg, 2),
	}

	if x2.Parts() != 3 {
		t.Errorf("expected the 144 character message to have 3 parts, got %d", x1.Parts())
	}
}

func TestPhonenumber_Valid_SimpleFail(t *testing.T) {
	invalid := []string{
		"",
		"+",
		"+1", // too short
		"text",
		"0700000000",
		"070-000 00 00",
		"+4670000weird",
		"+123456789123456789", // too long
	}

	for _, s := range invalid {
		err := Phonenumber(s).Valid()
		if err == nil {
			t.Errorf("expected %s to not be valid, but it was", s)
		}
	}
}

func TestPhonenumber_Valid_SimpleSuccess(t *testing.T) {
	valid := []string{
		"+1123555000",
		"+4670000000",
		"+123",
	}

	for _, s := range valid {
		err := Phonenumber(s).Valid()
		if err != nil {
			t.Errorf("expected %s to be valid, got error %v", s, err)
		}
	}
}

const validMsg = "New phone, who?"
const validPhonenumber = "+123"

func TestMessage_Valid_WrongPhonenumber_Fail(t *testing.T) {
	m := Message{
		Target:  "",
		Content: validMsg,
	}

	err := m.Valid()
	if err == nil {
		t.Errorf("expected validation to fail when the target isn't a valid phonenumber")
	}
}

func TestMessage_Valid_ContentTooLong_Fail(t *testing.T) {
	m := Message{
		Target:  validPhonenumber,
		Content: strings.Repeat("A", 1600), // more than 10 parts
	}

	err := m.Valid()
	if err == nil {
		t.Errorf("expected validation to fail when content is too long")
	}
}

func TestMessage_Valid_SimpleValid_Success(t *testing.T) {
	m := Message{
		Target:  validPhonenumber,
		Content: validMsg,
	}

	err := m.Valid()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
