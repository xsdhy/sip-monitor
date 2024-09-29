package model

import (
	"sip-monitor/src/pkg/callbuffer"
)

func SaveToDBRunner() {
	for {
		select {
		case item := <-SaveToDBQueue:
			if item.CSeqMethod != "REGISTER" {
				Infra.SaveMsg(item)

				buffer, ok := CallBufferMap[item.UUID]
				if !ok {
					CallBufferMap[item.UUID] = callbuffer.NewCallBuffer()
				}
				result := buffer.Add(item)
				if result != nil {
					Infra.SaveCall(result)
				}
			} else {
				//saveRegister(item)
			}
		}
	}
}
