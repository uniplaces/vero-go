package vero

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientCreation(t *testing.T) {
	t.Parallel()

	NewClient("AUTH_TOKEN")
}

func TestIdentify(t *testing.T) {
	t.Parallel()

	serverMux := http.NewServeMux()
	configureHandler(
		serverMux,
		t,
		http.MethodPost,
		"/users/track",
		`{"auth_token":"AUTH_TOKEN","data":{"First name":"Jeff","Last name":"Kane","email":"jeff@yourdomain.com"},"id":"1234567890"}`,
	)
	server := httptest.NewServer(serverMux)
	defer server.Close()

	data := make(map[string]interface{})
	data["First name"] = "Jeff"
	data["Last name"] = "Kane"

	email := "jeff@yourdomain.com"

	client := getClient(server.URL)

	client.Identify("1234567890", data, &email)
}

func TestReidentify(t *testing.T) {
	t.Parallel()

	serverMux := http.NewServeMux()
	configureHandler(
		serverMux,
		t,
		http.MethodPut,
		"/users/reidentify",
		`{"auth_token":"AUTH_TOKEN","id":"1234567890","new_id":"123"}`,
	)
	server := httptest.NewServer(serverMux)
	defer server.Close()

	client := getClient(server.URL)

	client.Reidentify("1234567890", "123")
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	serverMux := http.NewServeMux()
	configureHandler(
		serverMux,
		t,
		http.MethodPut,
		"/users/edit",
		`{"auth_token":"AUTH_TOKEN","changes":{"city":"lisbon"},"id":"1234567890"}`,
	)
	server := httptest.NewServer(serverMux)
	defer server.Close()

	client := getClient(server.URL)
	changes := make(map[string]interface{})
	changes["city"] = "lisbon"

	client.Update("1234567890", changes)
}

func TestTags(t *testing.T) {
	t.Parallel()

	serverMux := http.NewServeMux()
	configureHandler(
		serverMux,
		t,
		http.MethodPut,
		"/users/tags/edit",
		`{"add":["Blog reader"],"auth_token":"AUTH_TOKEN","id":"1234567890","remove":[]}`,
	)
	server := httptest.NewServer(serverMux)
	defer server.Close()

	client := getClient(server.URL)
	add := []string{"Blog reader"}
	remove := []string{}

	client.Tags("1234567890", add, remove)
}

func TestUnsubscribe(t *testing.T) {
	t.Parallel()

	serverMux := http.NewServeMux()
	configureHandler(
		serverMux,
		t,
		http.MethodPost,
		"/users/unsubscribe",
		`{"auth_token":"AUTH_TOKEN","id":"1234567890"}`,
	)
	server := httptest.NewServer(serverMux)
	defer server.Close()

	client := getClient(server.URL)

	client.Unsubscribe("1234567890")
}

func TestResubscribe(t *testing.T) {
	t.Parallel()

	serverMux := http.NewServeMux()
	configureHandler(
		serverMux,
		t,
		http.MethodPost,
		"/users/resubscribe",
		`{"auth_token":"AUTH_TOKEN","id":"1234567890"}`,
	)
	server := httptest.NewServer(serverMux)
	defer server.Close()

	client := getClient(server.URL)

	client.Resubscribe("1234567890")
}

func TestTrack(t *testing.T) {
	t.Parallel()

	serverMux := http.NewServeMux()
	configureHandler(
		serverMux,
		t,
		http.MethodPost,
		"/events/track",
		`{"auth_token":"AUTH_TOKEN","data":{"city":"Lisbon"},"event_name":"booking-request","extras":{"neighborhood":"Alameda"},"identity":{"email":"example@email.com","id":"123"}}`,
	)
	server := httptest.NewServer(serverMux)
	defer server.Close()

	client := getClient(server.URL)

	identity := make(map[string]string)
	identity["id"] = "123"
	identity["email"] = "example@email.com"

	data := make(map[string]interface{})
	data["city"] = "Lisbon"

	extras := make(map[string]interface{})
	extras["neighborhood"] = "Alameda"

	client.Track("booking-request", identity, data, extras)
}

func assertEqual(t *testing.T, expected string, actual string) {
        if expected != actual {
                t.Errorf(`Expected: %v - Got: %v`, expected, actual)
        }
}

func configureHandler(
	serverMux *http.ServeMux,
	t *testing.T,
	method string,
	endpoint string,
	expectedRequestBody string,
) {
	serverMux.HandleFunc(
		endpoint,
		func(w http.ResponseWriter, r *http.Request) {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			body := buf.String()
                        
                        assertEqual(t, expectedRequestBody, body)
                        assertEqual(t, endpoint, r.URL.Path)
                        assertEqual(t, method, r.Method)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
		},
	)
}

func getClient(baseUrl string) VeroClient {
	return VeroClient{
		authToken: "AUTH_TOKEN",
		baseUrl:   buildBaseUrl(baseUrl),
	}
}

func buildBaseUrl(baseUrl string) string {
	return baseUrl + "/%v"
}
