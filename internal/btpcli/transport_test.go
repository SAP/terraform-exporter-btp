package btpcli

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBTPCLITransport(t *testing.T) {
	tests := []struct {
		description string

		redirectXtimes int

		returnStatus int

		returnAdditionalHeaders       map[string]string
		returnAdditionalHeadersInCall int
	}{
		{
			description:    "happy path",
			redirectXtimes: 3,
			returnStatus:   http.StatusOK,
		},
		{
			description:                   "happy path - forwards idtoken in case of redirect",
			redirectXtimes:                3,
			returnAdditionalHeadersInCall: 3,
			returnAdditionalHeaders: map[string]string{
				HeaderIDToken: "id-token-val",
			},
			returnStatus: http.StatusOK,
		},
		{
			description:                   "happy path - forwards subdomain in case of redirect",
			redirectXtimes:                3,
			returnAdditionalHeadersInCall: 3,
			returnAdditionalHeaders: map[string]string{
				HeaderCLISubdomain: "another-ga",
			},
			returnStatus: http.StatusOK,
		},
		{
			description:                   "happy path - forwards idtoken and subdomain",
			redirectXtimes:                2,
			returnAdditionalHeadersInCall: 2,
			returnAdditionalHeaders: map[string]string{
				HeaderIDToken:      "id-token-val",
				HeaderCLISubdomain: "another-ga",
			},
			returnStatus: http.StatusOK,
		},
		{
			description:  "error path",
			returnStatus: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			count := 1

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, fmt.Sprintf("/path-%d", count-1), r.URL.Path)

				if test.returnAdditionalHeadersInCall == count {
					for k, v := range test.returnAdditionalHeaders {
						w.Header().Set(k, v)
					}
				}

				if test.returnAdditionalHeadersInCall+1 == count {
					for k, v := range test.returnAdditionalHeaders {
						assert.Equal(t, v, r.Header.Get(k), fmt.Sprintf("value of header '%s' is expected to be '%s', but is '%s'", k, v, r.Header.Get(k)))
					}
				}

				if count < test.redirectXtimes+1 {
					http.RedirectHandler(fmt.Sprintf("/path-%d", count), http.StatusTemporaryRedirect).ServeHTTP(w, r)
					count = count + 1
					return
				}

				w.WriteHeader(test.returnStatus)
			}))
			defer srv.Close()

			client := injectBTPCLITransport(srv.Client(), mustParseURL(srv.URL))

			req, err := http.NewRequest(http.MethodGet, srv.URL+"/path-0", nil)

			assert.NoError(t, err)

			res, err := client.Do(req)

			assert.NoError(t, err)
			assert.Equal(t, test.returnStatus, res.StatusCode)
		})
	}
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func TestBTPCLITransportStripsHeadersOnCrossOriginRedirect(t *testing.T) {
	var otherReceivedIDToken, otherReceivedSubdomain string

	otherHost := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		otherReceivedIDToken = r.Header.Get(HeaderIDToken)
		otherReceivedSubdomain = r.Header.Get(HeaderCLISubdomain)
		w.WriteHeader(http.StatusOK)
	}))
	defer otherHost.Close()

	trusted := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Redirect to a different origin, emitting sensitive headers.
		w.Header().Set(HeaderIDToken, "secret-id-token")
		w.Header().Set(HeaderCLISubdomain, "secret-subdomain")
		http.Redirect(w, r, otherHost.URL+"/landing", http.StatusTemporaryRedirect)
	}))
	defer trusted.Close()

	client := injectBTPCLITransport(trusted.Client(), mustParseURL(trusted.URL))

	req, err := http.NewRequest(http.MethodGet, trusted.URL+"/start", nil)
	req.Header.Set(HeaderIDToken, "secret-id-token")
	req.Header.Set(HeaderCLISubdomain, "secret-subdomain")
	assert.NoError(t, err)

	res, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	assert.Empty(t, otherReceivedIDToken, "X-Id-Token must not be forwarded to a different origin")
	assert.Empty(t, otherReceivedSubdomain, "X-Cpcli-Subdomain must not be forwarded to a different origin")
}
