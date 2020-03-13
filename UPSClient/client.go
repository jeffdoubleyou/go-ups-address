// Copyright 2020 Jeffrey Weitz \&lt;jeffdoubleyou@gmail.com\&gt;
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package UPSClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type Client struct {
	accessKey  string
	username   string
	password   string
	httpClient *http.Client
	debug      bool
	testing    bool
}

const (
	API_URL      = "https://onlinetools.ups.com"
	TEST_API_URL = "https://wwwcie.ups.com"
	API_VERSION  = "v1"
)

/*
{
  "UPSSecurity": {
    "UsernameToken": {
      "Username": " username1",
      "Password": " password1"
    },
    "ServiceAccessToken": {
      "AccessLicenseNumber": "AccessLicenseNumber1"
    }
  },
  "LoginAcceptTermsAndConditionsRequest": {
    "Request": ""
  }
}*/

func NewClient(username string, password string, accessKey string) (c *Client) {
	return &Client{accessKey, username, password, &http.Client{}, false, false}
}

func (c *Client) Debug(debug ...bool) bool {
	if len(debug) == 1 {
		c.debug = debug[0]
	}
	return c.debug
}

func (c *Client) Testing(testing ...bool) bool {
	if len(testing) == 1 {
		c.testing = testing[0]
	}
	return c.testing
}

func (c *Client) Timeout(timeout ...int) int {
	if len(timeout) == 1 {
		c.httpClient.Timeout = time.Second * time.Duration(timeout[0])
	}
	return int(c.httpClient.Timeout)
}

func (c Client) dumpRequest(r *http.Request) {
	if r == nil {
		log.Print("dumpReq ok: <nil>")
		return
	}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Print("dumpReq err:", err)
	} else {
		log.Print("dumpReq ok:", string(dump))
	}
}

func (c Client) dumpResponse(r *http.Response) {
	if r == nil {
		log.Print("dumpResponse ok: <nil>")
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Print("dumpResponse err:", err)
	} else {
		log.Print("dumpResponse ok:", string(dump))
	}
}

func (c *Client) Post(resource string, requestName string, request interface{}, requestOption string, queryArgs ...map[string]string) (response []byte, err error) {
	var url string
	apiBase := API_URL

	if c.testing {
		apiBase = TEST_API_URL
	}

	url = fmt.Sprintf("%s/%s/%s", apiBase, resource, API_VERSION)

	if requestOption != "" {
		url += "/" + requestOption
	}

	if len(queryArgs) == 1 {
		query := "?"
		for k, v := range queryArgs[0] {
			query = query + fmt.Sprintf("%s=%s&", k, v)
		}
		url = url + query
	}

	if c.debug {
		log.Printf("POST to %s", url)
	}

	requestBody := make(map[string]interface{})
	requestBody[requestName] = request
	body, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		return
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Username", c.username)
	req.Header.Add("Password", c.password)
	req.Header.Add("AccessLicenseNumber", c.accessKey)

	if c.debug {
		c.dumpRequest(req)
	}

	resp, err := c.httpClient.Do(req)

	if c.debug {
		c.dumpResponse(resp)
	}

	if err != nil {
		return
	}

	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}

	return
}
