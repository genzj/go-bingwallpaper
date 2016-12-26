package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/genzj/gobingwallpaper/i18n"
	"github.com/genzj/gobingwallpaper/log"
	"golang.org/x/net/proxy"
)

var client *http.Client

var proxyConf struct {
	proxyType string
	proxyURL  string
}

func newClient() {
	transport := http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	nextClient := http.Client{
		Transport: &transport,
	}

	switch strings.ToLower(proxyConf.proxyType) {
	case "socks5":
		transport.DialContext = nil
		proxyURL, err := url.Parse(proxyConf.proxyURL)
		if err != nil {
			log.Errorf(i18n.T("proxy_url_parse_error", i18n.Fields{
				"URL":   proxyConf.proxyURL,
				"Error": err,
			}))
			return
		}

		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			log.Errorf(i18n.T("proxy_url_parse_error", i18n.Fields{
				"URL":   proxyURL,
				"Error": err,
			}))
			return
		} else {
			log.Debugf("SOCKS5 proxy %T installed from URL %v", dialer, proxyURL)
		}

		transport.Dial = dialer.Dial
		transport.Proxy = nil
	default:
		transport.Proxy = http.ProxyFromEnvironment
	}

	client = &nextClient

}

func logError(url string, err error) {
	fields := log.Fields{
		"URL":   url,
		"Error": err,
	}

	// respawn client on next request
	client = nil

	log.WithFields(fields).Error(i18n.T("http_get_error_url", fields))
}

// SetProxy changes proxy usage for all httpclient utilities
// proxyType shall be one of direct, url and env.
// direct - no proxy shall be used, can be used in CLI to override
// configuration settings
// url - the proxyURL must be set for this type.
// golang.org/x/net/proxy.FromURL will be called to get the final proxy
// env - proxy will be load from environment variable, default value
func SetProxy(proxyType string, proxyURL string) {
	proxyConf.proxyType = proxyType
	proxyConf.proxyURL = proxyURL
	log.Debugf("change proxy configuration to %+v", proxyConf)
}

func httpGet(url string) (*http.Response, error) {
	if client == nil {
		newClient()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logError(url, err)
		return nil, err
	}

	return client.Do(req)
}

func validateContentType(resp *http.Response, expected string) error {
	ct := strings.TrimSpace(strings.ToLower(resp.Header.Get("Content-Type")))
	if !strings.HasPrefix(ct, expected) {
		err := errors.New(i18n.T("http_content_type_error_url", i18n.Fields{
			"Expected": expected,
			"Type":     ct,
		}))
		return err
	}
	return nil
}

// HTTPGetJSON retrieves json from specified URL then unmarshal response body
// into data struct. Returns errors happened in http session, or response isn't
// application/json mime type
func HTTPGetJSON(url string, data interface{}) error {
	resp, err := httpGet(url)
	if err != nil {
		logError(url, err)
		return err
	}

	defer resp.Body.Close()

	err = validateContentType(resp, "application/json")
	if err != nil {
		logError(url, err)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logError(url, err)
		return err
	}
	//log.Debugf("response: %#v", string(body))

	err = json.Unmarshal(body, data)
	if err != nil {
		logError(url, err)
		return err
	}
	//log.Debugf("response JSON: %#v", data)

	return nil
}
