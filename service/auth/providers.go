package auth

import (
	"go.api.backend/schema"
	"net/http"
	"net/url"

	jsoniter "github.com/json-iterator/go"

	"go.api.backend/schema/dto"
)


type Provider interface {
	GrantIntent(userCredential *dto.UserCredIn) (*dto.SISECGrantIntentIn, error, string)
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
// - cred [*dto.UserCredIn] ~ User credential
func (p *ProviderSisec) GrantIntent(cred *dto.UserCredIn) (*dto.SISECGrantIntentIn, error, string) {

	// Building the request for grant intent against SISEC
	req := &http.Request {
		Method: "GET",	// TODO we need to do a POST method and use the user credential info in the request against SISEC
		URL: p.URL,
		Header: map[string][]string {
			"Content-Type": {"application/x-www-form-urlencoded"},
			"Accept-Encoding": {"gzip, deflate, br"},
			// TODO create the SISEC authorization header, see the doc file
		},
	}

	// Doing the request
	res, err := http.DefaultClient.Do(req)
	if err != nil { return nil, err, schema.ErrNetwork }

	defer res.Body.Close()								// ensuring closing the body reader

	// Parsing and unmarshalling the response data
	grantData := &dto.SISECGrantIntentIn{} 				// new(dto.SISECGrantIntentIn)
	if e := jsoniter.NewDecoder(res.Body).Decode(grantData); e != nil {
		{ return nil, e, schema.ErrJsonParse }
	}

	return grantData, nil, ""
}
// endregion =============================================================================
