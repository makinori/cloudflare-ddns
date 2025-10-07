package main

import (
	"errors"
	"io"
	"net/http"
	"regexp"
)

var (
	// https://www.oreilly.com/library/view/regular-expressions-cookbook/9780596802837/ch07s16.html
	validateIPV4Regexp = regexp.MustCompile(
		`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`,
	)
)

func getMyIP(ipv6 bool) (string, error) {
	// alternate method is to dns resolve txt o-o.myaddr.l.google.com
	// http doesnt redirect and returns faster

	var url string
	if ipv6 {
		url = "http://api64.ipify.org?format=text"
	} else {
		url = "http://api.ipify.org?format=text"
	}

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	ipBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	validIPV4 := validateIPV4Regexp.Match(ipBytes)

	if !ipv6 && !validIPV4 {
		return "", errors.New("failed to get valid ipv4: " + string(ipBytes))
	}

	if ipv6 && validIPV4 {
		return "", errors.New("failed to get valid ipv6: " + string(ipBytes))
	}

	return string(ipBytes), nil
}
