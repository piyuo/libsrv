package i18n

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/piyuo/libsrv/env"
	"github.com/piyuo/libsrv/file"
)

const (
	CacheKey = "i-"
)

var locales = [...]string{"en_US", "zh_CN", "zh_TW"}

// IsPredefined determine a locale is predefined in i18n, return true and predefined locale is exist
//
//   predefined, locale := isPredefined("en-us"); // true,en_US
//
func IsPredefined(locale string) (bool, string) {
	locale = strings.Replace(locale, "-", "_", 1)
	locale = strings.ToLower(locale)
	for _, l := range locales {
		if strings.ToLower(l) == locale {
			return true, l
		}
	}
	return false, ""
}

// GetLocaleFromContext return locale from current request, return en_US if anything else
//
//	locale := GetLocale(ctx)
//
func GetLocaleFromContext(ctx context.Context) string {
	value := ctx.Value(env.KeyContextRequest)
	if value == nil {
		return "en_US"
	}
	req := value.(*http.Request)
	return GetLocaleFromRequest(req)
}

// GetLocaleFromRequest parse http header Accept-Language field and return first one as default language, return en_US if anything else
//
//	defaultLocale := GetLocale(request)
//
func GetLocaleFromRequest(r *http.Request) string {
	return acceptLanguage(r.Header.Get("Accept-Language"))
}

// acceptLanguage parse http header Accept-Language field and match to i18n predefine locale, return 'en_US' if nothing match
//
//	locale := acceptLanguage("da, en-us;q=0.8, en;q=0.7") // "en_US"
//
func acceptLanguage(acptLang string) string {
	//if acptLang is locale like 'en_US', this will speed thing up
	exist, predefined := IsPredefined(acptLang)
	if exist {
		return predefined
	}

	type langQ struct {
		Lang string
		Q    float64
	}

	langQS := []*langQ{}
	accepts := strings.Split(acptLang, ",")
	for _, accept := range accepts {
		accept = strings.Trim(accept, " ")
		args := strings.Split(accept, ";")
		if len(args) == 1 {
			langQS = append(langQS, &langQ{
				Lang: args[0],
				Q:    1,
			})
		} else {
			qp := strings.Split(args[1], "=")
			q, err := strconv.ParseFloat(qp[1], 64)
			if err == nil {
				langQS = append(langQS, &langQ{
					Lang: args[0],
					Q:    q,
				})
			}
		}
	}

	sort.SliceStable(langQS, func(i, j int) bool {
		return langQS[i].Q > langQS[j].Q
	})

	for _, lq := range langQS {
		exist, predefined := IsPredefined(lq.Lang)
		if exist {
			return predefined
		}
	}
	return "en_US"
}

// LocaleFilename get resource key name
//
//   LocaleFilename(ctx, "file1",".json") // "name_zh_TW"
//
func LocaleFilename(ctx context.Context, name, ext string) string {
	return name + "_" + GetLocaleFromContext(ctx) + ext
}

// JSON get i18n resource file in JSON format in current locale
//
//	j, err := JSON(ctx, "filename")
//
func JSON(ctx context.Context, name, ext string, d time.Duration) (map[string]interface{}, error) {
	return file.I18nJSON(LocaleFilename(ctx, name, ext), d)
}

// Text get i18n resource file in text format in current locale
//
//	j, err := Text(ctx, "filename")
//
func Text(ctx context.Context, name, ext string, d time.Duration) (string, error) {
	return file.I18nText(LocaleFilename(ctx, name, ext), d)
}
