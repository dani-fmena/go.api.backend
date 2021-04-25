package auth

import (
	"encoding/base64"
	"errors"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"go.api.backend/schema"
	"go.api.backend/schema/dto"
	"go.api.backend/service/utils"
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

	tkData := base64.StdEncoding.EncodeToString([]byte(v.SisecClientId + ":" + v.SisecClientPass))                                                                          			// preparing SISEC basic authentication token with app (this) credential data, see the doc they deliver to us | FIX think about generate this token on system startup, this way you don't need to generate on every client login
	bodyData := "username=" + url.QueryEscape(uCred.Username) + "&password=" + url.QueryEscape(uCred.Password) + "&domain=" + url.QueryEscape(uCred.Domain) + "&grant_type=password" 	// preparing body with user credentials

	// Building the request for grant intent against SISEC | https://medium.com/rungo/making-external-http-requests-in-go-eb4c015f8839
	req := &http.Request {
		Method: "POST",
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
	defer res.Body.Close()												// ensuring closing the body reader

	// checking what happen with the request
	if err != nil {
		return nil, err, schema.ErrNetwork
	} else if res.StatusCode == iris.StatusNotFound || res.StatusCode == iris.StatusBadRequest {
		return nil, errors.New(schema.ErrDetHttpResError + " - " + strconv.Itoa(res.StatusCode)), schema.ErrHttpResError
	}

	// TODO This code block is temporal, we need to implement the unauthorized (wrong credentials) case, so we need to read the status code response. check the response body
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(3) == 0 { return nil, errors.New(schema.ErrDetInvalidCred), schema.ErrUnauthorized }
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
