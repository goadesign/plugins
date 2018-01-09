package cors

import (
	"testing"
)

func TestMatchOrigin(t *testing.T) {
	cases := map[string]map[string][]struct {
		origin string
		output bool
	}{
		"string-spec": {
			"domain": {
				{"domain", true},
				{"some-other-domain", false},
			},
		},
		"wildcard-spec": {
			"*": {
				{"any domain", true},
				{"/any regex domain/", true},
			},
			"some*domain": {
				{"some.other.domain", true},
				{"not.this.domain", false},
				{"some.domain.com", false},
			},
		},
		"regex-spec": {
			"/.*domain\\..+/": {
				{"some.domain.com", true},
				{"not.this.domain", false},
			},
		},
	}
	for name, tests := range cases {
		t.Run(name, func(t *testing.T) {
			for spec, tcs := range tests {
				for _, tc := range tcs {
					output := MatchOrigin(tc.origin, spec)
					if output != tc.output {
						t.Errorf("TestMatchOrigin(%s): MatchOrigin(%q, %q): Expected %t, Got %t", name, tc.origin, spec, tc.output, output)
					}
				}
			}
		})
	}
}
