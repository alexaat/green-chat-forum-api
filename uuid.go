package main

import (
	"github.com/google/uuid"
)
// generate a new session ID using either a UUID (if it is generated successfully)
// or a combination of the current time in milliseconds and a random integer (if the UUID generation fails)
func generateSessionId() string {
	return uuid.New().String()
}