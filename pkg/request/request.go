package request

import (
	"encoding/json"
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
	response	chan Response
}

// UnmarshalJSON Deserialization
func (r *Request) UnmarshalJSON(b []byte) error {
	type Tmp Request
	err := json.Unmarshal(b, (*Tmp)(r))
	if err ==nil{
		r.response=make(chan Response)
	}
	return err
}

// Response TODO
type Response struct {
	URL     string	`json:"URL" binding:"required"`
	Links	map[string]interface{} `json:"Links" binding:"required"`
}