package util

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/genzj/gobingwallpaper/i18n"
)

type DummyData struct {
	ToolTips struct {
		Loading  string
		Previous string
		Next     string
	}
}

func TestHttpGetJSON(t *testing.T) {
	var data DummyData
	err := HttpGetJSON(`http://bing.com/HPImageArchive.aspx?format=js&video=1&n=10&idx=1&mkt=ja-jp`, &data)

	if err != nil {
		t.Fatalf("error %v happened", err)
	}

	if data.ToolTips.Loading == "" || data.ToolTips.Next == "" || data.ToolTips.Previous == "" {
		t.Fatalf("response json has empty fields %#v", data)
	}
}

func TestHttpGetJSONIncorrectType(t *testing.T) {
	var data DummyData
	err := HttpGetJSON(`http://bing.com/HPImageArchive.aspx?format=xml&video=1&n=10&idx=1&mkt=ja-jp`, &data)

	if err == nil {
		t.Fatalf("expected error not happened")
	}

	if !strings.Contains(err.Error(), "returns unexpected content-type") {
		t.Fatalf("unexpected error %v happened instead of content-type error", err)
	}
}

func TestHttpSetProxy(t *testing.T) {
	SetProxy("socks5", "socks5://127.0.0.1:1080")
	TestHttpGetJSON(t)
}

func TestMain(m *testing.M) {
	i18n.SetLanguageFilePath(".")
	flag.Parse()
	os.Exit(m.Run())
}
