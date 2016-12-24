package i18n

import (
	"fmt"
	"strings"

	"github.com/genzj/go-bingwallpaper/log"
	"github.com/genzj/go-bingwallpaper/util"
	"github.com/nicksnyder/go-i18n/i18n"
)

const DEFAULT_LANGUAGE = "en-us"

var tFunc i18n.TranslateFunc

func translationFile(lang string) string {
	return fmt.Sprintf("%v/i18n/%v.all.json", util.ExecutableFolder(), lang)
}

func LoadTFunc(langCfg string) i18n.TranslateFunc {
	langCfg = strings.ToLower(langCfg)
	defaultLang := DEFAULT_LANGUAGE

	i18n.MustLoadTranslationFile(translationFile(defaultLang))
	i18n.LoadTranslationFile(translationFile(langCfg))
	T, lang := i18n.MustTfuncAndLanguage(langCfg, defaultLang)
	log.Debug(T("lang_debug_candidate_loaded", Fields{
		"LangCfg":    langCfg,
		"LangLoaded": lang,
	}))

	tFunc = T
	return T
}

type Fields map[string]interface{}

func T(translationID string, args ...interface{}) string {
	if tFunc == nil {
		LoadTFunc(DEFAULT_LANGUAGE)
	}
	return tFunc(translationID, args...)
}
