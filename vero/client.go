package vero

import (
	"fmt"

	"bytes"
	"encoding/json"
	"github.com/uniplaces/vero-go"
	"io/ioutil"
	"net/http"
)

const baseUrl = "https://api.getvero.com/api/v2/users/%v"

type client struct {
	authToken string
}

type requestData map[string]interface{}

// NewClient returns a new instance of the Vero client
func NewClient(authToken string) vero_go.Client {
	return client{authToken: authToken}
}

func (client) Identify(userId string, data map[string]interface{}, email *string) ([]byte, error) {
	endpoint := client.buildEndpoint("track")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["data"] = data
	if email != nil {
		data["email"] = *email
	}

	return client.send(endpoint, data, http.MethodPost)
}

func (client) Reidentify(userId string, newUserId string) ([]byte, error) {
	endpoint := client.buildEndpoint("reidentify")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["new_id"] = newUserId

	return client.send(endpoint, data, http.MethodPut)
}

func (client) Update(userId string, changes map[string]interface{}) ([]byte, error) {
	endpoint := client.buildEndpoint("edit")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["changes"] = changes

	return client.send(endpoint, data, http.MethodPut)
}

func (client) Tags(userId string, add map[int]string, remove map[int]string) ([]byte, error) {
	endpoint := client.buildEndpoint("edit")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["add"] = add
	requestData["remove"] = remove

	return client.send(endpoint, data, http.MethodPut)
}

func (client) Unsubscribe(userId string) ([]byte, error) {
	endpoint := client.buildEndpoint("unsubscribe")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId

	return client.send(endpoint, data, http.MethodPost)
}

func (client) Resubscribe(userId string) ([]byte, error) {
	endpoint := client.buildEndpoint("resubscribe")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId

	return client.send(endpoint, data, http.MethodPost)
}

func (client) Track(
	eventName string,
	identity map[string]string,
	data map[string]interface{},
	extras map[string]interface{},
) (
	[]byte,
	error,
) {
	endpoint := client.buildEndpoint("track")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["identity"] = identity
	requestData["event_name"] = eventName
	requestData["data"] = data
	requestData["extras"] = extras

	return client.send(endpoint, data, http.MethodPost)
}

func (client) buildEndpoint(endpoint string) string {
	return fmt.Sprintf(baseUrl, endpoint)
}

func (client) send(url string, data map[string]interface{}, method string) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}

	request, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return []byte{}, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	response, err := http.Client{}.Do(request)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
