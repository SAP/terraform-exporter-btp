package btpcli

import (
	"net/http"
	"net/url"
)

// injectBTPCLITransport wraps the transport of the given client with the btpcliTransport.
// trustedServerURL identifies the host that session/identity headers may be sent to.
func injectBTPCLITransport(client *http.Client, trustedServerURL *url.URL) *http.Client {
	parentTransport := http.DefaultTransport

	if client.Transport != nil {
		parentTransport = client.Transport
	}

	var trustedHost string
	if trustedServerURL != nil {
		trustedHost = trustedServerURL.Host
	}

	client.Transport = &btpcliTransport{transport: parentTransport, trustedHost: trustedHost}
	return client
}

// btpcliTransport implements the http.RoundTripper interface. Its purpose is to copy headers from parent responses (in case of redirects)
// to the actual request, scoped to the trusted server host so that session and
// identity headers are never forwarded to a different origin on redirects.
type btpcliTransport struct {
	transport   http.RoundTripper
	trustedHost string
}

// copyResponseHeaderToRequestHeader copies a header from the parent response (source) to the request (target).
func (bt *btpcliTransport) copyResponseHeaderToRequestHeader(req *http.Request, source string, target string) {
	if req.Response == nil {
		return
	}

	if value := req.Response.Header.Get(source); len(value) > 0 {
		req.Header.Set(target, value)
	}
}

// isTrustedHost reports whether the request targets the trusted server host.
// When no trusted host is known the request is treated as trusted to preserve
// behaviour for callers that do not supply a server URL.
func (bt *btpcliTransport) isTrustedHost(req *http.Request) bool {
	if bt.trustedHost == "" {
		return true
	}
	return req.URL != nil && req.URL.Host == bt.trustedHost
}

func (bt *btpcliTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if bt.isTrustedHost(req) {
		bt.copyResponseHeaderToRequestHeader(req, HeaderCLISubdomain, HeaderCLISubdomain)
		bt.copyResponseHeaderToRequestHeader(req, HeaderIDToken, HeaderIDToken)
	} else {
		// Never forward session/identity headers to a host other than the
		// trusted server, e.g. after a cross-origin redirect.
		req.Header.Del(HeaderCLISubdomain)
		req.Header.Del(HeaderIDToken)
		req.Header.Del(HeaderCLISessionId)
	}

	return bt.transport.RoundTrip(req)
}
