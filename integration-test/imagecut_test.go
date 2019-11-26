package main

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"net/http"
)

type ImagecutTest struct {
	ResponseFromOrigin *http.Response
}

func (i *ImagecutTest) iSendRequestToImagecutService(url string) error {
	res, err := http.Get(url)

	if err != nil {
		return err
	}

	i.ResponseFromOrigin = res

	return err
}

func (i *ImagecutTest) responseStatusIsAndHeaderEqual(status int, headerKey, headerValue string ) error {
	if i.ResponseFromOrigin.StatusCode != status {
		return fmt.Errorf("expected response from origin %d but got %d",
			status, i.ResponseFromOrigin.StatusCode)
	}

	return nil
}



func FeatureContext(s *godog.Suite) {
	test := ImagecutTest{}
	s.Step(`^I send request to imagecut service: "([^"]*)"$`, test.iSendRequestToImagecutService)
	s.Step(`^response status is "([^"]*)" and header "([^"]*)" equal "([^"]*)"$`,
		test.responseStatusIsAndHeaderEqual)
}
