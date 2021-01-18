package request

import (
	"encoding/json"
	"errors"
	"net/url"
	"path"
	"regexp"
	"strings"
	pgextract "github.com/Ghamster0/page-extraction/src/extractor"
	"github.com/cfhamlet/os-rq-pod/pkg/sth"
)

// BaseRequest is base request
type BaseRequest struct {
	URL         string            `json:"URL" binding:"required"`
	ContentType string            `json:"ContentType" binding:"required"` //Content-Type in http header
	Content     string            `json:"Content" binding:"required"`     //base64
	UrlParsed            *url.URL
	BaseURL              *url.URL
	Responsec            chan interface{}
}

func (b BaseRequest) Url() string {
	return b.URL
}

// ExtractRequest is Extract related request
type ExtractRequest struct {
	BaseRequest
	Template    *pgextract.Template `json:"template" binding:"required"`
	Exfunc      func (*ExtractRequest)
}

// UnmarshalJSON Deserialization
func (e *ExtractRequest) UnmarshalJSON(b []byte) error {
	type Tmp ExtractRequest
	err := json.Unmarshal(b, (*Tmp)(e))
	if err == nil {
		e.Responsec = make(chan interface{}, 1)
		if e.UrlParsed, err = url.Parse(e.URL); err == nil && e.UrlParsed.Host == "" {
			err = errors.New("empty host")
		}
	}
	return err
}

func(e ExtractRequest) Executor(req interface{}){
	e.Exfunc(req.(*ExtractRequest))
}

var ignoredExt = map[string]struct{}{
	".7z":struct{}{}, ".7zip":struct{}{}, ".bz2":struct{}{}, ".rar":struct{}{}, ".tar":struct{}{}, ".tar.gz":struct{}{}, ".xz":struct{}{}, ".zip":struct{}{},
	".mng":struct{}{}, ".pct":struct{}{}, ".bmp":struct{}{}, ".gif":struct{}{}, ".jpg":struct{}{}, ".jpeg":struct{}{}, ".png":struct{}{}, ".pst":struct{}{},
	".psp":struct{}{}, ".tif":struct{}{}, ".tiff":struct{}{}, ".ai":struct{}{}, ".drw":struct{}{}, ".dxf":struct{}{}, ".eps":struct{}{}, ".ps":struct{}{},
	".svg":struct{}{}, ".cdr":struct{}{}, ".ico":struct{}{},
	".mp3":struct{}{}, ".wma":struct{}{}, ".ogg":struct{}{}, ".wav":struct{}{}, ".ra":struct{}{}, ".aac":struct{}{}, ".mid":struct{}{}, ".au":struct{}{},
	".aiff":struct{}{}, ".3gp":struct{}{}, ".asf":struct{}{}, ".asx":struct{}{}, ".avi":struct{}{}, ".mov":struct{}{}, ".mp4":struct{}{}, ".mpg":struct{}{},
	".qt":struct{}{}, ".rm":struct{}{}, ".swf":struct{}{}, ".wmv":struct{}{}, ".m4a":struct{}{}, ".m4v":struct{}{}, ".flv":struct{}{}, ".webm":struct{}{},
	".xls":struct{}{}, ".xlsx":struct{}{}, ".ppt":struct{}{}, ".pptx":struct{}{}, ".pps":struct{}{}, ".doc":struct{}{}, ".docx":struct{}{}, ".odt":struct{}{},
	".ods":struct{}{}, ".odg":struct{}{}, ".odp":struct{}{}, 
	".css":struct{}{}, ".pdf":struct{}{}, ".exe":struct{}{}, ".bin":struct{}{}, ".rss":struct{}{}, ".dmg":struct{}{}, ".iso":struct{}{}, ".apk":struct{}{},


}

// Request TODO
type Request struct {
	BaseRequest
	IfRegexp             bool     `json:"IfRegexp,omitempty"`
	OnlyHomeSite         bool     `json:"OnlyHomeSite,omitempty"`
	CSSSelectors         []string `json:"CSSSelectors,omitempty"`
	XPathQuerys          []string `json:"XPathQuerys,omitempty"`
	AllowedDomains       []string `json:"AllowedDomains,omitempty"`
	DisallowedDomains    []string `json:"DisallowedDomains,omitempty"`
	DisallowedURLFilters []string `json:"DisallowedURLFilters,omitempty"`
	AllowedURLFilters    []string `json:"AllowedURLFilters,omitempty"`
	AllowedExts          []string `json:"AllowedExts,omitempty"`
	disallowedURLFilters []*regexp.Regexp
	allowedURLFilters    []*regexp.Regexp
	Exfunc               func (*Request)
}

func (r *Request) IsAllowed(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	if r.OnlyHomeSite {
		return parsedURL.Hostname() == r.UrlParsed.Hostname() ||
			strings.HasSuffix(parsedURL.Hostname(), "."+r.UrlParsed.Hostname())
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
	if !r.isDomainAllowed(parsedURL.Hostname()) {
		return false
	}
	if urlTExt := path.Ext(parsedURL.Path); !r.isExtAllowed(urlTExt) {
		return false
	}
	return true
}

func (r *Request) isExtAllowed(urlTExt string) bool {
	for _, allExt := range r.AllowedExts {
		if allExt == urlTExt {
			return true
		}
	}
	if _, ok := ignoredExt[urlTExt];ok{
		return false
	}
	return true
}

func (r *Request) isDomainAllowed(domain string) bool {
	for _, d2 := range r.DisallowedDomains {
		if d2 == domain || strings.HasSuffix(domain, "."+d2) {
			return false
		}
	}
	if r.AllowedDomains == nil || len(r.AllowedDomains) == 0 {
		return true
	}
	for _, d2 := range r.AllowedDomains {
		if d2 == domain || strings.HasSuffix(domain, "."+d2) {
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
	u = strings.TrimSpace(u)
	absURL, err := base.Parse(u)
	if err != nil || absURL.Host == "" {
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
	if err == nil {
		r.Responsec = make(chan interface{}, 1)
		if r.UrlParsed, err = url.Parse(r.URL); err == nil && r.UrlParsed.Host == "" {
			err = errors.New("empty host")
		}
		if len(r.DisallowedURLFilters) > 0 {
			r.disallowedURLFilters = make([]*regexp.Regexp, len(r.DisallowedURLFilters))
			for i, f := range r.DisallowedURLFilters {
				if r.disallowedURLFilters[i], err = regexp.Compile(f); err != nil {
					return err
				}
			}
		}

		if len(r.AllowedURLFilters) > 0 {
			r.allowedURLFilters = make([]*regexp.Regexp, len(r.AllowedURLFilters))
			for i, f := range r.AllowedURLFilters {
				if r.allowedURLFilters[i], err = regexp.Compile(f); err != nil {
					return err
				}
			}
		}
	}
	return err
}

func(e Request) Executor(req interface{}){
	e.Exfunc(req.(*Request))
}

// Response TODO
type Response = sth.Result