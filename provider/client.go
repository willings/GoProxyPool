package provider

import "net/http"

func addBotHeader(h http.Header) {
	h.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	h.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	h.Add("Accept-Encoding", "gzip, deflate, sdch")
	h.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.6,en;q=0.4")
	h.Add("Cache-control", "upgrade-insecure-requests")
	h.Add("Upgrade-Insecure-Requests", "1")
	h.Add("Connection", "Keep-Alive")
}
