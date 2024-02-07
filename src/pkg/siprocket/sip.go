package siprocket

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sip-monitor/src/entity"
	"strings"
)

var sip_type = 0
var keep_src = true

type SipMsg struct {
	Req      sipReq
	From     sipFrom
	To       sipTo
	Contact  sipContact
	Via      []sipVia
	Cseq     sipCseq
	Ua       sipVal
	Exp      sipVal
	MaxFwd   sipVal
	CallId   sipVal
	ContType sipVal
	ContLen  sipVal

	Sdp SdpMsg

	SessionID string //自定义的Header头

	Raw string
}

type SdpMsg struct {
	MediaDesc sdpMediaDesc
	Attrib    []sdpAttrib
	ConnData  sdpConnData
}

type sipVal struct {
	Value []byte // Sip Value
	Src   []byte // Full source if needed
}

func (s *sipVal) ToString() string {
	return string(s.Value)
}

func ParseSIP(v []byte) (output *entity.SIP) {
	parse := Parse(v)
	if parse == nil {
		return nil
	}

	output = &entity.SIP{
		CallID:    string(parse.CallId.Value),
		SessionID: parse.SessionID,

		ResponseCode: BytesToInt(parse.Req.StatusCode),
		ResponseDesc: string(parse.Req.StatusDesc),

		ToUser:   string(parse.To.User),
		FromUser: string(parse.From.User),

		CSeqNumber: BytesToInt(parse.Cseq.Id),
		CSeqMethod: string(parse.Cseq.Method),
		UserAgent:  string(parse.Ua.Value),

		Raw: &parse.Raw,
	}

	method := string(parse.Req.Method)
	if method == "SIP/2.0" {
		output.Title = string(parse.Req.StatusCode)
		output.IsRequest = false
	} else {
		output.Title = string(parse.Req.Method)
		output.IsRequest = true
	}

	return output
}

// Main parsing routine, passes by value
func Parse(v []byte) (output *SipMsg) {
	if len(v) <= 0 {
		return nil
	}
	output = new(SipMsg)
	output.Raw = string(v)

	// Allow multiple vias and media Attribs
	via_idx := 0
	output.Via = make([]sipVia, 0, 8)
	attr_idx := 0
	output.Sdp.Attrib = make([]sdpAttrib, 0, 8)

	// 分隔SIP头部和SDP内容
	parts := bytes.SplitN(v, []byte("\r\n\r\n"), 2)

	// 获取SIP头部部分
	headerBytes := parts[0]

	// 处理第一行（请求行或状态行）
	firstLineEnd := bytes.Index(headerBytes, []byte("\r\n"))
	if firstLineEnd == -1 {
		// 如果没有换行符，整个消息就是请求行
		firstLineEnd = len(headerBytes)
	}

	firstLine := headerBytes[:firstLineEnd]
	parseSipReq(firstLine, &output.Req)

	// 处理剩余的SIP头部
	if firstLineEnd < len(headerBytes) {
		headers := headerBytes[firstLineEnd+2:] // 跳过\r\n
		headerLines := bytes.Split(headers, []byte("\r\n"))

		for _, line := range headerLines {
			if len(line) == 0 {
				continue
			}

			// 查找冒号分隔符
			colonPos := bytes.Index(line, []byte(":"))
			if colonPos == -1 {
				continue // 跳过无效行
			}

			headerName := strings.ToLower(string(bytes.TrimSpace(line[:colonPos])))
			headerVal := bytes.TrimSpace(line[colonPos+1:])

			// 根据头部名称处理
			switch headerName {
			case "from", "f":
				parseSipFrom(headerVal, &output.From)
			case "to", "t":
				parseSipTo(headerVal, &output.To)
			case "via", "v":
				var tmpVia sipVia
				output.Via = append(output.Via, tmpVia)
				parseSipVia(headerVal, &output.Via[via_idx])
				via_idx++
			case "contact", "m":
				parseSipContact(headerVal, &output.Contact)
			case "call-id", "i":
				output.CallId.Value = headerVal
			case "content-type", "c":
				output.ContType.Value = headerVal
			case "content-length":
				output.ContLen.Value = headerVal
			case "user-agent":
				output.Ua.Value = headerVal
			case "expires":
				output.Exp.Value = headerVal
			case "max-forwards":
				output.MaxFwd.Value = headerVal
			case "cseq":
				parseSipCseq(headerVal, &output.Cseq)
			case "x-jcallid":
				output.SessionID = string(headerVal)
			}
		}
	}

	// 处理SDP内容（如果存在）
	if len(parts) > 1 && len(parts[1]) > 0 {
		sdpContent := parts[1]
		sdpLines := bytes.Split(sdpContent, []byte("\r\n"))

		for _, line := range sdpLines {
			line = bytes.TrimSpace(line)
			if len(line) < 2 || line[1] != '=' {
				continue // 无效SDP行
			}

			// SDP行格式为 x=value
			sdpType := line[0]
			sdpValue := bytes.TrimSpace(line[2:])

			switch sdpType {
			case 'm':
				parseSdpMediaDesc(sdpValue, &output.Sdp.MediaDesc)
			case 'c':
				parseSdpConnectionData(sdpValue, &output.Sdp.ConnData)
			case 'a':
				var tmpAttrib sdpAttrib
				output.Sdp.Attrib = append(output.Sdp.Attrib, tmpAttrib)
				parseSdpAttrib(sdpValue, &output.Sdp.Attrib[attr_idx])
				attr_idx++
			}
		}
	}

	return
}

// Finds the first valid Seperate or notes its type
func indexSep(s []byte) (int, byte) {

	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			return i, ':'
		}
		if s[i] == '=' {
			return i, '='
		}
	}
	return -1, ' '
}

// Get a string from a slice of bytes
// Checks the bounds to avoid any range errors
func getString(sl []byte, from, to int) string {
	// Remove negative values
	if from < 0 {
		from = 0
	}
	if to < 0 {
		to = 0
	}
	// pageSize if over len
	if from > len(sl) || from > to {
		return ""
	}
	if to > len(sl) {
		return string(sl[from:])
	}
	return string(sl[from:to])
}

// Get a slice from a slice of bytes
// Checks the bounds to avoid any range errors
func getBytes(sl []byte, from, to int) []byte {
	// Remove negative values
	if from < 0 {
		from = 0
	}
	if to < 0 {
		to = 0
	}
	// pageSize if over len
	if from > len(sl) || from > to {
		return nil
	}
	if to > len(sl) {
		return sl[from:]
	}
	return sl[from:to]
}

// Function to print all we know about the struct in a readable format
func (data *SipMsg) PrintSipStruct() {
	fmt.Println("-SIP --------------------------------")

	fmt.Println("  [REQ]")
	fmt.Println("    [UriType] =>", data.Req.UriType)
	fmt.Println("    [Method] =>", string(data.Req.Method))
	fmt.Println("    [StatusCode] =>", string(data.Req.StatusCode))
	fmt.Println("    [User] =>", string(data.Req.User))
	fmt.Println("    [Host] =>", string(data.Req.Host))
	fmt.Println("    [Port] =>", string(data.Req.Port))
	fmt.Println("    [UserType] =>", string(data.Req.UserType))
	fmt.Println("    [Src] =>", string(data.Req.Src))

	// FROM
	fmt.Println("  [FROM]")
	fmt.Println("    [UriType] =>", data.From.UriType)
	fmt.Println("    [Name] =>", string(data.From.Name))
	fmt.Println("    [User] =>", string(data.From.User))
	fmt.Println("    [Host] =>", string(data.From.Host))
	fmt.Println("    [Port] =>", string(data.From.Port))
	fmt.Println("    [Tag] =>", string(data.From.Tag))
	fmt.Println("    [Src] =>", string(data.From.Src))
	// TO
	fmt.Println("  [TO]")
	fmt.Println("    [UriType] =>", data.To.UriType)
	fmt.Println("    [Name] =>", string(data.To.Name))
	fmt.Println("    [User] =>", string(data.To.User))
	fmt.Println("    [Host] =>", string(data.To.Host))
	fmt.Println("    [Port] =>", string(data.To.Port))
	fmt.Println("    [Tag] =>", string(data.To.Tag))
	fmt.Println("    [UserType] =>", string(data.To.UserType))
	fmt.Println("    [Src] =>", string(data.To.Src))
	// TO
	fmt.Println("  [Contact]")
	fmt.Println("    [UriType] =>", data.Contact.UriType)
	fmt.Println("    [Name] =>", string(data.Contact.Name))
	fmt.Println("    [User] =>", string(data.Contact.User))
	fmt.Println("    [Host] =>", string(data.Contact.Host))
	fmt.Println("    [Port] =>", string(data.Contact.Port))
	fmt.Println("    [Transport] =>", string(data.Contact.Tran))
	fmt.Println("    [Q] =>", string(data.Contact.Qval))
	fmt.Println("    [Expires] =>", string(data.Contact.Expires))
	fmt.Println("    [Src] =>", string(data.Contact.Src))
	// UA
	fmt.Println("  [Cseq]")
	fmt.Println("    [Id] =>", string(data.Cseq.Id))
	fmt.Println("    [Method] =>", string(data.Cseq.Method))
	fmt.Println("    [Src] =>", string(data.Cseq.Src))
	// UA
	fmt.Println("  [User Agent]")
	fmt.Println("    [Value] =>", string(data.Ua.Value))
	fmt.Println("    [Src] =>", string(data.Ua.Src))
	// Exp
	fmt.Println("  [Expires]")
	fmt.Println("    [Value] =>", string(data.Exp.Value))
	fmt.Println("    [Src] =>", string(data.Exp.Src))
	// MaxFwd
	fmt.Println("  [Max Forwards]")
	fmt.Println("    [Value] =>", string(data.MaxFwd.Value))
	fmt.Println("    [Src] =>", string(data.MaxFwd.Src))
	// CallId
	fmt.Println("  [Call-ID]")
	fmt.Println("    [Value] =>", string(data.CallId.Value))
	fmt.Println("    [Src] =>", string(data.CallId.Src))
	// Content-Type
	fmt.Println("  [Content-Type]")
	fmt.Println("    [Value] =>", string(data.ContType.Value))
	fmt.Println("    [Src] =>", string(data.ContType.Src))

	// Via - Multiple
	fmt.Println("  [Via]")
	for i, via := range data.Via {
		fmt.Println("    [", i, "]")
		fmt.Println("      [Tansport] =>", via.Trans)
		fmt.Println("      [Host] =>", string(via.Host))
		fmt.Println("      [Port] =>", string(via.Port))
		fmt.Println("      [Branch] =>", string(via.Branch))
		fmt.Println("      [Rport] =>", string(via.Rport))
		fmt.Println("      [Maddr] =>", string(via.Maddr))
		fmt.Println("      [ttl] =>", string(via.Ttl))
		fmt.Println("      [Recevied] =>", string(via.Rcvd))
		fmt.Println("      [Src] =>", string(via.Src))
	}

	fmt.Println("-SDP --------------------------------")
	// Media Desc
	fmt.Println("  [MediaDesc]")
	fmt.Println("    [MediaType] =>", string(data.Sdp.MediaDesc.MediaType))
	fmt.Println("    [Port] =>", string(data.Sdp.MediaDesc.Port))
	fmt.Println("    [Proto] =>", string(data.Sdp.MediaDesc.Proto))
	fmt.Println("    [Fmt] =>", string(data.Sdp.MediaDesc.Fmt))
	fmt.Println("    [Src] =>", string(data.Sdp.MediaDesc.Src))
	// Connection Data
	fmt.Println("  [ConnData]")
	fmt.Println("    [AddrType] =>", string(data.Sdp.ConnData.AddrType))
	fmt.Println("    [ConnAddr] =>", string(data.Sdp.ConnData.ConnAddr))
	fmt.Println("    [Src] =>", string(data.Sdp.ConnData.Src))

	// Attribs - Multiple
	fmt.Println("  [Attrib]")
	for i, attr := range data.Sdp.Attrib {
		fmt.Println("    [", i, "]")
		fmt.Println("      [Cat] =>", string(attr.Cat))
		fmt.Println("      [Val] =>", string(attr.Val))
		fmt.Println("      [Src] =>", string(attr.Src))
	}
	fmt.Println("-------------------------------------")

}

func (data *SipMsg) ToJson() string {
	marshal, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func BytesToInt(bys []byte) int {
	if bys == nil {
		return 0
	}

	// 尝试将字节解析为字符串数字
	str := string(bys)
	var i int
	if _, err := fmt.Sscanf(str, "%d", &i); err == nil {
		return i
	}

	// 如果不是字符串数字，尝试二进制读取
	var data int64
	_ = binary.Read(bytes.NewBuffer(bys), binary.BigEndian, &data)
	return int(data)
}

func BytesToInt64(bys []byte) int64 {
	if bys == nil {
		return 0
	}

	// 尝试将字节解析为字符串数字
	str := string(bys)
	var i int64
	if _, err := fmt.Sscanf(str, "%d", &i); err == nil {
		return i
	}

	// 如果不是字符串数字，尝试二进制读取
	var data int64
	_ = binary.Read(bytes.NewBuffer(bys), binary.BigEndian, &data)
	return data
}

const FIELD_NULL = 0
const FIELD_BASE = 1
const FIELD_VALUE = 2
const FIELD_NAME = 3
const FIELD_NAMEQ = 4
const FIELD_USER = 5
const FIELD_HOST = 6
const FIELD_PORT = 7
const FIELD_TAG = 8
const FIELD_ID = 9
const FIELD_METHOD = 10
const FIELD_TRAN = 11
const FIELD_BRANCH = 12
const FIELD_RPORT = 13
const FIELD_MADDR = 14
const FIELD_TTL = 15
const FIELD_REC = 16
const FIELD_EXPIRES = 17
const FIELD_Q = 18
const FIELD_USERTYPE = 19
const FIELD_STATUS = 20
const FIELD_STATUSDESC = 21

const FIELD_ADDRTYPE = 40
const FIELD_CONNADDR = 41
const FIELD_MEDIA = 42
const FIELD_PROTO = 43
const FIELD_FMT = 44
const FIELD_CAT = 45

const FIELD_IGNORE = 255
