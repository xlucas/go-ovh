package ovh

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// OVH endpoints list
const (
	ENDPOINT_CA_OVHCOM     = "https://ca.api.ovh.com/1.0"
	ENDPOINT_CA_KIMSUFI    = "https://ca.api.kimsufi.com/1.0"
	ENDPOINT_CA_RUNABOVE   = "https://api.runabove.com/1.0"
	ENDPOINT_CA_SOYOUSTART = "https://ca.api.soyoustart.com/1.0"
	ENDPOINT_EU_OVHCOM     = "https://eu.api.ovh.com/1.0"
	ENDPOINT_EU_KIMSUFI    = "https://eu.api.kimsufi.com/1.0"
	ENDPOINT_EU_RUNABOVE   = "https://api.runabove.com/1.0"
	ENDPOINT_EU_SOYOUSTART = "https://eu.api.soyoustart.com/1.0"
)

// Client helps interacting with OVH API endpoints.
type Client struct {
	AppKey      string
	AppSecret   string
	ConsumerKey string
	Endpoint    string
	TimeShift   time.Duration
}

// NewClient builds up a new client link to the specified endpoint
// with given authentication information and no timeshift.
func NewClient(endpoint, ak, as, ck string) *Client {
	return &Client{
		AppKey:      ak,
		AppSecret:   as,
		ConsumerKey: ck,
		Endpoint:    endpoint,
		TimeShift:   0,
	}
}

func computeSignature(appSecret, consumerKey, method, url string, body []byte, timestamp int64) string {
	hasher := sha1.New()
	pattern := fmt.Sprintf("%s+%s+%s+%s+%x+%d",
		appSecret,
		consumerKey,
		method,
		url,
		body,
		timestamp)
	hasher.Write([]byte(pattern))
	return fmt.Sprintf("$1$%x", hasher.Sum(nil))
}

func sendRequest(appKey, consumerKey, signature string, timestamp int64, method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Ovh-Application", appKey)
	req.Header.Add("X-Ovh-Consumer", consumerKey)
	req.Header.Add("X-Ovh-Signature", signature)
	req.Header.Add("X-Ovh-Timestamp", fmt.Sprintf("%d", timestamp))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Unexpected HTTP return code (%s).", resp.Status))
	}

	outBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return outBytes, err
}

// PollTimeshift calculates the difference between
// local and remote system time through a call to
// the API. It may be useful to call this function
// to avoid the signature to be rejected due to
// timeshift or network delay.
func (c *Client) PollTimeshift() error {
	sysTime := time.Now()
	resp, err := http.Get(c.Endpoint + "/auth/time")
	if err != nil {
		return err
	}
	outPayload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	apiTime, err := strconv.ParseInt(string(outPayload), 10, 64)
	if err != nil {
		return err
	}
	c.TimeShift = time.Unix(apiTime, 0).Sub(sysTime)
	return err
}

// Call is a helper for OVH API interaction that returns and use interfaces
func (c *Client) Call(method, path string, in interface{}, out interface{}) error {
	var inBytes, outBytes []byte
	var err error
	if in != nil {
		inBytes, err = json.Marshal(in)
	}
	if err != nil {
		return err
	}

	url := c.Endpoint + path
	timestamp := time.Now().Add(c.TimeShift).Unix()
	signature := computeSignature(c.AppSecret, c.ConsumerKey, method, url, inBytes, timestamp)

	outBytes, err = sendRequest(c.AppKey, c.ConsumerKey, signature, timestamp, method, url, inBytes)
	if err != nil {
		return err
	}

	json.Unmarshal(outBytes, &out)
	return nil
}
