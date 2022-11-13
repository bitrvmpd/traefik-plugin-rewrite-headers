//nolint
package traefik_plugin_rewrite_headers

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Rewrite holds one rewrite body configuration.
type Rewrite struct {
	Header      string `json:"header,omitempty"`
	Regex       string `json:"regex,omitempty"`
	Replacement string `json:"replacement,omitempty"`
}

// Config holds the plugin configuration.
type Config struct {
	Rewrites Rewrites `json:"rewrites,omitempty"`
}

type Rewrites struct {
	Request  []Rewrite `json:"request,omitempty"`
	Response []Rewrite `json:"response,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type rewrite struct {
	header      string
	regex       *regexp.Regexp
	replacement string
}

type rewriteBody struct {
	name      string
	next      http.Handler
	responses []rewrite
	requests  []rewrite
}

// New creates and returns a new rewrite body plugin instance.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	responses := make([]rewrite, len(config.Rewrites.Response))
	requests := make([]rewrite, len(config.Rewrites.Request))

	for i, rewriteConfig := range config.Rewrites.Response {
		regex, err := regexp.Compile(rewriteConfig.Regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", rewriteConfig.Regex, err)
		}

		responses[i] = rewrite{
			header:      rewriteConfig.Header,
			regex:       regex,
			replacement: rewriteConfig.Replacement,
		}
	}

	for i, rewriteConfig := range config.Rewrites.Request {
		regex, err := regexp.Compile(rewriteConfig.Regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", rewriteConfig.Regex, err)
		}

		requests[i] = rewrite{
			header:      rewriteConfig.Header,
			regex:       regex,
			replacement: rewriteConfig.Replacement,
		}
	}

	return &rewriteBody{
		name:      name,
		next:      next,
		responses: responses,
		requests:  requests,
	}, nil
}

func (r *rewriteBody) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrappedWriter := &responseWriter{
		writer:   rw,
		rewrites: r.responses,
	}

	// Modify requests headers before passing the response writer.
	for _, rewrite := range r.requests {

		headers := req.Header.Values(rewrite.header)

		if len(headers) == 0 {
			continue
		}

		req.Header.Del(rewrite.header)

		for _, header := range headers {
			value := rewrite.regex.ReplaceAllString(header, rewrite.replacement)
			req.Header.Add(rewrite.header, strings.TrimSpace(value))
		}
	}

	r.next.ServeHTTP(wrappedWriter, req)
}

type responseWriter struct {
	writer   http.ResponseWriter
	rewrites []rewrite
}

func (r *responseWriter) Header() http.Header {
	return r.writer.Header()
}

func (r *responseWriter) Write(bytes []byte) (int, error) {
	return r.writer.Write(bytes)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	for _, rewrite := range r.rewrites {
		headers := r.writer.Header().Values(rewrite.header)

		if len(headers) == 0 {
			continue
		}

		r.writer.Header().Del(rewrite.header)

		for _, header := range headers {
			value := rewrite.regex.ReplaceAllString(header, rewrite.replacement)
			r.writer.Header().Add(rewrite.header, value)
		}
	}

	r.writer.WriteHeader(statusCode)
}
