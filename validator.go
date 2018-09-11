package main

type Validator struct {
}

func (v *Validator) isLegalOriginator(originator string) bool {
	return true
}

func (v *Validator) validate(rawMessage *RawMessage) error {
	return nil
}
