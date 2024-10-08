package mongo

import (
	"context"
	"fmt"

	"sip-monitor/src/entity"
)

func (m *MgInfra) SaveCall(item *entity.SIPRecordCall) {
	if m.CollectionRecordCall == nil || item == nil {
		return
	}
	_, err := m.CollectionRecordCall.InsertOne(context.Background(), item)
	if err != nil {
		m.logger.WithError(err).Error("Save Item call Error:")
		return
	}
}

func (m *MgInfra) SaveMsg(item *entity.SIP) {
	if m.CollectionRecord == nil {
		return
	}
	if item.Raw == nil {
		return
	}
	timestamp := uint64(item.CreateTime.Unix()*1000) + item.TimestampMicro/1000
	ctx := context.TODO()
	_, err := m.CollectionRecord.InsertOne(ctx, entity.Record{
		UUID:       item.UUID,
		NodeID:     item.NodeID,
		NodeIP:     item.NodeIP,
		CallID:     item.CallID,
		Method:     item.Title,
		Src:        item.SrcAddr,
		Dst:        item.DstAddr,
		CreateTime: item.CreateTime,
		Timestamp:  timestamp,
		Body:       *item.Raw,
	})

	if err != nil {
		m.logger.WithError(err).Error("Save Item Sip Message Error:")
		return
	}
	m.logger.Debug(fmt.Sprintf("Save Item msg%s(%s) %s->%s", item.CSeqMethod, item.CallID, item.FromUsername+item.FromDomain, item.ToUsername+item.ToDomain))
}
