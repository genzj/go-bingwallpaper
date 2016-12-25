package bingwallpaper

import (
	"testing"
)

func TestBaseProviderGetBase(t *testing.T) {
	p := BaseProvider{
		BaseURL: "https://bing.com/HPImageArchive.aspx",
	}
	u := p.GetBaseURL()
	if u == nil {
		t.Fatal("nothing returned")
	}
	if u.String() != "https://bing.com/HPImageArchive.aspx" {
		t.Fatal("got ", u.String(), " but expected ", p.BaseURL)
	}
}

func TestBaseProviderGetParameters(t *testing.T) {
	var exists bool
	var value string

	p := BaseProvider{
		Index:  -1,
		Video:  1,
		Market: "ja-jp",
	}
	u := p.GetParameters()
	if u == nil {
		t.Fatal("nothing returned")
	}
	if value, exists = u["format"]; !exists || value != "js" {
		t.Fatal("got format=", value, " but expected js")
	}
	if value, exists = u["market"]; !exists || value != "ja-jp" {
		t.Fatal("got market=", value, " but expected ja-jp")
	}
	if value, exists = u["n"]; !exists || value != "10" {
		t.Fatal("got n=", value, " but expected 10")
	}
	if value, exists = u["idx"]; !exists || value != "-1" {
		t.Fatal("got idx=", value, " but expected -1")
	}
	if value, exists = u["video"]; !exists || value != "1" {
		t.Fatal("got video=", value, " but expected 1")
	}
}

func TestBaseProviderFormURL(t *testing.T) {
	const expected = "https://bing.com/HPImageArchive.aspx?format=js&idx=-1&market=ja-jp&mbl=1&n=10&video=1"
	p := BaseProvider{
		BaseURL: "https://bing.com/HPImageArchive.aspx",
		Market:  "ja-jp",
		Index:   -1,
		Video:   1,
	}
	u := p.FormURL()
	if u == nil {
		t.Fatal("nothing returned")
	}
	if u.String() != expected {
		t.Fatal("got url ", u, " but expected ", expected)
	}
}
