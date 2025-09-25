package internal

import (
	"fmt"
	"net/http"
	"strings"
)

// cookieInjectorTransport 实现 http.RoundTripper，
// 在每次请求前自动将 cookie 注入到 header 中。
type cookieInjectorTransport struct {
	base    http.RoundTripper
	cookies []*http.Cookie
}

func (t *cookieInjectorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 自动注入 Cookie header
	if len(t.cookies) > 0 {
		var cookiePairs []string
		for _, c := range t.cookies {
			cookiePairs = append(cookiePairs, fmt.Sprintf("%s=%s", c.Name, c.Value))
		}
		req.Header.Set("Cookie", strings.Join(cookiePairs, "; "))
	}
	return t.base.RoundTrip(req)
}
