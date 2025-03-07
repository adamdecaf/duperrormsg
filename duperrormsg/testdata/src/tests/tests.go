package tests

import (
	"errors"
	"fmt"
	"log"
)

func duplicateErrorsNew() {
	// These should be flagged as duplicates
	errors.New("connection failed") // want "duplicate error message"
	errors.New("connection failed") // want "duplicate error message"
}

func duplicateErrorsNewAndErrorf() {
	// These should be flagged as duplicates
	errors.New("validation error") // want "duplicate error message"
	fmt.Errorf("validation error") // want "duplicate error message"
}

func formatStringVariants() {
	// These should be treated as the same message
	fmt.Errorf("user %s not found", "john") // want "duplicate error message"
	fmt.Errorf("user %v not found", "jane") // want "duplicate error message"
}

func uniqueErrors() {
	// These should not be flagged
	errors.New("error one")
	errors.New("error two")
	errors.New("error three")
}

func duplicateLogging() {
	// These should be flagged as duplicates
	log.Printf("failed to process item") // want "duplicate error message"
	log.Printf("failed to process item") // want "duplicate error message"
}

func createCustomError() {
	// Custom error constructor pattern
	NewUserError("invalid input") // want "duplicate error message"
	NewItemError("invalid input") // want "duplicate error message"
}

// Mock functions
func NewUserError(msg string) error {
	return errors.New(msg)
}

func NewItemError(msg string) error {
	return errors.New(msg)
}
