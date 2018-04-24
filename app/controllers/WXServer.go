package controllers

import (
	"bytes"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/revel/revel"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"github.com/tsingakbar/leanote/app/db"
	"gopkg.in/mgo.v2/bson"
)

type WXServer struct {
	BaseController
}

var (
	alphaSpace = regexp.MustCompile(`[^0-9a-z ]`)
)

type DictWordItem struct {
	Id    string `bson:"_id,omitempty"`
	Strip string `bson:"strip"`
	Word  string `bson:"word"`
	Trans string `bson:"trans"`
}

type textXML struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	Others       []Node   `xml:",any"`
}

type Node struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

func (c WXServer) verifyWXReq() bool {
	var params [3]string
	params[0], _ = revel.Config.String("wxToken")
	params[1] = c.Params.Query.Get("timestamp")
	params[2] = c.Params.Query.Get("nonce")
	var paramsSorted = params[:]
	sort.Strings(paramsSorted)
	params2Hash := strings.Join(paramsSorted, "")
	hasher := sha1.New()
	io.WriteString(hasher, params2Hash)
	return fmt.Sprintf("%x", hasher.Sum(nil)) == c.Params.Query.Get("signature")
}

func (c WXServer) EchoStr() revel.Result {
	if !c.verifyWXReq() {
		return c.Forbidden("")
	}
	return c.RenderText(c.Params.Query.Get("echostr"))
}

type InvalidRequestResponse string

func (r InvalidRequestResponse) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusBadRequest, "text/plain")
	resp.GetWriter().Write([]byte(r))
}

func ProcessTextReqContent(xmlRecv *textXML) {
	if len(xmlRecv.Content) == 0 {
		xmlRecv.Content = "I felt the void on ur lips."
		return
	}

	//if xmlRecv.FromUserName == "oHFH2wCS_20Qdnd_ztqPX1gsJUyM" || xmlRecv.FromUserName == "oHFH2wOTU8M9zE0eN6BBu3qM7wGc" {
	//	runeContent := []rune(xmlRecv.Content)
	//	prefix := runeContent[0]
	//}

	stripKey := strings.Trim(alphaSpace.ReplaceAllString(strings.ToLower(xmlRecv.Content), ""), " \t\r\n")

	queryMgo := db.DictEnEn.Find(bson.M{"strip": bson.M{"$regex": "^" + stripKey}}).Limit(9)
	cnt, err := queryMgo.Count()
	if err != nil {
		xmlRecv.Content = err.Error()
		return
	}

	if cnt == 0 {
		// value of nil or other types than []byte will failed here
		var u url.URL
		u.Scheme = "http"
		u.Host = "m.youdao.com"
		u.Path = "dict"
		var q = make(url.Values)
		q.Set("le", "eng")
		q.Set("q", xmlRecv.Content)
		u.RawQuery = q.Encode()
		xmlRecv.Content = "没找到哎，去<a href =\"" + u.String() + "\">有道</a>查查吧"
	} else {
		// NOTE: wx has a 2048 bytes limit for text message,
		// make sure ur data set is trimmed.
		var candiBuffer bytes.Buffer
		items := queryMgo.Iter()
		item := DictWordItem{}
		for items.Next(&item) {
			if cnt == 1 || item.Word == xmlRecv.Content {
				xmlRecv.Content = item.Trans
				return
			}
			candiBuffer.WriteString(item.Word)
			candiBuffer.WriteString("\n")
		}
		xmlRecv.Content = candiBuffer.String()
	}
}

func (c WXServer) DoProcess() revel.Result {
	if !c.verifyWXReq() {
		return c.Forbidden("")
	}
	var xmlRecv textXML
	if err := xml.NewDecoder(c.Request.GetBody()).Decode(&xmlRecv); err != nil {
		return InvalidRequestResponse(err.Error())
	}
	revel.AppLog.Debug("wxincom",
		"from", xmlRecv.FromUserName,
		"to", xmlRecv.ToUserName,
		"type", xmlRecv.MsgType)
	revel.AppLog.Debugf("%+v", xmlRecv)

	if xmlRecv.MsgType == "text" {
		ProcessTextReqContent(&xmlRecv)
	} else {
		xmlRecv.Content = "[暂不支持回复此类型消息]"
		if xmlRecv.MsgType == "event" {
			for i, _ := range xmlRecv.Others {
				if xmlRecv.Others[i].XMLName.Local == "Event" && xmlRecv.Others[i].Content == "<![CDATA[subscribe]]>" {
					xmlRecv.Content = "你刚刚订阅了一个英文词典公众号，表达了你对生活的热爱和对知识的渴求。"
					break
				}
			}
		}
		xmlRecv.MsgType = "text"
	}
	xmlRecv.FromUserName, xmlRecv.ToUserName = xmlRecv.ToUserName, xmlRecv.FromUserName
	xmlRecv.CreateTime = time.Now().Unix()
	xmlRecv.Others = nil // should not marshal the "Others" back as response xml
	return c.RenderXML(xmlRecv)
}
