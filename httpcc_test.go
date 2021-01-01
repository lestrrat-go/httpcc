package httpcc_test

import (
	"testing"

	httpcc "github.com/lestrrat-go/httpcc"
	"github.com/stretchr/testify/assert"
)

func compareRequestDirective(t *testing.T, expected map[string]interface{}, dir *httpcc.RequestDirective) bool {
	t.Helper()

	for k, v := range expected {
		switch k {
		case httpcc.MaxAge:
			got, ok := dir.MaxAge()
			if !assert.True(t, ok, `dir.MaxAge() should return true`) {
				return false
			}

			if !assert.Equal(t, v, got) {
				return false
			}
		case httpcc.NoStore:
			got := dir.NoStore()
			if !assert.Equal(t, v, got) {
				return false
			}
		case "extensions":
			got := dir.Extensions()
			if !assert.Equal(t, v, got) {
				return false
			}
		default:
			assert.Fail(t, `unhandled field %s`, k)
			return false
		}
	}
	return true
}

func compareResponseDirective(t *testing.T, expected map[string]interface{}, dir *httpcc.ResponseDirective) bool {
	t.Helper()

	for k, v := range expected {
		switch k {
		case httpcc.MaxAge:
			got, ok := dir.MaxAge()
			if !assert.True(t, ok, `dir.MaxAge() should return true`) {
				return false
			}

			if !assert.Equal(t, v, got) {
				return false
			}
		case httpcc.NoStore:
			got := dir.NoStore()
			if !assert.Equal(t, v, got) {
				return false
			}
		case "extensions":
			got := dir.Extensions()
			if !assert.Equal(t, v, got) {
				return false
			}
		default:
			assert.Fail(t, `unhandled field %s`, k)
			return false
		}
	}
	return true
}

func TestParseDirective(t *testing.T) {
	testcases := []struct {
		Source    string
		Error     bool
		Expected  *httpcc.TokenPair
		IsRequest bool
	}{
		{
			Source: `no-store="foo"`,
			Error:  true,
		},
		{
			Source:   `s-maxage=4649`,
			Expected: &httpcc.TokenPair{Name: `s-maxage`, Value: `4649`},
		},
		{
			Source:    `s-maxage=4649`,
			Expected:  &httpcc.TokenPair{Name: `s-maxage`, Value: `4649`},
			IsRequest: true,
		},
		{
			Source:   "no-store",
			Expected: &httpcc.TokenPair{Name: "no-store"},
		},
		{
			Source:    "max-age=4649",
			Expected:  &httpcc.TokenPair{Name: "max-age", Value: "4649"},
			IsRequest: true,
		},
		{
			Source:    `max-age="4649"`,
			Error:     true,
			IsRequest: true,
		},
		{
			Source:    `max-age="4649`,
			Error:     true,
			IsRequest: true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Source, func(t *testing.T) {
			var pair *httpcc.TokenPair
			var err error
			if tc.IsRequest {
				pair, err = httpcc.ParseRequestDirective(tc.Source)
			} else {
				pair, err = httpcc.ParseResponseDirective(tc.Source)
			}
			if tc.Error {
				if !assert.Error(t, err, `expected to return an error`) {
					return
				}
			} else {
				if !assert.NoError(t, err, `expected to succeed`) {
					return
				}
				if !assert.Equal(t, tc.Expected, pair, `expected to return pair`) {
					return
				}
			}
		})
	}
}

func TestParseDirectives(t *testing.T) {
	testcases := []struct {
		Source    string
		Error     bool
		Expected  []*httpcc.TokenPair
		IsRequest bool
	}{
		{
			Source:    `max-age=4649, no-store`,
			IsRequest: true,
			Expected: []*httpcc.TokenPair{
				{Name: `max-age`, Value: `4649`},
				{Name: `no-store`},
			},
		},
		{
			Source:    `   max-age=4649    , no-store     `,
			IsRequest: true,
			Expected: []*httpcc.TokenPair{
				{Name: `max-age`, Value: `4649`},
				{Name: `no-store`},
			},
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Source, func(t *testing.T) {
			var tokens []*httpcc.TokenPair
			var err error
			if tc.IsRequest {
				tokens, err = httpcc.ParseRequestDirectives(tc.Source)
			} else {
				tokens, err = httpcc.ParseResponseDirectives(tc.Source)
			}
			if tc.Error {
				if !assert.Error(t, err, `expected to return an error`) {
					return
				}
			} else {
				if !assert.NoError(t, err, `expected to succeed`) {
					return
				}
				if !assert.Equal(t, tc.Expected, tokens, `expected to return list of tokens`) {
					return
				}
			}
		})
	}
}

func TestParseRequest(t *testing.T) {
	testcases := []struct {
		Source   string
		Error    bool
		Expected map[string]interface{}
	}{
		{
			Source: `max-age=4649, no-store`,
			Expected: map[string]interface{}{
				httpcc.MaxAge:  uint64(4649),
				httpcc.NoStore: true,
			},
		},
		{
			Source: `max-age="4649"`,
			Error:  true,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Source, func(t *testing.T) {
			dir, err := httpcc.ParseRequest(tc.Source)
			if tc.Error {
				if !assert.Error(t, err, `expected to return an error`) {
					return
				}
			} else {
				if !assert.NoError(t, err, `expected to succeed`) {
					return
				}
				if !compareRequestDirective(t, tc.Expected, dir) {
					return
				}
			}
		})
	}
}

func TestParseResponse(t *testing.T) {
	testcases := []struct {
		Source   string
		Error    bool
		Expected map[string]interface{}
	}{
		{
			Source: `max-age=4649, no-store, community="UCI"`,
			Expected: map[string]interface{}{
				httpcc.MaxAge:  uint64(4649),
				httpcc.NoStore: true,
				"extensions": map[string]string{
					"community": "UCI",
				},
			},
		},
		{
			Source: `max-age="4649"`,
			Error:  true,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Source, func(t *testing.T) {
			dir, err := httpcc.ParseResponse(tc.Source)
			if tc.Error {
				if !assert.Error(t, err, `expected to return an error`) {
					return
				}
			} else {
				if !assert.NoError(t, err, `expected to succeed`) {
					return
				}
				if !compareResponseDirective(t, tc.Expected, dir) {
					return
				}
			}
		})
	}
}
