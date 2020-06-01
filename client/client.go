package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client token, host, htpp.Client
type Client struct {
	Token  string
	Host   string
	Client *http.Client
}

// RequestOptions  path, method, etc
type RequestOptions struct {
	Path   string
	Method string
	Body   []byte
	QS     map[string]string
}

// NewClient returns a new client configured to communicate on a server with the
// given hostname and to send an Authorization Header with the value of
// token
func NewClient(hostname string, token string) *Client {
	return &Client{
		Host:   hostname,
		Token:  token,
		Client: &http.Client{},
	}

}

// RequestAPI http request to Codefresh API
func (client *Client) RequestAPI(opt *RequestOptions) ([]byte, error) {
	finalURL := fmt.Sprintf("%s%s", client.Host, opt.Path)
	if opt.QS != nil {
		finalURL += ToQS(opt.QS)
	}
	request, err := http.NewRequest(opt.Method, finalURL, bytes.NewBuffer(opt.Body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", client.Token)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Client.Do(request)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body %v %v", resp.StatusCode, resp.Status)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%v, %s", resp.Status, string(body))
	}
	return body, nil
}

// ToQS add extra parameters to path
func ToQS(qs map[string]string) string {
	var arr = []string{}
	for k, v := range qs {
		arr = append(arr, fmt.Sprintf("%s=%s", k, v))
	}
	return "?" + strings.Join(arr, "&")
}

// DecodeResponseInto json Unmarshall
func DecodeResponseInto(body []byte, target interface{}) error {
	return json.Unmarshal(body, target)
}

// EncodeToJSON json Marshal
func EncodeToJSON(object interface{}) ([]byte, error) {
	body, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return body, nil
}
