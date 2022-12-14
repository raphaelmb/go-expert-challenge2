package main

import "testing"

func TestCheckInput(t *testing.T) {
	invalidInputs := []string{"1", "a", "123456789", "1234567", "1234abcd", "12345-1671", "123425-167", "abcdefg"}
	for _, input := range invalidInputs {
		err := IsValidCep(input)
		if err == nil {
			t.Error()
		}
	}
}
