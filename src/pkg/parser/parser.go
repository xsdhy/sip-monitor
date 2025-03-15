package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"sip-monitor/src/config"
	"sip-monitor/src/entity"
	"sip-monitor/src/pkg/hep"
	"sip-monitor/src/pkg/util"
)

var empty = struct{}{}
var acceptMethods map[string]struct{}
var discardMethods map[string]struct{}

const CRLF = "\r\n"

const (
	ParseOk = iota
	ECanNotFindHeader
	EBadHeaderValue
)

const EmptyStr = ""

type Parser struct {
	cfg    *config.Config
	hepMsg *hep.HepMsg
	sip    entity.SIP
}

func NewParser(cfg *config.Config, hepMsg *hep.HepMsg) *Parser {
	raw := string(hepMsg.Body)
	return &Parser{
		cfg:    cfg,
		hepMsg: hepMsg,
		sip: entity.SIP{
			Raw: &raw,
		},
	}
}

func (p *Parser) ParseSIPMsg() (s *entity.SIP, err error) {
	p.ParseCseq()
	if p.sip.CSeqMethod == "" {
		return nil, errors.New("cseq_is_empty")
	}
	if strings.Contains(p.cfg.DiscardMethods, p.sip.CSeqMethod) {
		return nil, errors.New("method_discarded")
	}

	p.ParseCallID()
	if p.sip.CallID == "" {
		return nil, errors.New("callid_is_empty")
	}

	p.ParseFirstLine()
	if p.sip.Title == "" {
		return nil, errors.New("title_is_empty")
	}

	if p.sip.RequestURL != "" {
		p.ParseRequestURL()
	}

	p.ParseFrom()
	p.ParseTo()
	p.ParseUserAgent()
	p.sip.CreateAt = time.Unix(int64(p.hepMsg.Timestamp), 0)
	p.sip.TimestampMicro = p.sip.CreateAt.Add(time.Microsecond * time.Duration(p.hepMsg.TimestampMicro)).UnixMicro()

	if p.cfg.HeaderFSCallIDName != "" {
		p.ParseFSCallID(p.cfg.HeaderFSCallIDName)
	}

	if p.cfg.HeaderUIDName != "" {
		p.ParseUID(p.cfg.HeaderUIDName)
	}

	p.sip.Protocol = int(p.hepMsg.IPProtocolID)

	p.sip.SrcAddr = fmt.Sprintf("%s_%d", p.hepMsg.IP4SourceAddress, p.hepMsg.SourcePort)
	p.sip.SrcPort = int(p.hepMsg.SourcePort)
	p.sip.SrcHost = p.hepMsg.IP4SourceAddress

	p.sip.DstAddr = fmt.Sprintf("%s_%d", p.hepMsg.IP4DestinationAddress, p.hepMsg.DestinationPort)
	p.sip.DstHost = p.hepMsg.IP4DestinationAddress
	p.sip.DstPort = int(p.hepMsg.DestinationPort)

	p.sip.NodeID = strconv.Itoa(int(p.hepMsg.CaptureAgentID))

	return &p.sip, nil
}

// Request 	: INVITE bob@example.com SIP/2.0
// Response 	: SIP/2.0 200 OK
// Response	: SIP/2.0 501 Not Implemented
func (p *Parser) ParseFirstLine() {
	if p.sip.Raw == nil {
		return
	}
	if *p.sip.Raw == EmptyStr {
		return
	}

	firstLineIndex := strings.Index(*p.sip.Raw, CRLF)
	if firstLineIndex == -1 {
		return
	}
	firstLine := (*p.sip.Raw)[:firstLineIndex]
	firstLineMeta := strings.SplitN(firstLine, " ", 3)

	if len(firstLineMeta) != 3 {
		return
	}
	if strings.HasPrefix(firstLineMeta[0], "SIP") {
		p.sip.IsRequest = false
		p.sip.Title = firstLineMeta[1]
		p.sip.ResponseCode = util.StrToInt(firstLineMeta[1])
		p.sip.ResponseDesc = firstLineMeta[2]
		return
	}
	p.sip.IsRequest = true
	p.sip.Title = firstLineMeta[0]
	p.sip.RequestURL = firstLineMeta[1]
}

func (p *Parser) ParseRequestURL() {
	if p.sip.RequestURL == "" {
		return
	}
	user, domain := ParseSIPURL(p.sip.RequestURL)
	p.sip.RequestDomain = domain
	p.sip.RequestUsername = user
}

func (p *Parser) ParseFrom() {
	v := p.GetHeaderValue(entity.HeaderFrom)
	if v == EmptyStr {
		return
	}
	user, domain := ParseSIPURL(v)
	p.sip.FromUsername = user
	p.sip.FromDomain = domain
}

func (p *Parser) ParseTo() {
	v := p.GetHeaderValue(entity.HeaderTo)
	if v == EmptyStr {
		return
	}
	user, domain := ParseSIPURL(v)
	p.sip.ToUsername = user
	p.sip.ToDomain = domain
}

func (p *Parser) ParseUserAgent() {
	v := p.GetHeaderValue(entity.HeaderUA)
	if v == EmptyStr {
		return
	}
	p.sip.UserAgent = v
}

func (p *Parser) ParseCallID() {
	v := p.GetHeaderValue(entity.HeaderCallID)
	if v == EmptyStr {
		return
	}
	p.sip.CallID = v
}

// "Bob" <sips:bob@biloxi.com> ;tag=a48s
// sip:+12125551212@phone2net.com;tag=887s
// Anonymous <sip:c8oqz84zk7z@privacy.org>;tag=hyh8
// Carol <sip:carol@chicago.com>
// sip:carol@chicago.com
func ParseSIPURL(s string) (string, string) {
	if s == "" {
		return "", ""
	}

	newURL := s

	if strings.Contains(s, "<") {
		start := strings.Index(s, "<")
		end := strings.Index(s, ">")
		if start > end {
			return "", ""
		}
		newURL = s[start:end]
	}

	a := strings.Index(newURL, ":")
	b := strings.Index(newURL, "@")
	c := strings.Index(newURL, ";")

	if a == -1 {
		return "", ""
	}

	if b == -1 && b < len(newURL) {
		if c == -1 {
			return "", newURL[a+1:]
		}
		if c > b {
			return "", newURL[a+1 : c]
		}
	}

	if c == -1 {
		c = len(newURL)
	}

	user := newURL[a+1 : b]
	domain := newURL[b+1 : c]
	return user, domain
}

func (p *Parser) ParseCseq() {
	cseqValue := p.GetHeaderValue(entity.HeaderCSeq)
	if cseqValue == EmptyStr {
		return
	}
	cs := strings.SplitN(cseqValue, " ", 2)
	if len(cs) != 2 {
		return
	}
	p.sip.CSeqNumber = util.StrToInt(cs[0])
	p.sip.CSeqMethod = cs[1]
}

func (p *Parser) GetHeaderValue(header string) (v string) {
	if *p.sip.Raw == EmptyStr || header == EmptyStr {
		return EmptyStr
	}

	if strings.Contains(header, CRLF) || strings.Contains(header, " ") {
		return EmptyStr
	}

	startIndex := strings.Index(*p.sip.Raw, header+":")

	if startIndex == -1 {
		return EmptyStr
	}

	newStr := (*p.sip.Raw)[startIndex:]

	endIndex := strings.Index(newStr, CRLF)

	if endIndex == -1 {
		return EmptyStr
	}

	if len(header)+1 > endIndex {
		return EmptyStr
	}

	return strings.TrimSpace(newStr[len(header)+1 : endIndex])
}

func (p *Parser) ParseUID(HeaderUIDName string) {
	if HeaderUIDName == "" {
		return
	}

	v := p.GetHeaderValue(HeaderUIDName)
	if v == EmptyStr {
		return
	}
	p.sip.UID = v
}

func (p *Parser) ParseFSCallID(FSCallID string) {
	if FSCallID == "" {
		return
	}

	v := p.GetHeaderValue(FSCallID)
	if v == EmptyStr {
		return
	}
	p.sip.FSCallID = v
}

func init() {
	am := map[string]struct{}{
		"INVITE":    empty,
		"CANCEL":    empty,
		"ACK":       empty,
		"BYE":       empty,
		"INFO":      empty,
		"OPTIONS":   empty,
		"UPDATE":    empty,
		"REGISTER":  empty,
		"MESSAGE":   empty,
		"SUBSCRIBE": empty,
		"NOTIFY":    empty,
		"PRACK":     empty,
		"REFER":     empty,
		"PUBLISH":   empty,
	}

	// may be read from env
	dm := map[string]struct{}{
		"INFO":      empty,
		"OPTIONS":   empty,
		"REGISTER":  empty,
		"MESSAGE":   empty,
		"SUBSCRIBE": empty,
		"PUBLISH":   empty,
	}

	acceptMethods = am
	discardMethods = dm
}
