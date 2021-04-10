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
// - providers [Array] ~ Array of providers string token / identifiers
//
// - appConf [*SvcConfig] ~ App conf instance pointer
func NewSvcAuthentication(providers []string, c *utils.SvcConfig) *SvcAuthentication {

	k := &SvcAuthentication{AuthProviders: make(map[string]Provider)}

	for _, v := range providers {

		if v == "sisec" {							// ===== SISEC CASE =======
			_url, err := url.Parse(c.SisecUrl)
			if err != nil { panic(err) }

			(*k).AuthProviders[v] = &ProviderSisec {
				URL:        _url,
				ClientId:   c.SisecClientId,
				ClientPass: c.SisecClientPass,
			}
		} else if v == "default" {					// ===== DEFAULT CASE =======
			// TODO implement normal database login authentication case
		}
	}

	return k
}
