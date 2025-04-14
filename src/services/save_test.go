package services

import (
	"encoding/json"
	"sip-monitor/src/entity"
	"testing"
)

// 测试被叫返回486 CallID:c3a31dd0-911c-123e-0e90-00163e0fafd1
func TestSaveService_Save486(t *testing.T) {
	//todo:: 486 Busy Here
}

// 测试被叫返回408 CallID:18c16b42-911d-123e-cba3-fa163ebe21a3
// 408 Request Timeout
func TestSaveService_Save408(t *testing.T) {
	//todo:: 408 Request Timeout
	dataJson := `[
  {
    "id": 19249666,
    "node_ip": "127.0.0.1",
    "sip_call_id": "18a6fee6-911d-123e-0e90-00163e0fafd1",
    "method": "INVITE",
    "response_desc": "",
    "to_user": "mx400529371068715",
    "from_user": "5597184217",
    "src_addr": "47.89.248.143:5060",
    "dst_addr": "192.168.11.141:5080",
    "create_time": "2025-04-11T10:11:18+08:00",
    "timestamp_micro": 1744337478729701
  },
  {
    "id": 19249663,
    "node_ip": "127.0.0.1",
    "sip_call_id": "18a6fee6-911d-123e-0e90-00163e0fafd1",
    "method": "INVITE",
    "response_desc": "",
    "to_user": "mx400529371068715",
    "from_user": "5597184217",
    "src_addr": "47.89.248.143:5060",
    "dst_addr": "192.168.11.141:5080",
    "create_time": "2025-04-11T10:11:18+08:00",
    "timestamp_micro": 1744337478729778
  },
  {
    "id": 19249667,
    "node_ip": "127.0.0.1",
    "sip_call_id": "18a6fee6-911d-123e-0e90-00163e0fafd1",
    "method": "100",
    "response_desc": "Trying",
    "to_user": "mx400529371068715",
    "from_user": "5597184217",
    "src_addr": "192.168.11.141:5080",
    "dst_addr": "47.89.248.143:5060",
    "create_time": "2025-04-11T10:11:18+08:00",
    "timestamp_micro": 1744337478730616
  },
  {
    "id": 19249664,
    "node_ip": "127.0.0.1",
    "sip_call_id": "18a6fee6-911d-123e-0e90-00163e0fafd1",
    "method": "100",
    "response_desc": "Trying",
    "to_user": "mx400529371068715",
    "from_user": "5597184217",
    "src_addr": "192.168.11.141:5080",
    "dst_addr": "47.89.248.143:5060",
    "create_time": "2025-04-11T10:11:18+08:00",
    "timestamp_micro": 1744337478730623
  },
  {
    "id": 19249673,
    "node_ip": "127.0.0.1",
    "sip_call_id": "18a6fee6-911d-123e-0e90-00163e0fafd1",
    "method": "183",
    "response_desc": "Session Progress",
    "to_user": "mx400529371068715",
    "from_user": "5597184217",
    "src_addr": "192.168.11.141:5080",
    "dst_addr": "47.89.248.143:5060",
    "create_time": "2025-04-11T10:11:18+08:00",
    "timestamp_micro": 1744337478875625
  },
  {
    "id": 19249672,
    "node_ip": "127.0.0.1",
    "sip_call_id": "18a6fee6-911d-123e-0e90-00163e0fafd1",
    "method": "183",
    "response_desc": "Session Progress",
    "to_user": "mx400529371068715",
    "from_user": "5597184217",
    "src_addr": "192.168.11.141:5080",
    "dst_addr": "47.89.248.143:5060",
    "create_time": "2025-04-11T10:11:18+08:00",
    "timestamp_micro": 1744337478875632
  },
  {
    "id": 19249714,
    "node_ip": "127.0.0.1",
    "sip_call_id": "18a6fee6-911d-123e-0e90-00163e0fafd1",
    "method": "408",
    "response_desc": "Request Timeout",
    "to_user": "mx400529371068715",
    "from_user": "5597184217",
    "src_addr": "192.168.11.141:5080",
    "dst_addr": "47.89.248.143:5060",
    "create_time": "2025-04-11T10:11:44+08:00",
    "timestamp_micro": 1744337504834595
  },
  {
    "id": 19249713,
    "node_ip": "127.0.0.1",
    "sip_call_id": "18a6fee6-911d-123e-0e90-00163e0fafd1",
    "method": "408",
    "response_desc": "Request Timeout",
    "to_user": "mx400529371068715",
    "from_user": "5597184217",
    "src_addr": "192.168.11.141:5080",
    "dst_addr": "47.89.248.143:5060",
    "create_time": "2025-04-11T10:11:44+08:00",
    "timestamp_micro": 1744337504834604
  },
  {
    "id": 19249716,
    "node_ip": "127.0.0.1",
    "sip_call_id": "18a6fee6-911d-123e-0e90-00163e0fafd1",
    "method": "ACK",
    "response_desc": "",
    "to_user": "mx400529371068715",
    "from_user": "5597184217",
    "src_addr": "47.89.248.143:5060",
    "dst_addr": "192.168.11.141:5080",
    "create_time": "2025-04-11T10:11:44+08:00",
    "timestamp_micro": 1744337504916974
  },
  {
    "id": 19249715,
    "node_ip": "127.0.0.1",
    "sip_call_id": "18a6fee6-911d-123e-0e90-00163e0fafd1",
    "method": "ACK",
    "response_desc": "",
    "to_user": "mx400529371068715",
    "from_user": "5597184217",
    "src_addr": "47.89.248.143:5060",
    "dst_addr": "192.168.11.141:5080",
    "create_time": "2025-04-11T10:11:44+08:00",
    "timestamp_micro": 1744337504917046
  }
]`

	var item []entity.Record
	err := json.Unmarshal([]byte(dataJson), &item)
	if err != nil {
		t.Fatal(err)
	}

	var records []entity.SIP
	for _, v := range item {
		records = append(records, entity.SIP{
			CallID:         v.SIPCallID,
			Title:          v.Method,
			ResponseCode:   v.ResponseCode,
			ResponseDesc:   v.ResponseDesc,
			FromUser:       v.FromUser,
			ToUser:         v.ToUser,
			SrcAddr:        v.SrcAddr,
			DstAddr:        v.DstAddr,
			CreateTime:     v.CreateTime,
			TimestampMicro: v.TimestampMicro,
		})
	}

	t.Log(records)
}

// 测试主叫主动挂断 CallID:bdc167e7-911c-123e-cba3-fa163ebe21a3 / CallID:103cb329-911e-123e-0e90-00163e0fafd1
func TestSaveService_Save487(t *testing.T) {
	//todo:: 487 Busy Everywhere
}

// 测试正常通话-主叫挂断 CallID:f28dcd80-911d-123e-0e90-00163e0fafd1
func TestSaveService_NormalCall(t *testing.T) {
	//todo:: 正常通话
}
