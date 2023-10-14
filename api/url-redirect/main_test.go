package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsURLExist(t *testing.T) {
	type input struct {
		urlDetails []UrlInfo
		shortUrl   string
	}
	tests := []struct {
		name  string
		input input
		url   string
		exist bool
	}{
		{
			name: "success case",
			input: input{
				urlDetails: []UrlInfo{{Url: "www.google.com", ShortUrl: "xyz", Domain: "google"}},
				shortUrl:   "xyz",
			},
			url:   "www.google.com",
			exist: true,
		},
		{
			name: "failure case - url does not exist",
			input: input{
				urlDetails: []UrlInfo{{Url: "www.google.com", ShortUrl: "xyz", Domain: "google"}},
				shortUrl:   "abc",
			},
			url:   "",
			exist: false,
		},
	}
	for _, tt := range tests {
		url, exist := isUrlExist(tt.input.urlDetails, tt.input.shortUrl)
		assert.Equal(t, url, tt.url)
		assert.Equal(t, exist, tt.exist)
	}
}
