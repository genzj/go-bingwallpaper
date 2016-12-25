package bingwallpaper

import (
	"github.com/genzj/gobingwallpaper/log"
	"net/url"
	"strconv"
)

type Provider interface {
	GetBaseURL() *url.URL
	GetParameters() map[string]string
	FormURL() *url.URL
}

var (
	DefaultBaseProvider = BaseProvider{
		BaseURL: "http://www.bing.com/HPImageArchive.aspx",
	}
)

type BaseProvider struct {
	// BaseURL specifies the bing.com wallpaper meta data serving endpoint
	BaseURL string

	// ItemCount specifies number of (historical) pictures to be listed
	// zero value stands for default count which is 10 because n=0 leads to
	// meaningless null
	ItemCount int

	// Index specifies offset of pictures, 0 is today's pic, -1 is tomorrows, 1
	// is yesterday, etc.
	Index int

	// Video is non-zero if video shall be retrieved or vice versa
	Video int

	// Market is the country pictures will be displayed in. It should be locale
	// code like "en-us". For full list of bing markets, login bing.com then
	// visit the Country/Region section of http://www.bing.com/account/general
	// Zero value omits this parameter from query url and server will normally
	// decides the market by client IP
	Market string
}

func (p BaseProvider) GetBaseURL() *url.URL {
	u, err := url.Parse(p.BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Base URL parsed: %v from %v", u.String(), p.BaseURL)
	return u
}

func (b BaseProvider) GetParameters() map[string]string {
	var n int
	if n = b.ItemCount; n <= 0 {
		n = 10
	}
	ans := map[string]string{
		"format": "js",
		"idx":    strconv.Itoa(b.Index),
		"mbl":    "1",
		"n":      strconv.Itoa(n),
		"video":  strconv.Itoa(b.Video),
	}
	if market := b.Market; market != "" {
		ans["market"] = market
	}
	log.Debugf("Parameters loaded: %v", ans)
	return ans
}

func (p BaseProvider) FormURL() *url.URL {
	u := p.GetBaseURL()
	q := u.Query()
	for k, v := range p.GetParameters() {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u
}
