//nolint
package traefik_plugin_rewrite_headers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTPRequestResponse(t *testing.T) {
	tests := []struct {
		desc      string
		responses struct {
			rewrite       []Rewrite
			reqHeader     http.Header
			expRespHeader http.Header
		}
		requests struct {
			rewrite       []Rewrite
			reqHeader     http.Header
			expRespHeader http.Header
		}
	}{
		{
			desc: "should replace foo by bar in location request and response header",
			responses: struct {
				rewrite       []Rewrite
				reqHeader     http.Header
				expRespHeader http.Header
			}{
				[]Rewrite{
					{
						Header:      "Location",
						Regex:       "foo",
						Replacement: "bar",
					},
				},
				map[string][]string{
					"Location": {"foo", "anotherfoo"},
				},
				map[string][]string{
					"Location": {"bar", "anotherbar"},
				},
			},
			requests: struct {
				rewrite       []Rewrite
				reqHeader     http.Header
				expRespHeader http.Header
			}{[]Rewrite{
				{
					Header:      "Location",
					Regex:       "^http://(.+)$",
					Replacement: "https://$1",
				},
			},
				map[string][]string{
					"Location": {"http://test:1000"},
				},
				map[string][]string{
					"Location": {"https://test:1000"},
				},
			},
		},
		{
			desc: "should replace http by https in location response header",
			responses: struct {
				rewrite       []Rewrite
				reqHeader     http.Header
				expRespHeader http.Header
			}{[]Rewrite{
				{
					Header:      "Location",
					Regex:       "^http://(.+)$",
					Replacement: "https://$1",
				},
			},
				map[string][]string{
					"Location": {"http://test:1000"},
				},
				map[string][]string{
					"Location": {"https://test:1000"},
				},
			},
		},
		{
			desc: "should replace http by https in location response header",
			responses: struct {
				rewrite       []Rewrite
				reqHeader     http.Header
				expRespHeader http.Header
			}{[]Rewrite{
				{
					Header:      "Location",
					Regex:       "^http://(.+)$",
					Replacement: "https://$1",
				},
			},
				map[string][]string{
					"Location": {"http://test:1000"},
				},
				map[string][]string{
					"Location": {"https://test:1000"},
				},
			},
		},
		{
			desc: "should replace http by https in location response header",
			requests: struct {
				rewrite       []Rewrite
				reqHeader     http.Header
				expRespHeader http.Header
			}{[]Rewrite{
				{
					Header:      "Location",
					Regex:       "^http://(.+)$",
					Replacement: "https://$1",
				},
			},
				map[string][]string{
					"Location": {"http://test:1000"},
				},
				map[string][]string{
					"Location": {"https://test:1000"},
				},
			},
		},
		{
			desc: "should reorder Content-Type header in request header",
			requests: struct {
				rewrite       []Rewrite
				reqHeader     http.Header
				expRespHeader http.Header
			}{[]Rewrite{
				{
					Header:      "Content-Type",
					Regex:       `(.*)(charset=.+){1};\s(boundary=.+\"){1}(.*)`,
					Replacement: "$1$3; $2 $4",
				},
			},
				map[string][]string{
					"Content-Type": {`multipart/form-data; charset=utf-8; boundary="__boundary__"`},
				},
				map[string][]string{
					"Content-Type": {`multipart/form-data; boundary="__boundary__"; charset=utf-8`},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			config := &Config{
				Rewrites: Rewrites{
					Request:  test.requests.rewrite,
					Response: test.responses.rewrite,
				},
			}

			// Response setup if any
			next := func(rw http.ResponseWriter, req *http.Request) {
				for k, v := range test.responses.reqHeader {
					for _, h := range v {
						rw.Header().Add(k, h)
					}
				}
				rw.WriteHeader(http.StatusOK)
			}

			rewriteBody, err := New(context.Background(), http.HandlerFunc(next), config, "rewriteHeader")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// Request setup if any
			for k, v := range test.requests.reqHeader {
				for _, h := range v {
					req.Header.Add(k, h)
				}
			}

			rewriteBody.ServeHTTP(recorder, req)

			// Validate response headers
			for k, expected := range test.responses.expRespHeader {
				values := recorder.Header().Values(k)
				if !testEq(values, expected) {
					t.Errorf("Slice arent equals: expect: %+v, result: %+v", expected, values)
				}
			}

			// Validate request headers
			for k, expected := range test.requests.expRespHeader {
				values := req.Header.Values(k)

				if !testEq(values, expected) {
					t.Errorf("Slice arent equals: expect: %+v, result: %+v", expected, values)
				}
			}
		})
	}
}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
