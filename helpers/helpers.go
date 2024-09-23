package helpers

import (
	"errors"
	"net/url"
)

func ValidateURL(urlToValidate string) error {
	parsedURL, err := url.Parse(urlToValidate)
	if err != nil {
		return err
	}

	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return errors.New("invalid url scheme")
	}

	return nil
}
