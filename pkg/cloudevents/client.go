package cloudevents

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Builder
	Target    string
	SendError error
}

func NewClient(eventType, source, target string) *Client {
	c := &Client{
		Builder: Builder{
			Source:    source,
			EventType: eventType,
		},
		Target: target,
	}
	return c
}

func (c *Client) RequestSend(data interface{}) (*http.Response, error) {
	req, err := c.Build(c.Target, data)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}

func (c *Client) Send(data interface{}) bool {
	resp, err := c.RequestSend(data)
	if err != nil {
		c.SendError = err
		return false
	}
	if Accepted(resp) {
		c.SendError = nil
		return true
	}
	c.SendError = fmt.Errorf("error sending cloudevent: %s", Status(resp))
	return false
}

func Accepted(resp *http.Response) bool {
	if resp.StatusCode == 204 {
		return true
	}
	return false
}

func Status(resp *http.Response) string {
	if Accepted(resp) {
		return "sent"
	}

	status := resp.Status
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Sprintf("Status[%s] error reading response body: %v", status, err)
	}

	return fmt.Sprintf("Status[%s] %s", status, body)
}