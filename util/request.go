package util

import (
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"

	useragent "github.com/mileusna/useragent"
	"github.com/piyuo/libsrv/i18n"
)

// GetUserAgentID return short id from user agent. no version in here cause we used this for refresh token
//
//	txt := GetUserAgentID(request) // "iPhone,iOS,Safari"
//
func GetUserAgentID(r *http.Request) string {
	u := useragent.Parse(r.UserAgent())
	return u.Device + "," + u.OS + "," + u.Name
}

// GetUserAgentString return short string with version info from user agent
//
//	txt := GetUserAgentString(request) // "iPhone,iOS 7.0,Safari 6.0"
//
func GetUserAgentString(r *http.Request) string {
	u := useragent.Parse(r.UserAgent())
	return u.Device + "," + u.OS + " " + u.OSVersion + "," + u.Name + " " + u.Version
}

// ParseUserAgent return browser name,browser version,os name,os version,device from user agent
//
//	browserName,browserVer,osName,osVer,device := ParseUserAgent(ua)
//
func ParseUserAgent(ua string) (string, string, string, string, string) {
	u := useragent.Parse(ua)
	return u.Name, u.Version, u.OS, u.OSVersion, u.Device
}

// GetUserAgent return user agent
//
//	ua := GetUserAgent(request) // "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit/546.10 (KHTML, like Gecko) Version/6.0 Mobile/7E18WD Safari/8536.25"
//
func GetUserAgent(r *http.Request) string {
	return r.UserAgent()
}

// GetIP return ip from request
//
//	ip := GetIP(request)
//
func GetIP(r *http.Request) string {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	if ips != "" {
		splitIps := strings.Split(ips, ",")
		for _, ip := range splitIps {
			netIP := net.ParseIP(ip)
			if netIP != nil {
				return ip
			}
		}
	}

	//Get IP from RemoteAddr
	if r.RemoteAddr != "" {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return ""
		}
		netIP = net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}
	return ""
}

// GetLocale parse http header Accept-Language field and return first one as default language
//
//	defaultLocale := GetLocale(request)
//
func GetLocale(r *http.Request) string {
	return acceptLanguage(r.Header.Get("Accept-Language"))
}

// acceptLanguage parse http header Accept-Language field and match to i18n predefine locale, return 'en_US' if nothing match
//
//	locale := acceptLanguage("da, en-us;q=0.8, en;q=0.7") // "en_US"
//
func acceptLanguage(acptLang string) string {
	//if acptLang is locale like 'en_US', this will speed thing up
	exist, predefined := i18n.IsPredefined(acptLang)
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
		exist, predefined := i18n.IsPredefined(lq.Lang)
		if exist {
			return predefined
		}
	}
	return "en_US"
}
