package tdapi

import (
	_ "embed"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"testing"
)

// RoundTripFunc is a type representing a function that handles HTTP round-trip.
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip executes the provided function for an HTTP request and returns the response.
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns an *http.Client with the Transport replaced to avoid making real calls. See https://hassansin.github.io/Unit-Testing-http-client-in-Go.
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// mockResponseFunc creates a function that returns a mock HTTP response.
// It expects the request URL and path to a file with the mock response body.
func mockResponseFunc(t *testing.T, expectedURL string, mockBodyFile string) func(req *http.Request) *http.Response {
	file, err := os.Open(mockBodyFile)
	if err != nil {
		t.Fatalf("Cannot open file %q: %v", mockBodyFile, err)
	}

	readCloser := io.NopCloser(file)

	return func(req *http.Request) *http.Response {
		// Check if the request URL matches the expected URL
		if req.URL.Path != expectedURL {
			t.Errorf("Unexpected request URL.\n Got: %s\nWant: %s", req.URL.Path, expectedURL)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       readCloser,
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	}
}

func TestGetAllPersonalLabels(t *testing.T) {
	expectedURL := "/labels"

	// Prepend path for the base URL to the expected URL
	apiURL, err := url.Parse(apiBase)
	if err != nil {
		t.Fatalf("Failed to parse API base URL: %v", err)
	}
	expectedURL = path.Join(apiURL.Path, expectedURL)

	client := NewTestClient(mockResponseFunc(t, expectedURL, "testdata/labels.json"))

	api := TodoistClient{
		httpClient: client,
	}

	labels, err := api.GetAllPersonalLabels()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedLabels := []PersonalLabel{
		{ID: "1", Name: "Label 1", Color: "red", Order: 1, IsFavorite: true},
		{ID: "2", Name: "Label 2", Color: "blue", Order: 2, IsFavorite: false},
	}

	// Compare the returned labels with the expected labels
	if !reflect.DeepEqual(labels, expectedLabels) {
		t.Errorf("Labels do not match.\n Got: %+v\nWant: %+v\n", labels, expectedLabels)
	}
}

func TestGetAllSharedLabels(t *testing.T) {
	expectedURL := "/labels/shared"

	// Prepend path for the base URL to the expected URL
	apiURL, err := url.Parse(apiBase)
	if err != nil {
		t.Fatalf("Failed to parse API base URL: %v", err)
	}
	expectedURL = path.Join(apiURL.Path, expectedURL)

	client := NewTestClient(mockResponseFunc(t, expectedURL, "testdata/shared_labels.json"))

	api := TodoistClient{
		httpClient: client,
	}

	labels, err := api.GetAllSharedLabels()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedLabels := []string{"Shared Label 1", "Shared Label 2"}

	// Compare the returned labels with the expected labels
	if !reflect.DeepEqual(labels, expectedLabels) {
		t.Errorf("Labels do not match.\n Got: %+v\nWant: %+v\n", labels, expectedLabels)
	}
}

func TestGetPersonalLabel(t *testing.T) {
	expectedURL := "/labels/1234567890"

	// Prepend path for the base URL to the expected URL
	apiURL, err := url.Parse(apiBase)
	if err != nil {
		t.Fatalf("Failed to parse API base URL: %v", err)
	}
	expectedURL = path.Join(apiURL.Path, expectedURL)

	client := NewTestClient(mockResponseFunc(t, expectedURL, "testdata/label1.json"))

	api := TodoistClient{
		httpClient: client,
	}

	labels, err := api.GetPersonalLabel("1234567890")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedLabels := PersonalLabel{ID: "1234567890", Name: "Personal_Label_1", Order: 1, Color: "charcoal", IsFavorite: false}

	// Compare the returned labels with the expected labels
	if !reflect.DeepEqual(labels, expectedLabels) {
		t.Errorf("Labels do not match.\n Got: %+v\nWant: %+v\n", labels, expectedLabels)
	}
}
