package request

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strings"
	"github.com/cfhamlet/os-rq-pod/pkg/sth"
)

// Request TODO
type Request struct {
	URL     string	`json:"URL" binding:"required"`
	IfRegexp 	bool	`json:"IfRegexp,omitempty"`
	OnlyHomeSite bool   `json:"OnlyHomeSite,omitempty"`
	//Content-Type in http header
	ContentType string	`json:"ContentType" binding:"required"`
	//base64
	Content    string	`json:"Content" binding:"required"`
	CSSSelectors []string `json:"CSSSelectors,omitempty"`
	XPathQuerys []string `json:"XPathQuerys,omitempty"`
	AllowedDomains []string `json:"AllowedDomains,omitempty"`
	DisallowedDomains []string `json:"DisallowedDomains,omitempty"`
	DisallowedURLFilters []string `json:"DisallowedURLFilters,omitempty"`
	AllowedURLFilters []string `json:"AllowedURLFilters,omitempty"`
	disallowedURLFilters []*regexp.Regexp
	allowedURLFilters []*regexp.Regexp
	UrlParsed	*url.URL
	BaseURL *url.URL
	Responsec	chan interface{}
}

func (r *Request)IsAllowed(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	if r.OnlyHomeSite {
		return parsedURL.Hostname()==r.BaseURL.Hostname()||
		strings.HasSuffix(parsedURL.Hostname(), "."+r.BaseURL.Hostname())
	}
	if len(r.disallowedURLFilters) > 0 {
		if r.isMatchingFilter(r.disallowedURLFilters, []byte(u)) {
			return false
		}
	}
	if len(r.allowedURLFilters) > 0 {
		if !r.isMatchingFilter(r.allowedURLFilters, []byte(u)) {
			return false
		}
	}
	return r.isDomainAllowed(parsedURL.Hostname())
}

func (r *Request) isDomainAllowed(domain string) bool {
	for _, d2 := range r.DisallowedDomains {
		if d2 == domain {
			return false
		}
	}
	if r.AllowedDomains == nil || len(r.AllowedDomains) == 0 {
		return true
	}
	for _, d2 := range r.AllowedDomains {
		if d2 == domain {
			return true
		}
	}
	return false
}

func (r *Request) isMatchingFilter(fs []*regexp.Regexp, d []byte) bool {
	for _, r := range fs {
		if r.Match(d) {
			return true
		}
	}
	return false
}

func (r *Request) AbsoluteURL(u string) string {
	if strings.HasPrefix(u, "#") {
		return ""
	}
	var base *url.URL
	if r.BaseURL != nil {
		base = r.BaseURL
	} else {
		base = r.UrlParsed
	}
	absURL, err := base.Parse(u)
	if err != nil {
		return ""
	}
	absURL.Fragment = ""
	if absURL.Scheme == "//" {
		absURL.Scheme = r.UrlParsed.Scheme
	}
	return absURL.String()
}

// UnmarshalJSON Deserialization
func (r *Request) UnmarshalJSON(b []byte) error {
	type Tmp Request
	err := json.Unmarshal(b, (*Tmp)(r))
	if err ==nil{
		r.Responsec=make(chan interface{},1)
		r.UrlParsed, err = url.Parse(r.URL)
		if(len(r.DisallowedURLFilters)>0){
			r.disallowedURLFilters=make([]*regexp.Regexp, len(r.DisallowedURLFilters))
			for i, f := range r.DisallowedURLFilters{
				if r.disallowedURLFilters[i],err=regexp.Compile(f);err!=nil{
					return err
				}
			}
		}

		if(len(r.AllowedURLFilters)>0){
			r.allowedURLFilters=make([]*regexp.Regexp, len(r.AllowedURLFilters))
			for i, f := range r.AllowedURLFilters{
				if r.allowedURLFilters[i],err=regexp.Compile(f);err!=nil{
					return err
				}
			}
		}
	}
	return err
}

// Response TODO
type Response sth.Result