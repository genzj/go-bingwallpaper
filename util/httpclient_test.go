package util

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"github.com/genzj/gobingwallpaper/i18n"
)

type DummyData struct {
	Images []struct {
		URL string
	}
	ToolTips struct {
		Loading  string
		Previous string
		Next     string
	}
}

func TestHttpGetJSON(t *testing.T) {
	var data DummyData
	err := HTTPGetJSON(`http://bing.com/HPImageArchive.aspx?format=js&video=1&n=10&idx=1&mkt=ja-jp`, &data)

	if err != nil {
		t.Fatalf("error %v happened", err)
	}

	if data.ToolTips.Loading == "" || data.ToolTips.Next == "" || data.ToolTips.Previous == "" {
		t.Fatalf("response json has empty fields %#v", data)
	}
}

func TestHttpGetJSONIncorrectType(t *testing.T) {
	var data DummyData
	err := HTTPGetJSON(`http://bing.com/HPImageArchive.aspx?format=xml&video=1&n=10&idx=1&mkt=ja-jp`, &data)

	if err == nil {
		t.Fatalf("expected error not happened")
	}

	if !strings.Contains(err.Error(), "http_content_type_error_url") {
		t.Fatalf("unexpected error %v happened instead of content-type error", err)
	}
}

func TestHttpSetProxy(t *testing.T) {
	// TODO need mock of socks5
	SetProxy("socks5", "socks5://127.0.0.1:1080")
	TestHttpGetJSON(t)
}

func TestHttpSetProxyCreateSingleton(t *testing.T) {
	client1 := getClient()
	SetProxy("socks5", "socks5://127.0.0.1:1080")
	client2 := getClient()
	if client1 == client2 {
		t.Fatal("no new client created after setting proxy")
	}
}

func TestMain(m *testing.M) {
	var responder httpmock.Responder
	var jsonResponse interface{}
	var err error

	i18n.LoadMockTFunc()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	err = json.Unmarshal([]byte(`
{
    "images": [
        {
            "startdate": "20161224",
            "fullstartdate": "201612241500",
            "enddate": "20161225",
            "url": "/az/hprichbg/rb/SnowGlobeVideo_JA-JP8461656803_1920x1080.jpg",
            "urlbase": "/az/hprichbg/rb/SnowGlobeVideo_JA-JP8461656803",
            "copyright": "｢トレントのクリスマス・マーケット｣イタリア,　南ティロル自治州 (© Carlo Trolese/500px)",
            "copyrightlink": "http://www.bing.com/search?q=%E3%83%88%E3%83%AC%E3%83%B3%E3%83%88%E3%80%80%E3%82%A4%E3%82%BF%E3%83%AA%E3%82%A2&form=hpcapt&filters=HpDate:%2220161224_1500%22",
            "quiz": "/search?q=Bing+homepage+quiz&filters=WQOskey:%22HPQuiz_20161224_SnowGlobeVideo%22&FORM=HPQUIZ",
            "wp": false,
            "hsh": "66e96180b6671f4cdd9ec6ff654bfa14",
            "drk": 1,
            "top": 1,
            "bot": 1,
            "hs": [],
            "vid": {
                "sources": [
                    [
                        "mp4",
                        "video/mp4; codecs=\"avc1.42E01E, mp4a.40.2\"",
                        "//az29176.vo.msecnd.net/videocontent/Xmas_Market_Globes_500px_RF_69446295_768_JA-JP.mp4"
                    ],
                    [
                        "mp4hd",
                        "video/mp4; codecs=\"avc1.42E01E, mp4a.40.2\"",
                        "//az29176.vo.msecnd.net/videocontent/Xmas_Market_Globes_500px_RF_69446295_1080_HD_JA-JP.mp4"
                    ],
                    [
                        "mp4mobile",
                        "video/mp4; codecs=\"avc1.42E01E, mp4a.40.2\"",
                        "//az29176.vo.msecnd.net/videocontent/Xmas_Market_Globes_500px_RF_mobile_360_JA-JP.mp4"
                    ]
                ],
                "loop": true,
                "image": "//az29176.vo.msecnd.net/videocontent/Xmas_Market_Globes_FF_768_HD_JA-JP1853590220.jpg",
                "caption": "｢トレントのクリスマス・マーケット｣イタリア,　南ティロル自治州 (© Carlo Trolese/500px)",
                "captionlink": "",
                "dark": true
            }
        }
    ],
    "tooltips": {
        "loading": "正在加载...",
        "previous": "上一个图像",
        "next": "下一个图像",
        "walle": "此图片不能下载用作壁纸。",
        "walls": "下载今日美图。仅限用作桌面壁纸。",
        "play": "播放视频",
        "pause": "暂停视频"
    }
}
	`), &jsonResponse)
	if err == nil {
		responder, err = httpmock.NewJsonResponder(200, jsonResponse)
	}

	httpmock.RegisterResponder("GET", "http://bing.com/HPImageArchive.aspx?format=js&video=1&n=10&idx=1&mkt=ja-jp", responder)
	if err != nil {
		panic(fmt.Sprintf("load mock json response failed: %v", err))
	}

	// test content-type awareness, nil result is enough
	responder, err = httpmock.NewXmlResponder(200, nil)
	httpmock.RegisterResponder("GET", "http://bing.com/HPImageArchive.aspx?format=xml&video=1&n=10&idx=1&mkt=ja-jp", responder)
	if err != nil {
		panic(fmt.Sprintf("load mock xml response failed: %v", err))
	}

	newClient()
	clientSingleton.Transport = httpmock.DefaultTransport

	flag.Parse()
	os.Exit(m.Run())
}
