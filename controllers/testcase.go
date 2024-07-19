package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/low4ey/OJ/Golang-backend/models"
)

func getTestCases(url string) ([]models.TestCase, error) {
	// Create an HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the response into a slice of TestCase structs
	var testCases []models.TestCase
	err = json.Unmarshal(body, &testCases)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}
	return testCases, nil
}

func getSampleTestCase(url string) ([]models.TestCase, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var sampleTestCase []models.TestCase
	err = json.Unmarshal(body, &sampleTestCase)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	if len(sampleTestCase) > 0 {
		return []models.TestCase{sampleTestCase[0]}, nil
	}

	return nil, fmt.Errorf("no test cases found in the response")
}
