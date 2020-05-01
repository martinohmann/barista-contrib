package ipify

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/martinohmann/barista-contrib/modules/ip"
)

// New create a new *ip.Module using https://ipify.org to look up the current
// public ip address.
func New() *ip.Module {
	return ip.New(Provider)
}

// Provider is an ip.Provider which retrieves the public ip via
// https://api.ipify.org.
var Provider = ip.ProviderFunc(func() (net.IP, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.ipify.org", nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		if netErr, ok := err.(net.Error); ok {
			if netErr.Temporary() || netErr.Timeout() {
				// Transient errors and timeouts indicate that we are offline.
				return nil, nil
			}
		}

		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(string(buf)), nil
})
