// Package i18n offers utilities to simplify i18n translations.
// github.com/nicksnyder/go-i18n/i18n is used under the hood
package i18n

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nicksnyder/go-i18n/i18n/language"
)

// DefaultLanguage used for i18n.
const DefaultLanguage = "en-us"

var (
	tFunc        i18n.TranslateFunc
	langFilePath string
	langLoaded   *language.Language
)

func translationFile(lang string) string {
	return fmt.Sprintf("%s/%s.all.json", langFilePath, lang)
}

// SetLanguageFilePath changes path to *.all.json files and returns previous
// settings
func SetLanguageFilePath(path string) string {
	lastPath := langFilePath
	langFilePath = path
	return lastPath
}

// LoadTFunc loads translation file and return the T function according to
// the specified language code which is in go-i18n/i18n compatible notation.
// It has side-effect that updates package global T function. The defaultLanguage
// will always be used as fallback
func LoadTFunc(langCandidates ...string) i18n.TranslateFunc {
	i18n.MustLoadTranslationFile(translationFile(DefaultLanguage))
	for _, langCode := range langCandidates {
		i18n.LoadTranslationFile(translationFile(langCode))
	}
	langCandidates = append(langCandidates, DefaultLanguage)
	T, lang := i18n.MustTfuncAndLanguage(langCandidates[0], langCandidates[1:]...)

	tFunc = T
	langLoaded = lang
	return T
}

// GetLoadedLang returns globally loaded language or nil if no language loaded
func GetLoadedLang() *language.Language {
	return langLoaded
}

// Fields is shortcut as translation template fields to save some typing.
// E.g.
//  log.Debug(T("lang_debug_candidate_loaded", Fields{
//    "LangCfg":    langCfg,
//    "LangLoaded": lang,
//  }))
type Fields map[string]interface{}

// T offers shortcut to get translation with go-i18n/i18n.TranslateFunc a
// package global TranslateFunc would be loaded either during caller
// initialization by invoking LoadTFunc or in first T invocation with default
// language (en-us)
func T(translationID string, args ...interface{}) string {
	if tFunc == nil {
		LoadTFunc(DefaultLanguage)
	}
	return tFunc(translationID, args...)
}
