package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	_ "image/jpeg" // for jpeg decoder registration
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/genzj/gobingwallpaper/i18n"
	"github.com/genzj/gobingwallpaper/log"
	"golang.org/x/net/proxy"
)

var once sync.Once
var clientSingleton *http.Client

var proxyConf struct {
	sync.Mutex
	proxyType string
	proxyURL  string
}

func getClient() *http.Client {
	once.Do(newClient)
	return clientSingleton
}

func newClient() {
	proxyConf.Lock()
	defer proxyConf.Unlock()

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
		}

		log.Debugf("SOCKS5 proxy %T installed from URL %v", dialer, proxyURL)
		transport.Dial = dialer.Dial
		transport.Proxy = nil
	default:
		transport.Proxy = http.ProxyFromEnvironment
	}

	clientSingleton = &nextClient

}

func logError(url string, err error) {
	fields := log.Fields{
		"URL":   url,
		"Error": err,
	}

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
	proxyConf.Lock()
	defer proxyConf.Unlock()

	proxyConf.proxyType = proxyType
	proxyConf.proxyURL = proxyURL

	// respawn client on next request
	once = sync.Once{}

	log.Debugf("change proxy configuration to %+v", proxyConf)
}

func issueGetRequest(url string) (*http.Response, error) {
	client := getClient()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
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

func getAndRead(url string, expectedType string) ([]byte, error) {
	resp, err := issueGetRequest(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if expectedType != "" {
		err = validateContentType(resp, expectedType)
		if err != nil {
			return nil, err
		}
	}

	return ioutil.ReadAll(resp.Body)
}

// HTTPGet issues a get request to specified HTTP endpoint then return
// response content. If expectedType is not zero value, Content-Type in
// response header will be checked to match the specified type prefix
func HTTPGet(url string, expectedType string) ([]byte, error) {
	body, err := getAndRead(url, expectedType)
	if err != nil {
		logError(url, err)
		return nil, err
	}

	return body, nil
}

// HTTPGetJSON retrieves json from specified URL then unmarshal response body
// into data struct. Returns errors happened in http session, or response isn't
// application/json mime type
func HTTPGetJSON(url string, data interface{}) error {
	body, err := HTTPGet(url, "application/json")
	if err != nil {
		logError(url, err)
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		logError(url, err)
		return err
	}
	//log.Debugf("response JSON: %#v", data)

	return nil
}

// HTTPGetJpeg downloads jpeg image from given URL then return parsed image
// struct. Error will be returned if any exception occurs during http request
// or response isn't in image/jpeg type.
func HTTPGetJpeg(url string) (image.Image, error) {
	body, err := HTTPGet(url, "image/jpeg")
	if err != nil {
		logError(url, err)
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		logError(url, err)
		return nil, err
	}
	//log.Debugf("response JSON: %#v", data)

	return img, nil
}
