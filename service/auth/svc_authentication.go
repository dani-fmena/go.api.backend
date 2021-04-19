package auth

import (
	"net/url"

	"go.api.backend/service/utils"
)

type SvcAuthentication struct {
	AuthProviders map[string]Provider 			// similar to slices, maps are reference types.
}

// NewSvcAuthentication creates the authentication service. It provides the methods to make the
// authentication intent with the register providers.
//
// - providers [Array] ~ Maps of providers string token / identifiers
//
// - svcConfig [*SvcConfig] ~ App conf instance pointer
func NewSvcAuthentication(providers map[string]bool, svcConfig *utils.SvcConfig) *SvcAuthentication {

	k := &SvcAuthentication{AuthProviders: make(map[string]Provider)}

	for v, _ := range providers {

		if v == "sisec" {							// ===== SISEC CASE =======
			_url, err := url.Parse(svcConfig.SisecUrl)
			if err != nil { panic(err) }

			(*k).AuthProviders[v] = &ProviderSisec {
				URL:        _url,
				ClientId:   svcConfig.SisecClientId,
				ClientPass: svcConfig.SisecClientPass,
			}
		} else if v == "default" {					// ===== DEFAULT CASE, NORMAL DATABASE LOGIN =======
			// TODO implement normal database login authentication case
		}
	}

	return k
}
