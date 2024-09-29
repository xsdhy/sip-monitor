package mongo

import (
	"context"
	"fmt"
	"log/slog"
	"sip-monitor/src/entity"
)

func (m *MgInfra) SaveCall(item *entity.SIPRecordCall) {
	if m.CollectionRecordCall == nil || item == nil {
		return
	}
	_, err := m.CollectionRecordCall.InsertOne(context.Background(), item)
	if err != nil {
		slog.Error("Save Item call Error:", err.Error())
		return
	}
	slog.Debug("Save Item call", slog.String("msg", fmt.Sprintf("%s(%s) %s->%s", "inteve", item.CallID, item.FromUser+item.SrcHost, item.ToUser+item.DstHost)))
}

func (m *MgInfra) SaveMsg(item *entity.SIP) {
	if m.CollectionRecord == nil {
		return
	}
	if item.Raw == nil {
		return
	}
	ctx := context.TODO()
	_, err := m.CollectionRecord.InsertOne(ctx, entity.Record{
		UUID:       item.UUID,
		NodeID:     item.NodeID,
		NodeIP:     item.NodeIP,
		CallID:     item.CallID,
		CreateTime: item.CreateTime,
		Body:       *item.Raw,
	})

	if err != nil {
		slog.Error("Save Item Sip Message Error:", err.Error())
		return
	}
	slog.Debug("Save Item", slog.String("msg", fmt.Sprintf("%s(%s) %s->%s", item.CSeqMethod, item.CallID, item.FromUsername+item.FromDomain, item.ToUsername+item.ToDomain)))
}
