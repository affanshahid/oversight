package httpprobe

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/affanshahid/oversight/prober/probe"
	"github.com/affanshahid/oversight/prober/registrar"
	"github.com/go-redis/redis/v7"
)

// HTTPOptions are the HTTPProbe specific options
type HTTPOptions struct {
	Method string
	URL    string
}

// HTTPProbe fetches data from remote URLs
type HTTPProbe struct {
	*probe.BaseProbe
	parsedOpts HTTPOptions
}

// Fetch sends a request to a remote URL
func (p *HTTPProbe) Fetch() (interface{}, error) {
	req, err := http.NewRequest(
		p.parsedOpts.Method,
		p.parsedOpts.URL,
		nil,
	)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// New creates a new HTTPProbe
func New(config *probe.Config, redisClient *redis.Client) (probe.Probe, error) {
	var opts HTTPOptions

	err := json.Unmarshal(config.Options.RawMessage, &opts)

	if err != nil {
		return nil, err
	}

	return &HTTPProbe{
		BaseProbe: &probe.BaseProbe{
			Config:      config,
			RedisClient: redisClient,
		},
		parsedOpts: opts,
	}, nil
}

func init() {
	registrar.Register("http", New)
}
