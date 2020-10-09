package extract

import (
	"io/ioutil"
	"bytes"
	"net/url"
	"encoding/base64"
	"strings"
	"regexp"
	xhtml "html"

	"github.com/lichao-mobanche/go-extractor-server/pkg/request"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xmlquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)
// ContantFormat type
type ContantFormat int

const (
	unknown         ContantFormat = iota
	html       // html
	xml        // xml
)

// Extract extract links
func Extract(r *request.Request){
	
	c, e := base64.StdEncoding.DecodeString(r.Content)
	if e!=nil {
		r.Responsec<-Base64Error(e.Error())
		return
	}
	r.Content=string(c)
	if e:=fixCharset(r);e!=nil{
		r.Responsec<-e
		return
	}
	resp:=request.Response{}
	if r.IfRegexp {
		regexpHandler(r,resp)
		r.Responsec<-resp
	}
	switch getFormat(r){
	case html:
		if e:=cssHandler(r,resp);e==nil{
			if e=xpathHtmlHandler(r,resp);e==nil{
				r.Responsec<-resp
			}
		}
		if e!=nil {
			r.Responsec<-e
		}
	case xml:
		if xpathXmlHandler(r,resp);e!=nil{
			r.Responsec<-e
		} else {
			r.Responsec<-resp
		}
	default:
		r.Responsec<-ContentTypeError("unknown")
	}
	return
}

func regexpHandler(r *request.Request, resp request.Response) {
	var baseUrlRe = regexp.MustCompile(`(?i)<base\s[^>]*href\s*=\s*[\"\']\s*([^\"\'\s]+)\s*[\"\']`)
	var l int
	if l=len(r.Content);l>4096{
		l=4096
	}
	if m:=baseUrlRe.FindString(r.Content[0:l]);m!=""{
		r.BaseURL, _ = r.UrlParsed.Parse(m)
	}
	var LinksRe = regexp.MustCompile(`(?is)<a\s.*?href=(\"[.#]+?\"|'[.#]+?'|[^\s]+?)(>|\s.*?>)(.*?)<[/ ]?a>`)
	linksAndTxts := LinksRe.FindAllStringSubmatch(r.Content,-1)
	links:=make([]string,len(linksAndTxts))
	for i, l := range linksAndTxts {
		link:=xhtml.EscapeString(strings.Trim(l[1], "\t\r\n '\"\x0c"))
		if r.IsAllowed(link) {
			links[i]=link
		}
	}
	resp["re"]=links
	return
}

func cssHandler(r *request.Request, resp request.Response) error {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(r.Content))
	if err != nil {
		return DocError(err.Error())
	}
	if href, found := doc.Find("base[href]").Attr("href"); found {
		r.BaseURL, _ = r.UrlParsed.Parse(href)
	}
	for _, selector := range r.CSSSelectors {
		tmpLinks := make([]string, 0)
		doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
			for _, n := range s.Nodes {
				for _, a := range n.Attr {
					if a.Key == "href" {
						link:=r.AbsoluteURL(a.Val)
						if link != "" && r.IsAllowed(link){
							tmpLinks=append(tmpLinks, link)
						}
					}
				}
			}
		})
		resp["css_"+selector]=tmpLinks
	}
	return nil
}

func xpathHtmlHandler(r *request.Request, resp request.Response) error {
	doc, err := htmlquery.Parse(bytes.NewBufferString(r.Content))
	if err != nil {
		return DocError(err.Error())
	}

	if e := htmlquery.FindOne(doc, "//base"); e != nil {
		for _, a := range e.Attr {
			if a.Key == "href" {
				r.BaseURL, _ = r.UrlParsed.Parse(a.Val)
				break
			}
		}
	}
	if len(r.XPathQuerys)==0{
		tmpLinks := make([]string, 0)
		for _, n := range htmlquery.Find(doc, "//a") {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link:=r.AbsoluteURL(a.Val)
					if link != "" && r.IsAllowed(link){
						tmpLinks=append(tmpLinks, link)
					}
				}
			}
		}
		resp["xpath"]=tmpLinks
		return nil
	}
	for _, query := range r.XPathQuerys {
		tmpLinks := make([]string, 0)
		for _, n := range htmlquery.Find(doc, query) {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link:=r.AbsoluteURL(a.Val)
					if link != "" && r.IsAllowed(link){
						tmpLinks=append(tmpLinks, link)
					}
				}
			}
		}
		resp["xpath_"+query]=tmpLinks
	}
	return nil
}

func xpathXmlHandler(r *request.Request, resp request.Response) error {
	doc, err := xmlquery.Parse(bytes.NewBufferString(r.Content))
	if err != nil {
		return err
	}
	if len(r.XPathQuerys)==0{
		tmpLinks := make([]string, 0)
		xmlquery.FindEach(doc, "//a", func(i int, n *xmlquery.Node) {
			for _, a := range n.Attr {
				if a.Name.Local == "href" {
					link:=r.AbsoluteURL(a.Value)
					if link != "" && r.IsAllowed(link){
						tmpLinks=append(tmpLinks, link)
					}
				}
			}
		})
		resp["xpath"]=tmpLinks
		return nil
	}

	for _, query := range r.XPathQuerys {
		tmpLinks := make([]string, 0)
		xmlquery.FindEach(doc, query, func(i int, n *xmlquery.Node) {
			for _, a := range n.Attr {
				if a.Name.Local == "href" {
					link:=r.AbsoluteURL(a.Value)
					if link != "" && r.IsAllowed(link){
						tmpLinks=append(tmpLinks, link)
					}
				}
			}
		})
		resp["xpath_"+query]=tmpLinks
	}
	return nil
}

func xmlFile(u *url.URL) (bool, error) {
	return strings.HasSuffix(strings.ToLower(u.Path), ".xml") || strings.HasSuffix(strings.ToLower(u.Path), ".xml.gz"), nil
}

func getFormat(r *request.Request) ContantFormat {
	res := unknown
	contentType := strings.ToLower(r.ContentType)
	if strings.Contains(contentType, "html"){
		res=html
	} else if strings.Contains(contentType, "xml"){
		res=xml
	} else if isXMLFile, err:=xmlFile(r.UrlParsed);err==nil&&isXMLFile{
		res=xml
	}
	return res
}
func fixCharset(r *request.Request) error {
	if len(r.Content) == 0 {
		return ContentEmptyError(r.URL)
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
		r, err := d.DetectBest([]byte(r.Content))
		if err != nil {
			return DetectorError(err.Error())
		}
		contentType = "text/plain; charset=" + r.Charset
	}
	if strings.Contains(contentType, "utf-8") || strings.Contains(contentType, "utf8") {
		return nil
	}
	tmpContent, err := encodeBytes([]byte(r.Content), contentType)
	if err != nil {
		return EncodeError(err.Error())
	}
	r.Content = string(tmpContent)
	return nil
}

func encodeBytes(b []byte, contentType string) ([]byte, error) {
	r, err := charset.NewReader(bytes.NewReader(b), contentType)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}