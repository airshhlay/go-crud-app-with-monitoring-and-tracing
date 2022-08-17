package server

import (
	"testing"
)

func BenchmarkTestGetPasswordHash(b *testing.B) {
	handler := Handler{}
	passwordString := "password"
	for i := 0; i < b.N; i++ {
		handler.getPasswordHash(passwordString)
	}
}

func BenchmarkTestCheckPasswordHash(b *testing.B) {
	handler := Handler{}
	passwordString := "password"
	passwordHash, _ := handler.getPasswordHash(passwordString)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.checkPasswordMatch(passwordHash, passwordString)
	}
}
