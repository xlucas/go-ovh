package ovh

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

import "time"

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

// PollTimeshift calculates the difference between
// current system time and remote time through a call
// to API. It may be useful to call this function
// to avoid the signature to be rejected due to
// timeshifting or network delay.
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

// Call is a helper for OVH API interaction
func (c *Client) Call(method, path string, in map[string]interface{}) (map[string]interface{}, error) {
	var inBytes, outBytes []byte
	var err error
	var out map[string]interface{}
	if in != nil {
		inBytes, err = json.Marshal(in)
	}
	if err != nil {
		return nil, err
	}

	// Compute the signature value
	hasher := sha1.New()
	timestamp := strconv.FormatInt(time.Now().Add(c.TimeShift).Unix(), 10)
	url := c.Endpoint + path
	pattern := fmt.Sprintf("%s+%s+%s+%s+%x+%s",
		c.AppSecret,
		c.ConsumerKey,
		method,
		url,
		inBytes,
		timestamp,
	)
	hasher.Write([]byte(pattern))
	signature := fmt.Sprintf("$1$%x", hasher.Sum(nil))

	// Build the request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(inBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Ovh-Application", c.AppKey)
	req.Header.Add("X-Ovh-Consumer", c.ConsumerKey)
	req.Header.Add("X-Ovh-Signature", signature)
	req.Header.Add("X-Ovh-Timestamp", timestamp)

	// Interact with the endpoint
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Parse the result
	outBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(outBytes, &out)

	return out, err
}
