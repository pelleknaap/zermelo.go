package zermelo

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiVersion string = "v3"

type core struct {
	accessCode string
	school string
}

type errorReply struct {
	Response struct {
		Status int
		Message string
	}
}

type accessTokenReply struct {
	AccessToken string `json:"access_token"`
}

func getURI(school string, action string, apiVer string) string {
	return fmt.Sprintf("https://%s.zportal.nl/api/%s/%s", school, apiVer, action)
}

type JSONApiKeyWrapper struct {
	AccessToken string `json:"access_token"`
}

func (c *core) GetAccessToken(authCode string) error {
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", authCode)

	fmt.Println(form.Encode())

	var d accessTokenReply

	err := c.request(http.MethodPost, "oauth/token?grant_type=authorization_code&code=" + authCode, &d, nil)
	if err != nil {
		return err
	}

	fmt.Println(d)

	c.accessCode = d.AccessToken

	return nil
}

func (c core) Get(action string, d interface{}) error {
	return c.request(http.MethodGet, action, d, nil)
}

func (c core) request(method, action string, d interface{}, reader io.Reader) error {
	fmt.Println(getURI(c.school, action, apiVersion))
	req, err := http.NewRequest(method, getURI(c.school, action, apiVersion), reader)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "zermelo-api/v1")
	req.Header.Set("Authorization", "Bearer "+c.accessCode)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := resp.Body.Close(); err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		if len(data) > 0 {
			r := &errorReply{}
			err := json.Unmarshal(data, r)
			if err == nil && len(r.Response.Message) > 0 {
				return errors.New(fmt.Sprintf("%d: %s", r.Response.Status, r.Response.Message))
			}
		}
		return errors.New(resp.Status)
	}

	return json.Unmarshal(data, &d)
}
