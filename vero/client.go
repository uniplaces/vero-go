package vero

import (
	"fmt"

	"bytes"
	"encoding/json"
	"github.com/uniplaces/vero-go"
	"io/ioutil"
	"net/http"
)

const (
	usersBaseUrl  = "https://api.getvero.com/api/v2/users/%v"
	eventsBaseUrl = "https://api.getvero.com/api/v2/events/%v"
)

type client struct {
	authToken string
}

type requestData map[string]interface{}

// NewClient returns a new instance of the Vero client
func NewClient(authToken string) vero_go.Client {
	return client{authToken: authToken}
}

func (client client) Identify(userId string, data map[string]interface{}, email *string) ([]byte, error) {
	endpoint := client.buildEndpoint(usersBaseUrl, "track")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["data"] = data
	if email != nil {
		data["email"] = *email
	}

	return client.send(endpoint, requestData, http.MethodPost)
}

func (client client) Reidentify(userId string, newUserId string) ([]byte, error) {
	endpoint := client.buildEndpoint(usersBaseUrl, "reidentify")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["new_id"] = newUserId

	return client.send(endpoint, requestData, http.MethodPut)
}

func (client client) Update(userId string, changes map[string]interface{}) ([]byte, error) {
	endpoint := client.buildEndpoint(usersBaseUrl, "edit")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["changes"] = changes

	return client.send(endpoint, requestData, http.MethodPut)
}

func (client client) Tags(userId string, add []string, remove []string) ([]byte, error) {
	endpoint := client.buildEndpoint(usersBaseUrl, "tags/edit")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["add"] = add
	requestData["remove"] = remove

	return client.send(endpoint, requestData, http.MethodPut)
}

func (client client) Unsubscribe(userId string) ([]byte, error) {
	endpoint := client.buildEndpoint(usersBaseUrl, "unsubscribe")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId

	return client.send(endpoint, requestData, http.MethodPost)
}

func (client client) Resubscribe(userId string) ([]byte, error) {
	endpoint := client.buildEndpoint(usersBaseUrl, "resubscribe")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId

	return client.send(endpoint, requestData, http.MethodPost)
}

func (client client) Track(
	eventName string,
	identity map[string]string,
	data map[string]interface{},
	extras map[string]interface{},
) (
	[]byte,
	error,
) {
	endpoint := client.buildEndpoint(eventsBaseUrl, "track")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["identity"] = identity
	requestData["event_name"] = eventName
	requestData["data"] = data
	requestData["extras"] = extras

	return client.send(endpoint, requestData, http.MethodPost)
}

func (client) buildEndpoint(baseUrl string, endpoint string) string {
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

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
