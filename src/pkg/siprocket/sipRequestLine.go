package siprocket

/*
 RFC 3261 - https://www.ietf.org/rfc/rfc3261.txt

INVITE sip:01798300765@87.252.61.202;user=phone SIP/2.0
SIP/2.0 200 OK

*/

type sipReq struct {
	Method     []byte // Sip Method eg INVITE etc
	UriType    string // Type of URI sip, sips, tel etc
	StatusCode []byte // Status Code eg 100
	StatusDesc []byte // Status Code Description eg trying
	User       []byte // User part
	Host       []byte // Host part
	Port       []byte // Port number
	UserType   []byte // User Type
	Src        []byte // Full source if needed
}

func parseSipReq(v []byte, out *sipReq) {
	pos := 0
	state := FIELD_NULL

	// Init the output area
	out.UriType = ""
	out.Method = nil
	out.StatusCode = nil
	out.User = nil
	out.Host = nil
	out.Port = nil
	out.UserType = nil
	out.Src = nil

	// Keep the source line if needed
	if keep_src {
		out.Src = v
	}

	// 先检查是否为响应行 (SIP/2.0 开头)
	if len(v) > 7 && getString(v, 0, 7) == "SIP/2.0" {
		out.Method = []byte("SIP/2.0")

		// 跳过 "SIP/2.0 " 找到状态码
		pos = 8
		for pos < len(v) && v[pos] != ' ' && v[pos] != '\r' && v[pos] != '\n' {
			out.StatusCode = append(out.StatusCode, v[pos])
			pos++
		}

		// 跳过空格，找到状态描述
		if pos < len(v) && v[pos] == ' ' {
			pos++
			for pos < len(v) && v[pos] != '\r' && v[pos] != '\n' {
				out.StatusDesc = append(out.StatusDesc, v[pos])
				pos++
			}
		}

		return
	}

	// 如果不是响应行，则处理请求行
	for pos < len(v) {
		// FSM
		switch state {
		case FIELD_NULL:
			// 请求行
			if v[pos] >= 'A' && v[pos] <= 'Z' {
				state = FIELD_METHOD
				continue
			}

		case FIELD_METHOD:
			if v[pos] == ' ' {
				state = FIELD_BASE
				pos++
				continue
			}
			out.Method = append(out.Method, v[pos])

		case FIELD_BASE:
			if v[pos] != ' ' {
				// Not a space so check for uri types
				if getString(v, pos, pos+4) == "sip:" {
					state = FIELD_USER
					pos = pos + 4
					out.UriType = "sip"
					continue
				}
				if getString(v, pos, pos+5) == "sips:" {
					state = FIELD_USER
					pos = pos + 5
					out.UriType = "sips"
					continue
				}
				if getString(v, pos, pos+4) == "tel:" {
					state = FIELD_USER
					pos = pos + 4
					out.UriType = "tel"
					continue
				}
				if getString(v, pos, pos+5) == "user=" {
					state = FIELD_USERTYPE
					pos = pos + 5
					continue
				}
			}
		case FIELD_USER:
			if v[pos] == '@' {
				state = FIELD_HOST
				pos++
				continue
			}
			// 修复：在到达@之前的所有内容都添加到用户部分
			out.User = append(out.User, v[pos])

		case FIELD_HOST:
			if v[pos] == ':' {
				state = FIELD_PORT
				pos++
				continue
			}
			if v[pos] == ';' || v[pos] == '>' || v[pos] == ' ' {
				state = FIELD_BASE
				pos++
				continue
			}
			out.Host = append(out.Host, v[pos])

		case FIELD_PORT:
			if v[pos] == ';' || v[pos] == '>' || v[pos] == ' ' {
				state = FIELD_BASE
				pos++
				continue
			}
			out.Port = append(out.Port, v[pos])

		case FIELD_USERTYPE:
			if v[pos] == ';' || v[pos] == '>' || v[pos] == ' ' {
				state = FIELD_BASE
				pos++
				continue
			}
			out.UserType = append(out.UserType, v[pos])
		}
		pos++
	}
}
