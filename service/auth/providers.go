package auth

import (
	"encoding/base64"
	"errors"
	"go.api.backend/schema"
	"go.api.backend/service/utils"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"go.api.backend/schema/dto"
)


type Provider interface {
	GrantIntent(userCredential *dto.UserCredIn, data interface{}) (*dto.SISECGrantIntentIn, error, string)
}

// region ======== SISEC AUTHENTICATION PROVIDER =========================================

type ProviderSisec struct {
	URL        *url.URL
	ClientId   string
	ClientPass string
}

// password grant type

// GrantIntent make login / grant (password type) intent against SISEC Tecnomática auth system, using
// the given user credentials.
//
// It returns XX, error if any, and a error code to identify the kind of error on the caller.
//
// - uCred [*dto.UserCredIn] ~ User credential
//
// - options [interface{}] ~ Some options or data of a specific type to be used in the provider authentication method
func (p *ProviderSisec) GrantIntent(uCred *dto.UserCredIn, options interface{}) (*dto.SISECGrantIntentIn, error, string) {

	v, ok := options.(*utils.SvcConfig)						// Checking the options | Assertion check
	if !ok { return nil, nil, schema.ErrInvalidType }

	tkData := base64.StdEncoding.EncodeToString([]byte(v.SisecClientId + ":" + v.SisecClientPass))                                                                          			// preparing sisec basic authentication token, see the doc they deliver to us
	bodyData := "username=" + url.QueryEscape(uCred.Username) + "&password=" + url.QueryEscape(uCred.Password) + "&domain=" + url.QueryEscape(uCred.Domain) + "&grant_type=password" 	// preparing body with user credentials

	// Building the request for grant intent against SISEC | https://medium.com/rungo/making-external-http-requests-in-go-eb4c015f8839
	req := &http.Request {
		Method: "GET",	// TODO we need to do a POST method and use the user credential info in the request against SISEC
		URL: p.URL,
		Header: map[string][]string {
			"Content-Type": {"application/x-www-form-urlencoded"},
			"Accept-Encoding": {"gzip, deflate, br"},
			"Authorization": {"Basic " + tkData},
		},
		Body: ioutil.NopCloser(strings.NewReader(bodyData)),
	}

	// requesting the auth, access grant
	res, err := http.DefaultClient.Do(req)
	if err != nil { return nil, err, schema.ErrNetwork }
	defer res.Body.Close()												// ensuring closing the body reader

	// TODO This code block is temporal, we need to implement the unauthorized (wrong credentials) case, so we need to read the status code response
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(3) == 0 { return nil, errors.New(schema.ErrInvalidCred), schema.ErrUnauthorized }
	// see crypto/rand for more secure generations

	// Parsing and unmarshalling the response options
	grantData := &dto.SISECGrantIntentIn{} // new(dto.SISECGrantIntentIn)
	if e := jsoniter.NewDecoder(res.Body).Decode(grantData); e != nil {
		{ return nil, e, schema.ErrJsonParse }
	}

	// Generating the options to be tokenized
	return grantData, nil, ""
}
// endregion =============================================================================
