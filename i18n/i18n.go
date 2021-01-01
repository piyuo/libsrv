package i18n

import "strings"

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
