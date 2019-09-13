package vero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/uniplaces/vero-go"
)

const baseUrl = "https://api.getvero.com/api/v2/%v"

// VeroClient manages requests to the Vero API
type VeroClient struct {
	authToken string
	baseUrl   string
}

type requestData map[string]interface{}

// NewClient returns a new instance of the Vero client
func NewClient(authToken string) vero_go.Client {
	return VeroClient{authToken: authToken, baseUrl: baseUrl}
}

// Identify creates a new user profile if the user doesn’t exist yet
func (client VeroClient) Identify(userId string, data map[string]interface{}, email *string) ([]byte, error) {
	endpoint := client.buildEndpoint("users/track")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["data"] = data
	if email != nil {
		data["email"] = *email
	}

	return client.send(endpoint, requestData, http.MethodPost)
}

// Reidentify updates user id from existing user
func (client VeroClient) Reidentify(userId string, newUserId string) ([]byte, error) {
	endpoint := client.buildEndpoint("users/reidentify")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["new_id"] = newUserId

	return client.send(endpoint, requestData, http.MethodPut)
}

// Update updates information from existing user
func (client VeroClient) Update(userId string, changes map[string]interface{}) ([]byte, error) {
	endpoint := client.buildEndpoint("users/edit")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["changes"] = changes

	return client.send(endpoint, requestData, http.MethodPut)
}

// Tags lets you add or remove tags to or from any of your users
func (client VeroClient) Tags(userId string, add []string, remove []string) ([]byte, error) {
	endpoint := client.buildEndpoint("users/tags/edit")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId
	requestData["add"] = add
	requestData["remove"] = remove

	return client.send(endpoint, requestData, http.MethodPut)
}

// Unsubscribe unsubscribes a single user
func (client VeroClient) Unsubscribe(userId string) ([]byte, error) {
	endpoint := client.buildEndpoint("users/unsubscribe")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId

	return client.send(endpoint, requestData, http.MethodPost)
}

// Resubscribe lets you resubscribe a single user
func (client VeroClient) Resubscribe(userId string) ([]byte, error) {
	endpoint := client.buildEndpoint("users/resubscribe")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["id"] = userId

	return client.send(endpoint, requestData, http.MethodPost)
}

// Track endpoint tracks an event for a specific customer. If the customer profile doesn’t exist, Vero will create it
func (client VeroClient) Track(
	eventName string,
	identity map[string]string,
	data map[string]interface{},
	extras map[string]interface{},
) (
	[]byte,
	error,
) {
	endpoint := client.buildEndpoint("events/track")

	requestData := requestData{}
	requestData["auth_token"] = client.authToken
	requestData["identity"] = identity
	requestData["event_name"] = eventName
	requestData["data"] = data
	requestData["extras"] = extras

	return client.send(endpoint, requestData, http.MethodPost)
}

func (client VeroClient) buildEndpoint(endpoint string) string {
	return fmt.Sprintf(client.baseUrl, endpoint)
}

func (VeroClient) send(url string, data map[string]interface{}, method string) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}

	request, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return []byte{}, err
	}

	request.Close = true
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return []byte{}, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	if response.StatusCode != http.StatusOK {
		return []byte{}, createErrorWithHTTPResponseMessage(body)
	}

	return body, nil
}

func createErrorWithHTTPResponseMessage(responseBody []byte) error {
	var responseData map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return err
	}

	return errors.New(fmt.Sprintf("%s", responseData["message"]))
}
