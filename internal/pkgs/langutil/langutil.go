// Package langutil provides utility functions for language detection.
package langutil

import (
	"context"
	"net/http"
	"strings"
	"unicode"

	"github.com/mozillazg/go-unidecode"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

var languages = []language.Tag{
	language.English, // The first language is used as fallback.
	language.Turkish,
}
var matcher = language.NewMatcher(languages)

func LanguageFromRequest(r *http.Request) language.Tag {
	return LangFromAcceptLanguage(r.Header.Get("Accept-Language"))
}

func LangFromAcceptLanguage(str string) language.Tag {
	accept, _, _ := language.ParseAcceptLanguage(str)
	_, i, _ := matcher.Match(accept...)
	return languages[i]
}

type contextKey string

const langContextKey contextKey = "lang"

func WithLang(ctx context.Context, lang language.Tag) context.Context {
	return context.WithValue(ctx, langContextKey, lang)
}

func LanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(WithLang(r.Context(), LanguageFromRequest(r)))
			next.ServeHTTP(w, r)
		},
	)
}

func Language(ctx context.Context) language.Tag {
	l := ctx.Value(langContextKey)
	if l != nil {
		return l.(language.Tag)
	}
	return languages[0]
}

// Searchable converts text to more search friendly form by removing accent and converting to uppercase.
func Searchable(text string) string {

	// decode with unidecode to escape some known values
	text = unidecode.Unidecode(text)

	// remove punctuation, space etc.
	t := runes.Remove(runes.Predicate(func(r rune) bool {
		return unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsSymbol(r)
	}))
	text, _, _ = transform.String(t, text)

	return strings.ToUpper(text)
}
