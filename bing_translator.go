package main

import (
	_ "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"unicode/utf8"
)

type Config struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
	GrantType    string `json:"grant_type"`
}

var (
	Accessuri    = "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13"
	TranslateURI = "http://api.microsofttranslator.com/v2/Http.svc/Translate?text=%s&from=%s&to=%s"
	ClientSeret  = "PnBcvuSKqquTsqc7Zg5M/ad7B39swY03Uf5l0PJSRL8="
	ClientId     = "01635825664"
	Scope        = "http://api.microsofttranslator.com"
	GrantType    = "client_credentials"
)

type Token struct {
	TokenType   string `json:"token_type`
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	Scope       string `json:"scope"`
}

type TranslatorConfig struct {
	From        string
	To          string
	AccessToken string
}

func main() {

	// data, err := json.Marshal(InputParam)
	// body := bytes.NewBuffer(data)

	tokenChanOut := GetAccessToken(Accessuri, Config{ClientId, ClientSeret, Scope, GrantType})

	// DO SomeThing
	token := <-tokenChanOut

	value := BingTranslator("hoàng hôn", TranslatorConfig{
		From:        "vi",
		To:          "en",
		AccessToken: token.AccessToken,
	})

	fmt.Println(value)

}

func GetAccessToken(access_uri string, conf Config) <-chan Token {

	tokenChan := make(chan Token)

	go func() {
		values := url.Values{}

		values["client_secret"] = []string{ClientSeret}
		values["client_id"] = []string{ClientId}
		values["scope"] = []string{Scope}
		values["grant_type"] = []string{GrantType}

		resp, err := http.PostForm(access_uri, values)
		PannicOnError(err)

		var token Token
		data, err := ioutil.ReadAll(resp.Body)
		PannicOnError(err)

		err = json.Unmarshal(data, &token)
		PannicOnError(err)

		tokenChan <- token
	}()

	return tokenChan

}

func BingTranslator(text string, config TranslatorConfig) interface{} {

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf(TranslateURI, text, config.From, config.To), nil)
	req.Header.Add("Authorization", "Bearer"+" "+config.AccessToken)
	req.Header.Add("charset", "UTF-8")

	resp, err := client.Do(req)
	PannicOnError(err)

	data, err := ioutil.ReadAll(resp.Body)
	PannicOnError(err)

	return string(data)
}

func PannicOnError(e error) {
	if e != nil {
		panic(e)
	}
}
