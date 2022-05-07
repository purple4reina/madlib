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
		code  int
	}{
		{
			name:  "works fine",
			words: &words{Noun: "cat", Verb: "run", Adjective: "hot"},
			err:   nil,
			code:  200,
		},
		{
			name:  "empty responses",
			words: &words{},
			err:   nil,
			code:  200,
		},
		{
			name:  "error response",
			words: &words{Noun: "cat", Verb: "run", Adjective: "hot"},
			err:   fmt.Errorf("oops"),
			code:  500,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			client = newTestClient(test.words, test.err)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			var expect string
			if test.err == nil {
				expect = fmt.Sprintf(`{"madlib":"It was a %s day. I went `+
					`downstairs to see if I could %s dinner. I asked, \"Does `+
					`the stew need fresh %s?\""}`, test.words.Adjective,
					test.words.Verb, test.words.Noun)
				if resp.Body.String() != expect {
					t.Errorf("incorrect response body\nactual=%s\nexpect=%s",
						resp.Body, expect)
				}
			} else {
				// Since which of the three API calls responds first is
				// nondeterministic, test for all three
				expNoun := fmt.Sprintf(`{"error":"Get \"%snoun\": %s"}`,
					wordURL, test.err.Error())
				expVerb := fmt.Sprintf(`{"error":"Get \"%sverb\": %s"}`,
					wordURL, test.err.Error())
				expAdj := fmt.Sprintf(`{"error":"Get \"%sadjective\": %s"}`,
					wordURL, test.err.Error())
				if r := resp.Body.String(); r != expNoun && r != expVerb && r != expAdj {
					t.Errorf("incorrect response body\nactual=%s\nexpect=%s",
						resp.Body, expect)
				}
			}

			if resp.Code != test.code {
				t.Errorf("incorrect response code actual=%d expect=%d",
					resp.Code, test.code)
			}
		})
	}
}
