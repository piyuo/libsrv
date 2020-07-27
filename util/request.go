package util

import (
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

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

// GetLanguage parse http header Accept-Language field and return first one as default language
//
//	defaultLanguage := GetLanguage(request)
//
func GetLanguage(r *http.Request) string {
	languages := parseAcceptLanguage(r.Header.Get("Accept-Language"))
	return languages[0]
}

// GetAcceptLanguage parse http header Accept-Language field to sorted string list
//
//	list := GetAcceptLanguage(request)
//
func GetAcceptLanguage(r *http.Request) []string {
	return parseAcceptLanguage(r.Header.Get("Accept-Language"))
}

// parseAcceptLanguage parse http header Accept-Language field to sorted string list
//
//	list := parseAcceptLanguage("da, en-gb;q=0.8, en;q=0.7") // []string{"da","en-gb","en"}
//
func parseAcceptLanguage(acptLang string) []string {
	if acptLang == "" {
		return []string{"en-us"}
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

	result := []string{}
	for _, lq := range langQS {
		result = append(result, strings.ToLower(lq.Lang))
	}

	return result
}
