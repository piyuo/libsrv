package util

import (
	"net"
	"net/http"
	"strings"

	useragent "github.com/mileusna/useragent"
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
