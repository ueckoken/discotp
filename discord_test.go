package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTotpGenGenerateCode(t *testing.T) {
	type testCase struct {
		secret    string
		tokenCode string
		wantErr   bool
	}
	testCases := map[string]testCase{
		"with no white": {
			secret:    "xxxxxxxxxxxxxxxxxxxx",
			tokenCode: `000241`,
			wantErr:   false,
		},
		"with internal white": {
			secret:    "xxxxx xxxxx xxxxx xxxxx",
			tokenCode: `000241`,
			wantErr:   false,
		},
		"with internal white and leading, trailing spaces": {
			secret:    " xxxxx xxxxx xxxxx xxxxx ",
			tokenCode: `000241`,
			wantErr:   false,
		},
	}
	for k, v := range testCases {
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			tokGen := TotpGen{secret: v.secret}
			code, err := tokGen.GenerateCode(time.Time{})
			if v.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, v.tokenCode, code)
		})
	}
}
