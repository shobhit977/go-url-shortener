package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getMostShortenedUrl(t *testing.T) {
	type input struct {
		urlDetails []UrlInfo
		limit      int
	}
	tests := []struct {
		name  string
		input input
		resp  []Response
	}{
		{
			name: "success case - limit 3",
			input: input{
				urlDetails: []UrlInfo{
					{Url: "www.google.com", ShortUrl: "xyz3", Domain: "google"},
					{Url: "www.reddit.com", ShortUrl: "xyz", Domain: "reddit"},
					{Url: "www.google.com/abc", ShortUrl: "xyza", Domain: "google"},
					{Url: "www.facebook.com/abc", ShortUrl: "xyzs", Domain: "facebook"},
					{Url: "www.facebook.com", ShortUrl: "xyz1", Domain: "facebook"},
					{Url: "www.google.com/abcd", ShortUrl: "xyzxa", Domain: "google"},
				},
				limit: 3,
			},
			resp: []Response{
				{
					Domain: "google",
					Count:  3,
				},
				{
					Domain: "facebook",
					Count:  2,
				},
				{
					Domain: "reddit",
					Count:  1,
				},
			},
		},
		{
			name: "success case - limit 2",
			input: input{
				urlDetails: []UrlInfo{
					{Url: "www.google.com", ShortUrl: "xyz3", Domain: "google"},
					{Url: "www.reddit.com", ShortUrl: "xyz", Domain: "reddit"},
					{Url: "www.google.com/abc", ShortUrl: "xyza", Domain: "google"},
					{Url: "www.facebook.com/abc", ShortUrl: "xyzs", Domain: "facebook"},
					{Url: "www.facebook.com", ShortUrl: "xyz1", Domain: "facebook"},
					{Url: "www.google.com/abcd", ShortUrl: "xyzxa", Domain: "google"},
				},
				limit: 2,
			},
			resp: []Response{
				{
					Domain: "google",
					Count:  3,
				},
				{
					Domain: "facebook",
					Count:  2,
				},
			},
		},
		{
			name: "success case - limit 2",
			input: input{
				urlDetails: []UrlInfo{
					{Url: "www.google.com", ShortUrl: "xyz3", Domain: "google"},
					{Url: "www.reddit.com", ShortUrl: "xyz", Domain: "reddit"},
					{Url: "www.google.com/abc", ShortUrl: "xyza", Domain: "google"},
					{Url: "www.facebook.com/abc", ShortUrl: "xyzs", Domain: "facebook"},
					{Url: "www.facebook.com", ShortUrl: "xyz1", Domain: "facebook"},
					{Url: "www.google.com/abcd", ShortUrl: "xyzxa", Domain: "google"},
				},
				limit: 1,
			},
			resp: []Response{
				{
					Domain: "google",
					Count:  3,
				},
			},
		},
		{
			name: "failure case - empty url data",
			input: input{
				urlDetails: []UrlInfo{},
				limit:      2,
			},
			resp: []Response(nil),
		},
	}
	for _, tt := range tests {
		resp := getMostShortenedUrls(tt.input.urlDetails, tt.input.limit)
		assert.Equal(t, resp, tt.resp)
	}
}
