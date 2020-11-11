package extract

import (
	"encoding/base64"
	"github.com/lichao-mobanche/go-extractor-server/pkg/request"
	pgxtract "github.com/Ghamster0/page-extraction/src/extractor"
)

// Extract extract links
func Extract(r *request.ExtractRequest){
	c, e := base64.StdEncoding.DecodeString(r.Content)
	if e != nil {
		r.Responsec <- Base64Error(e.Error())
		return
	}
	r.Content = string(c)
	resp:=pgxtract.Extract(r.URL,r.Content,r.Template)
	r.Responsec <- resp
	return
}