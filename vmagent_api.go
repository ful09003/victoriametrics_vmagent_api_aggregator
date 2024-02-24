package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func fetchVMAgentTargets(c *http.Client, r *http.Request) (VMAgentAPIResponse, error) {
	vmRes := VMAgentAPIResponse{}

	resp, err := c.Do(r)
	if err != nil {
		return vmRes, nil
	}
	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return vmRes, nil
	}

	if err := json.Unmarshal(resBytes, &vmRes); err != nil {
		return VMAgentAPIResponse{}, err
	}

	return vmRes, nil
}
