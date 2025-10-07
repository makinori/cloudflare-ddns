package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type UnifiNetworkListUpdate struct {
	GroupMembers []string `json:"group_members"`
}

func unifiUpdateNetworkList(
	gatewayIP string, listID string, addresses []string, token string,
) error {
	data, err := json.Marshal(UnifiNetworkListUpdate{
		GroupMembers: addresses,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf(
			"https://%s/proxy/network/api/s/default/rest/firewallgroup/%s",
			gatewayIP, listID,
		),
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-API-KEY", token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	output, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return errors.New(string(output))
}
