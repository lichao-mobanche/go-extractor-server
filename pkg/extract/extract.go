package extract

import (
	"io/ioutil"
	"bytes"
	"net/url"
	"encoding/base64"
	"strings"
	"github.com/lichao-mobanche/go-extractor-server/pkg/request"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xmlquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

const (
	unknown         ContantFormat = iota
	html       // html
	xml        // xml
)

// Extract extract links
func Extract(r *request.Request){
	if r.Content, e := base64.StdEncoding.DecodeString(r.Content);e!=nil{
		r.responsec<-Base64Error(e)
		return
	}
	if e:=fixCharset(r);e!=nil{
		r.responsec<-e
		return
	}
	resp:=request.Response{}
	switch getFormat(r){
	case html:
		if e:=cssHandler(r,resp);e==nil{
			if e=xpathHtmlHandler(r,resp);e==nil{
				r.responsec<-resp
			}
		}
		if e!=nil {
			r.responsec<-e
		}
		
	case xml:
		resp, e:=xpathHandler(r)
		if e!=nil{
			r.responsec<-e
		} else {
			r.responsec<-resp
		}
	default:
		r.responsec<-ContentTypeError("unknown")
	}
	return
}

func cssHandler(r *request.Request, resp request.Response, ) error {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(r.Content))
	if err != nil {
		return DocError(err)
	}
	if href, found := doc.Find("base[href]").Attr("href"); found {
		r.baseURL, _ = r.url.Parse(href)
	}
	for _, selector := range r.CSSSelectors {
		tmpLinks := make([]string, 0)
		doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
			for _, n := range s.Nodes {
				link:=AbsoluteURL(n.Attr("href"))
				if link != ""{
					tmpLinks=append(tmpLinks, link)
				}
			}
		})
		resp["css_"+selector]=tmpLinks
	}
	return nil
}

func xpathHtmlHandler(r *request.Request, resp request.Response) error {
	doc, err := htmlquery.Parse(bytes.NewBuffer(resp.Body))
	if err != nil {
		return DocError(err)
	}

	if e := htmlquery.FindOne(doc, "//base"); e != nil {
		for _, a := range e.Attr {
			if a.Key == "href" {
				r.baseURL, _ = r.url.Parse(a.Val)
				break
			}
		}
	}

	for _, query := range r.XPathQuerys {
		for _, n := range htmlquery.Find(doc, query) {
			tmpLinks := make([]string, 0)
			for _, a := range n.Attr {
				if a.Key == "href" {
					link:=AbsoluteURL(n.Attr("href"))
					if link != ""{
						tmpLinks=append(tmpLinks, link)
					}
				}
			}
		}
		resp["xpath_"+query]=tmpLinks
	}
	return nil
}

func xmlFile(u string) (bool, err) {
	parsed, err := url.Parse(u)
	if err != nil {
		return parsedURL, err
	}
	return strings.HasSuffix(strings.ToLower(Parsed.Path), ".xml") || strings.HasSuffix(strings.ToLower(Parsed.Path), ".xml.gz"), nil
}

func getFormat(r *request.Request) ContantFormat {
	res := unknown
	contentType := strings.ToLower(r.ContentType)
	if strings.Contains(contentType, "html"){
		res=html
	} else if strings.Contains(contentType, "xml"){
		res=xml
	} else if isXMLFile, err:=xmlFile(r.URL);err==nil&&isXMLFile{
		res=xml
	}
	return res
}
func fixCharset(r *request.Request) error {
	if len(r.Content) == 0 {
		return ContentEmptyError()
	}

	contentType := strings.ToLower(r.ContentType)

	if strings.Contains(contentType, "image/") ||
		strings.Contains(contentType, "video/") ||
		strings.Contains(contentType, "audio/") ||
		strings.Contains(contentType, "font/") {

		return ContentTypeError(contentType)
	}

	if !strings.Contains(contentType, "charset") {
		d := chardet.NewTextDetector()
		r, err := d.DetectBest(r.Content)
		if err != nil {
			return DetectorError(err)
		}
		contentType = "text/plain; charset=" + r.Charset
	}
	if strings.Contains(contentType, "utf-8") || strings.Contains(contentType, "utf8") {
		return nil
	}
	tmpContent, err := encodeBytes(r.Body, contentType)
	if err != nil {
		return EncodeError(err)
	}
	r.ContentType = tmpContent
	return nil
}

func encodeBytes(b []byte, contentType string) ([]byte, error) {
	r, err := charset.NewReader(bytes.NewReader(b), contentType)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}