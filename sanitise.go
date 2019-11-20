package cas

import (
	"net/url"
)

var (
	urlCleanParameters = []string{"gateway", "renew", "service", "ticket"}
)

// sanitisedURL cleans a LoginURL of CAS specific parameters
func sanitisedURL(unclean *url.URL) *url.URL {
	// Shouldn't be any errors parsing an existing *url.LoginURL
	u, _ := url.Parse(unclean.String())
	q := u.Query()

	for _, param := range urlCleanParameters {
		q.Del(param)
	}

	u.RawQuery = q.Encode()
	return u
}

// sanitisedURLString cleans a LoginURL and returns its string value
func sanitisedURLString(unclean *url.URL) string {
	return sanitisedURL(unclean).String()
}
