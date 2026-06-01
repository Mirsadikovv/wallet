package utils

import (
	"math/rand"
	"strconv"
	"time"
)

type otpEntry struct {
	code      string
	expiresAt time.Time
}

type OTPCache struct {
	m map[string]otpEntry
}

func NewOTPCache() *OTPCache {
	return &OTPCache{m: make(map[string]otpEntry)}
}

func (o *OTPCache) Generate(email string) string {
	code := strconv.Itoa(100000 + rand.Intn(900000))
	o.m[email] = otpEntry{
		code:      code,
		expiresAt: time.Now().Add(5 * time.Minute),
	}
	return code
}

func (o *OTPCache) Verify(email, code string) bool {
	entry, ok := o.m[email]
	if !ok {
		return false
	}
	if time.Now().After(entry.expiresAt) {
		delete(o.m, email)
		return false
	}
	if entry.code != code {
		return false
	}
	delete(o.m, email)
	return true
}
