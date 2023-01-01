package main

import (
	"encoding"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTotpToksUnmarshalText(t *testing.T) {
	type inOut struct {
		in        []byte
		out       TotpToks
		expectErr bool
	}
	testCases := map[string]inOut{
		"one input": {
			in: []byte(`hoge:xxxxxxxxxxxx`),
			out: TotpToks{map[service]totpTok{
				"hoge": `xxxxxxxxxxxx`,
			}},
			expectErr: false,
		},
		"multi input": {
			in: []byte(`hoge:xxxxx,fuga:111111`),
			out: TotpToks{map[service]totpTok{
				"hoge": "xxxxx",
				"fuga": "111111",
			}},
			expectErr: false,
		},
		"no token with service name": {
			in:        []byte(`hoge`),
			out:       TotpToks{nil},
			expectErr: true,
		},
		"no token in last service": {
			in:        []byte(`hoge:xxxxxxx,fuga`),
			out:       TotpToks{nil},
			expectErr: true,
		},
		"empty input": {
			in:        []byte(``),
			out:       TotpToks{m: nil},
			expectErr: true,
		},
	}
	for k, v := range testCases {
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			var buf TotpToks
			err := buf.UnmarshalText(v.in)
			if v.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, v.out, buf)
		})
	}
}

func TestTotpToksImpl(t *testing.T) {
	assert.Implements(t, (*encoding.TextUnmarshaler)(nil), new(TotpToks))
}
