package client

import (
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/json"
)

type Client struct {
	Server		string
	Path 		string
	ValidatePath	string
}

type TokenReview struct {
	ApiVersion	string        `json:"apiVersion"`
	Kind 		string        `json:"kind"`
	Spec 		Token	      `json:"spec"`
}

type Token struct {
	Token	string        `json:"token"`
}

func (c *Client) Authenticate(username,password string) (token string, err error) {
	url := c.Server + c.Path

	req,_ := http.NewRequest("GET",url,nil)
	req.SetBasicAuth(username,password)

	client := &http.Client{}
	resp,err := client.Do(req)

	if err != nil {
		return
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	token = string(bodyText)
	return
}

const (
	apiVersion = "authentication.k8s.io/v1beta1"
	kind = "TokenReview"
)


func (c *Client) Validate(token string) (trr TokenReviewReturn, err error) {
	url := c.Server + c.ValidatePath

	t := TokenReview{
		ApiVersion: apiVersion,
		Kind: kind,
		Spec: Token{Token: token},
	}

	jsonString,_ := json.MarshalIndent(t,""," ")


	req,err := http.NewRequest("POST",url,bytes.NewBuffer(jsonString))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type","application/json")
	client := &http.Client{}
	resp,err := client.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	body,err := ioutil.ReadAll(resp.Body)

	if err != nil {
	  return
	}

	trr = TokenReviewReturn{}
	err = json.Unmarshal(body,&trr)

	return
}
