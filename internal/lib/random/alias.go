package random

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func GenerateAlias(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("error with parse URL: %w", err)
	}

	host := parsedURL.Host

	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimSuffix(host, ".com")

	reg := regexp.MustCompile("[^a-zA-Z0-9-]+")
	host = reg.ReplaceAllString(host, "")

	host = strings.ToLower(host)

	return host, nil
}
