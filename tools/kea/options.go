package kea

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/go-cleanhttp"

	"github.com/sirupsen/logrus"
)

type (
	options struct {
		httpTimeout *int
		logLevel    *logrus.Level
		proxyURL    *string
		auth        *auth
		remote      *string
	}

	// Option : Basic options allowed with this client.
	Option func(*options)
)

// SetHTTPTimeout : This option will set a custom http timeout
func SetHTTPTimeout(timeout int) Option {
	return func(o *options) {
		o.httpTimeout = &timeout
	}
}

// WithLogLevel : set custom client logging level
func WithLogLevel(lvl logrus.Level) Option {
	return func(o *options) {
		o.logLevel = &lvl
	}
}

// WithProxy : Will pass in a proxy URL to the init function.
func WithProxy(url string) Option {
	return func(o *options) {
		o.proxyURL = &url
	}
}

// WithAuth : Will pass in authentication credentials. Default will use ENV vars.
func WithAuth(user, pass string) Option {
	return func(o *options) {
		o.auth = &auth{
			username: user,
			password: pass,
		}
	}
}

// WithRemote : Will set a default remote to use with configuration-backend commands. Default postgresql.
func WithRemote(remote string) Option {
	return func(o *options) {
		o.remote = &remote
	}
}

func (c *Client) processOptions(opts ...Option) {
	o := new(options)
	for _, opt := range opts {
		opt(o)
	}
	transport := cleanhttp.DefaultPooledTransport()
	// nolint: gosec
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = &http.Client{Transport: transport}

	if o.httpTimeout != nil {
		c.client.Timeout = time.Duration(*o.httpTimeout) * time.Second
	}

	c.remote = "postgresql"
	if o.remote != nil {
		c.remote = *o.remote
	}

	if o.proxyURL != nil {
		pURL, err := url.Parse(*o.proxyURL)
		if err != nil {
			logrus.Fatal(err)
		}
		t := cleanhttp.DefaultPooledTransport()
		t.Proxy = http.ProxyURL(pURL)
		c.client.Transport = t
	}

	if o.auth != nil {
		c.auth = auth{
			username: o.auth.username,
			password: o.auth.password,
		}
	} else {
		c.auth = auth{
			username: os.Getenv(envKEAUSER),
			password: os.Getenv(envKEAPASS),
		}
	}
	if c.auth.username == "" || c.auth.password == "" {
		logrus.Fatalf("Missing auth! Use versa.WithAuth() or environment vars %s/%s", envKEAUSER, envKEAPASS)
	}

	c.log = func() *logrus.Logger {
		logger := logrus.New()
		logger.Level = logrus.InfoLevel
		if o.logLevel != nil {
			logger.Level = *o.logLevel
		}
		return logger
	}()
}
