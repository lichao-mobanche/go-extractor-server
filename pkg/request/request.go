package request

import (
	"encoding/json"
	"net/url"
)

// Request TODO
type Request struct {
	URL     string	`json:"URL" binding:"required"`
	ExType 	string	`json:"ExType" binding:"required"`
	CSSSelectors []string `json:"CSSSelectors,omitempty"`
	XPathQuerys []string `json:"XPathQuerys,omitempty"`
	ContentType string	`json:"ContentType" binding:"required"`
	//base64
	Content    string	`json:"Content" binding:"required"`
	url	*url.URL
	baseURL *url.URL
	responsec	chan interface{}
}
func (r *Request) AbsoluteURL(u string) string {
	if strings.HasPrefix(u, "#") {
		return ""
	}
	var base *url.URL
	if r.baseURL != nil {
		base = r.baseURL
	} else {
		base = r.url
	}
	absURL, err := base.Parse(u)
	if err != nil {
		return ""
	}
	absURL.Fragment = ""
	if absURL.Scheme == "//" {
		absURL.Scheme = r.url.Scheme
	}
	return absURL.String()
}

// UnmarshalJSON Deserialization
func (r *Request) UnmarshalJSON(b []byte) error {
	type Tmp Request
	err := json.Unmarshal(b, (*Tmp)(r))
	if err ==nil{
		r.responsec=make(chan interface{},1)
		r.url, err = url.Parse(r.URL)
	}
	return err
}

// Response TODO
// type Response struct {
// 	URL     string	`json:"URL" binding:"required"`
// 	Links	map[string]interface{} `json:"Links" binding:"required"`
// }

// Response TODO
type Response map[string]interface{}