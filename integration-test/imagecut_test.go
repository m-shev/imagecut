package main

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"io/ioutil"
	"net/http"
	"strings"
)

type ImagecutTest struct {
	ResImgFromOrigin        *http.Response
	ResImgFromCache         *http.Response
	ErrMessage              error
	ImageNotFoundStatusCode int

	response *http.Response
	err      error
	body     string
}

func (i *ImagecutTest) iSendRequestTo(url string) error {

	res, err := http.Get(url)

	i.response = res
	i.err = err

	if res != nil && strings.Contains(res.Header.Get("content-type"), "text/plain") {
		b, err := ioutil.ReadAll(res.Body)
		_ = res.Body.Close()

		if err != nil {
			return err
		}

		i.body = string(b)
	} else {
		i.body = ""
	}

	return nil
}

func (i *ImagecutTest) thereShouldBeStatusCodeAndHeaderEqual(status int, key, value string) error {
	if i.response == nil {
		return i.err
	}

	if i.response.StatusCode != status {
		return fmt.Errorf("expected response status code %d but got %d",
			status, i.response.StatusCode)
	}

	fromCache := i.response.Header.Get(key)

	if fromCache != value {
		return fmt.Errorf("expected header \"%s\" equal \"%s\" but got \"%s\"",
			key, value, fromCache)
	}

	return nil
}

func (i *ImagecutTest) thereShouldBeStatusCode(code int) error {
	if i.response == nil {
		return i.err
	}

	if i.response.StatusCode != code {
		return fmt.Errorf("expected status code \"%d\" but got \"%d\"",
			code, i.response.StatusCode)
	}

	return nil
}

func (i *ImagecutTest) thereShouldBeResponseThatContains(message string) error {

	if !strings.Contains(i.body, message) {
		return fmt.Errorf("expected error message contain: \"%s\", but got %s", message, i.body)
	}

	return nil
}


func FeatureContext(s *godog.Suite) {
	test := ImagecutTest{}
	s.Step(`^I send request to: "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^There should be status code (\d+) and header "([^"]*)" equal "([^"]*)"$`,
		test.thereShouldBeStatusCodeAndHeaderEqual)
	s.Step(`^I send request to: "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^There should be status code (\d+) and header "([^"]*)" equal "([^"]*)"$`,
		test.thereShouldBeStatusCodeAndHeaderEqual)

	s.Step(`^I send request to: "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^There should be status code (\d+)$`, test.thereShouldBeStatusCode)

	s.Step(`^I send request to: "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^There should be response that contains: "([^"]*)"$`,
		test.thereShouldBeResponseThatContains)

	s.Step(`^I send request to: "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^There should be status code (\d+)$`, test.thereShouldBeStatusCode)

	s.Step(`^I send request to: "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^There should be response that contains: "([^"]*)"$`,
		test.thereShouldBeResponseThatContains)
}
