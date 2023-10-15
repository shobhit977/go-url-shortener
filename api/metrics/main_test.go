package main

import (
	"go-url-shortener/lib/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getMostShortenedUrl(t *testing.T) {
	type input struct {
		urlDetails []models.UrlInfo
		limit      int
	}
	tests := []struct {
		name  string
		input input
		resp  []models.MetricsResponse
	}{
		{
			name: "success case - limit 3",
			input: input{
				urlDetails: []models.UrlInfo{
					{Url: "www.google.com", ShortUrl: "xyz3", Domain: "google"},
					{Url: "www.reddit.com", ShortUrl: "xyz", Domain: "reddit"},
					{Url: "www.google.com/abc", ShortUrl: "xyza", Domain: "google"},
					{Url: "www.facebook.com/abc", ShortUrl: "xyzs", Domain: "facebook"},
					{Url: "www.facebook.com", ShortUrl: "xyz1", Domain: "facebook"},
					{Url: "www.google.com/abcd", ShortUrl: "xyzxa", Domain: "google"},
				},
				limit: 3,
			},
			resp: []models.MetricsResponse{
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
				urlDetails: []models.UrlInfo{
					{Url: "www.google.com", ShortUrl: "xyz3", Domain: "google"},
					{Url: "www.reddit.com", ShortUrl: "xyz", Domain: "reddit"},
					{Url: "www.google.com/abc", ShortUrl: "xyza", Domain: "google"},
					{Url: "www.facebook.com/abc", ShortUrl: "xyzs", Domain: "facebook"},
					{Url: "www.facebook.com", ShortUrl: "xyz1", Domain: "facebook"},
					{Url: "www.google.com/abcd", ShortUrl: "xyzxa", Domain: "google"},
				},
				limit: 2,
			},
			resp: []models.MetricsResponse{
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
				urlDetails: []models.UrlInfo{
					{Url: "www.google.com", ShortUrl: "xyz3", Domain: "google"},
					{Url: "www.reddit.com", ShortUrl: "xyz", Domain: "reddit"},
					{Url: "www.google.com/abc", ShortUrl: "xyza", Domain: "google"},
					{Url: "www.facebook.com/abc", ShortUrl: "xyzs", Domain: "facebook"},
					{Url: "www.facebook.com", ShortUrl: "xyz1", Domain: "facebook"},
					{Url: "www.google.com/abcd", ShortUrl: "xyzxa", Domain: "google"},
				},
				limit: 1,
			},
			resp: []models.MetricsResponse{
				{
					Domain: "google",
					Count:  3,
				},
			},
		},
		{
			name: "failure case - empty url data",
			input: input{
				urlDetails: []models.UrlInfo{},
				limit:      2,
			},
			resp: []models.MetricsResponse(nil),
		},
	}
	for _, tt := range tests {
		resp := getMostShortenedUrls(tt.input.urlDetails, tt.input.limit)
		assert.Equal(t, resp, tt.resp)
	}
}
