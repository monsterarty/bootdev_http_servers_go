package main

import (
	"testing"
)

func TestProfaneFilter(t *testing.T) {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "replaces single profane word",
			s:    "This is a kerfuffle opinion I need to share with the world",
			want: "This is a **** opinion I need to share with the world",
		},
		{
			name: "leaves clean text unchanged",
			s:    "I had something interesting for breakfast",
			want: "I had something interesting for breakfast",
		},
		{
			name: "replaces single profane word in second sentence",
			s:    "I hear Mastodon is better than Chirpy. sharbert I need to migrate",
			want: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			name: "handles empty string",
			s:    "",
			want: "",
		},
		{
			name: "replaces multiple profane words case insensitive",
			s:    "KERFUFFLE , SHARbert and fornaX",
			want: "**** , **** and ****",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := profaneFilter(tt.s, badWords)

			if got != tt.want {
				t.Errorf(
					"profaneFilter(%q, %q) = %q, want = %q",
					tt.s,
					badWords,
					got,
					tt.want,
				)
			}
		})
	}
}
