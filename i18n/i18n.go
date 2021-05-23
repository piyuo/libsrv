package i18n

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/piyuo/libsrv/cache"
	"github.com/piyuo/libsrv/env"
	"github.com/piyuo/libsrv/file"
	"github.com/pkg/errors"
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

// ResourceKey get resource key name
//
//   So(ResourceKey(ctx, "file1",".json"), ShouldEqual, "name_zh_TW")
//
func ResourceKey(ctx context.Context, name, ext string) string {
	return name + "_" + GetLocaleFromContext(ctx) + ext
}

// ResourcePath get resource file path
//
//   So(ResourcePath(ctx, "name",".json"), ShouldEqual, "assets/i18n/name_zh_TW.json")
//
func ResourcePath(ctx context.Context, name, ext string) string {
	return "assets/i18n/" + ResourceKey(ctx, name, ext)
}

// Resource get i18n resource file in JSON format
//
//	json, err := Resource(ctx, "notExist",".json")
//
func Resource(ctx context.Context, name, ext string) (map[string]interface{}, error) {
	keyname := CacheKey + ResourceKey(ctx, name, ext)
	found, bytes, err := cache.Get(keyname)
	if err != nil {
		return nil, errors.Wrap(err, "get cache "+keyname)
	}
	if found {
		j := make(map[string]interface{})
		if err := json.Unmarshal(bytes, &j); err != nil {
			return nil, errors.Wrapf(err, "decode cache json %v", keyname)
		}
		return j, nil
	}

	j, bytes, err := ResourceWithoutCache(ctx, name, ext)
	if err != nil {
		return nil, errors.Wrapf(err, "get resource %v%v", name, ext)
	}

	if err := cache.Set(keyname, bytes, 0); err != nil {
		return nil, errors.Wrap(err, "set cache "+keyname)
	}
	return j, nil
}

// ResourceWithoutCache get i18n resource file without cache, no error if not found
//
//	json, err := ResourceWithoutCache(ctx, "notExist",".json")
//
func ResourceWithoutCache(ctx context.Context, name, ext string) (map[string]interface{}, []byte, error) {
	filepath, found := file.Lookup(ResourcePath(ctx, name, ext))
	if !found {
		return nil, nil, nil
	}

	bytes, err := file.Read(filepath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "read bytes %v", filepath)
	}
	j := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &j); err != nil {
		return nil, nil, errors.Wrapf(err, "decode cache json %v", name)
	}
	return j, bytes, nil
}
