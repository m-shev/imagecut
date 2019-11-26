package main

import (
	"github.com/DATA-DOG/godog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "pretty",
		Paths:     []string{"features"},
		Randomize: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
