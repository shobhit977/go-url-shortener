package main

import (
	"encoding/json"
	"go-url-shortener/lib/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getDomain_SuccessCases(t *testing.T) {
	tests := []struct {
		name     string
		inputUrl string
		domain   string
		err      error
	}{
		{
			name:     "success case",
			inputUrl: "https://www.google.com",
			domain:   "google",
			err:      nil,
		},
		{
			name:     "success case",
			inputUrl: "https://www.facebook.com",
			domain:   "facebook",
			err:      nil,
		},
		{
			name:     "success case",
			inputUrl: "https://facebook.com",
			domain:   "facebook",
			err:      nil,
		},
	}
	for _, tt := range tests {
		domain, err := getDomain(tt.inputUrl)
		assert.Equal(t, domain, tt.domain)
		assert.Nil(t, err)
	}
}

func Test_getDomain_FailureCases(t *testing.T) {
	tests := []struct {
		name     string
		inputUrl string
		domain   string
		err      error
	}{
		{
			name:     "failure case - malformed url",
			inputUrl: "postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require",
			domain:   "",
		},
		{
			name:     "failure case - unsupported character",
			inputUrl: "https://www.percent-off.com/_20_%+off_60000_",
			domain:   "",
		},
	}
	for _, tt := range tests {
		domain, err := getDomain(tt.inputUrl)
		assert.Equal(t, domain, tt.domain)
		assert.NotNil(t, err)
	}
}

func Test_generateUrlFileOutput_getDomain_Failure(t *testing.T) {
	testUrlData := []models.UrlInfo{}
	invalidUrl := "https://www.percent-off.com/_20_%+off_60000_"
	allUrlInfoBytes, urlInfo, err := generateUrlFileOutput(testUrlData, invalidUrl)
	assert.Nil(t, allUrlInfoBytes)
	assert.Empty(t, urlInfo)
	assert.NotNil(t, err)
	assert.Error(t, err, "invalid URL. Please specify a valid URL")
}

func Test_generateUrlFileOutput_Success_EmptyInitialData(t *testing.T) {
	testUrlData := []models.UrlInfo{}
	url := "https://www.google.com"
	_, urlInfo, err := generateUrlFileOutput(testUrlData, url)
	assert.Nil(t, err)
	assert.Equal(t, urlInfo, models.UrlInfo{
		Url:      "https://www.google.com",
		Domain:   "google",
		ShortUrl: "7378mDnD7g",
	})
}

func Test_generateUrlFileOutput_Success_ExistingUrlDataPresent(t *testing.T) {
	testUrlData := []models.UrlInfo{{
		Url:      "www.fb.com",
		ShortUrl: "Vb1V-0nHFy",
		Domain:   "fb",
	}}
	url := "https://www.google.com"
	allUrlInfoBytes, urlInfo, err := generateUrlFileOutput(testUrlData, url)
	assert.Nil(t, err)
	assert.Equal(t, urlInfo, models.UrlInfo{
		Url:      "https://www.google.com",
		Domain:   "google",
		ShortUrl: "7378mDnD7g",
	},
	)
	var allUrlData []models.UrlInfo
	json.Unmarshal(allUrlInfoBytes, &allUrlData)
	assert.Equal(t, allUrlData, []models.UrlInfo{{
		Url:      "www.fb.com",
		ShortUrl: "Vb1V-0nHFy",
		Domain:   "fb",
	},
		{
			Url:      "https://www.google.com",
			Domain:   "google",
			ShortUrl: "7378mDnD7g",
		}})
}

func Test_generateShortUrlLength(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "google",
			url:  "www.google.com",
		},
		{
			name: "facebook",
			url:  "www.facebook.com",
		},
		{
			name: "reddit",
			url:  "www.reddit.com",
		},
	}
	for _, tt := range tests {
		shortUrl := generateShortUrl(tt.url)
		assert.Equal(t, len(shortUrl), 10)
	}
}
