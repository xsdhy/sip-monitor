package hep

import (
	_ "embed"
	"reflect"
	"testing"

	"sip-monitor/src/entity"
)

func TestHepServer_saveCall(t *testing.T) {

	//messages := ReadContentAndSplit(call1)

	hep, _ := NewHepServer(nil, nil)

	//ip := net.ParseIP("127.0.0.1")

	sips := make([]*entity.SIP, 0)
	//for _, item := range messages {
	//	sipItem := hep.parse(item, ip)
	//	assert.Nil(t, sipItem)
	//	sips = append(sips, sipItem)
	//}

	tests := []struct {
		name string
		want *entity.SIPRecordCall
	}{
		{name: "1", want: &entity.SIPRecordCall{
			ID:              "",
			UUID:            "",
			NodeID:          "",
			NodeIP:          "",
			CallID:          "fb8t6s3gq27qep30e4os",
			FromUser:        "10001",
			ToUser:          "17311225659",
			UserAgent:       "JsSIP 3.9.0",
			SrcHost:         "",
			SrcPort:         0,
			SrcAddr:         "",
			SrcCountryName:  "",
			SrcCityName:     "",
			DstHost:         "",
			DstPort:         0,
			DstAddr:         "",
			DstCountryName:  "",
			DstCityName:     "",
			CreateTime:      nil,
			RingingTime:     nil,
			AnswerTime:      nil,
			EndTime:         nil,
			CallDuration:    0,
			RingingDuration: 0,
			TalkDuration:    0,
			HangupCode:      "",
			HangupReason:    "",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *entity.SIPRecordCall
			for _, item := range sips {
				got = hep.saveCall(item)
				if got != nil {
					break
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("saveCall() = %v, want %v", got, tt.want)
			}
		})
	}
}
