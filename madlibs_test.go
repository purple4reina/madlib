package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRoundTripper func(*http.Request) (*http.Response, error)

func (rt testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

func newTestClient(w *words, err error) *http.Client {
	return &http.Client{
		Transport: testRoundTripper(func(req *http.Request) (*http.Response, error) {
			var word string
			switch req.URL.Path {
			case "/noun":
				word = w.Noun
			case "/verb":
				word = w.Verb
			case "/adjective":
				word = w.Adjective
			default:
				panic("unexpected path found")
			}
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader([]byte(word))),
			}, err
		}),
	}
}

var req = func() *http.Request {
	req, err := http.NewRequest("GET", "/madlib", nil)
	if err != nil {
		panic(err)
	}
	return req
}()

func TestMadlibEndpoint(t *testing.T) {
	r := newRouter()
	testcases := []struct {
		name  string
		words *words
		err   error
	}{
		{
			name:  "works fine",
			words: &words{Noun: "cat", Verb: "run", Adjective: "hot"},
			err:   nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			client = newTestClient(test.words, test.err)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			expect := fmt.Sprintf(`{"madlib":"It was a %s day. I went `+
				`downstairs to see if I could %s dinner. I asked, \"Does `+
				`the stew need fresh %s?\""}`, test.words.Adjective,
				test.words.Verb, test.words.Noun)
			if resp.Body.String() != expect {
				t.Errorf("incorrect response body\nactual=%s\nexpect=%s", resp.Body, expect)
			}
		})
	}
}
