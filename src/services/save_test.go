package services

import (
	"context"
	"sip-monitor/src/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试记录
var testInsertions int
var testUpdates int

// 重置测试环境
func resetTestEnv() {
	testInsertions = 0
	testUpdates = 0
	callRecordCache = make(map[string]*entity.SIPRecordCall)
}

// 测试内存缓存中的通话生命周期管理
func TestCallLifecycle(t *testing.T) {
	// 准备测试环境
	resetTestEnv()

	// 准备测试数据
	callID := "test-call-id"
	nodeIP := "192.168.1.1"

	// 模拟INVITE消息
	inviteTime := time.Now()
	rawMsg := "SIP/2.0 INVITE"
	rawPtr := &rawMsg
	invite := entity.SIP{
		Title:        "INVITE",
		CSeqMethod:   "INVITE",
		CallID:       callID,
		NodeIP:       nodeIP,
		ToUsername:   "to-user",
		FromUsername: "from-user",
		UserAgent:    "test-ua",
		SrcHost:      "src-host",
		SrcPort:      5060,
		SrcAddr:      "src-addr",
		DstHost:      "dst-host",
		DstPort:      5060,
		DstAddr:      "dst-addr",
		CreateAt:     inviteTime,
		Raw:          rawPtr,
	}

	// 处理INVITE
	updateCallRecordInCache(invite)

	// 验证缓存创建
	assert.Equal(t, 1, len(callRecordCache))
	record, exists := callRecordCache[callID]
	assert.True(t, exists)
	assert.Equal(t, "to-user", record.ToUser)
	assert.Equal(t, "from-user", record.FromUser)
	assert.Equal(t, "test-ua", record.UserAgent)

	// 模拟180 Ringing
	ringingTime := inviteTime.Add(1 * time.Second)
	ringing := entity.SIP{
		Title:      "180",
		CSeqMethod: "INVITE",
		CallID:     callID,
		CreateAt:   ringingTime,
		Raw:        rawPtr,
	}

	// 处理Ringing
	updateCallRecordInCache(ringing)

	// 验证Ringing状态
	record, exists = callRecordCache[callID]
	assert.True(t, exists)
	assert.NotNil(t, record.RingingTime)
	assert.Equal(t, ringingTime, *record.RingingTime)

	// 验证Ringing记录未结束
	assert.Nil(t, record.EndTime)

	// 模拟200 OK答复
	answerTime := ringingTime.Add(2 * time.Second)
	answer := entity.SIP{
		Title:      "200",
		CSeqMethod: "INVITE",
		CallID:     callID,
		CreateAt:   answerTime,
		Raw:        rawPtr,
	}

	// 处理应答
	updateCallRecordInCache(answer)

	// 验证应答状态
	record, exists = callRecordCache[callID]
	assert.True(t, exists)
	assert.NotNil(t, record.AnswerTime)
	assert.Equal(t, answerTime, *record.AnswerTime)

	// 验证通话未结束
	assert.Nil(t, record.EndTime)
}

// 测试缓存过期策略
func TestCacheExpiry(t *testing.T) {
	// 清空缓存
	resetTestEnv()

	// 临时替换数据库更新函数
	origUpdateOne := internalUpdateOne
	internalUpdateOne = func(ctx context.Context, filter interface{}, update interface{}, opts ...interface{}) (interface{}, error) {
		testUpdates++
		return nil, nil
	}

	// 测试结束后恢复原始函数
	defer func() {
		internalUpdateOne = origUpdateOne
	}()

	// 创建三个记录：已完成、过期、活跃
	now := time.Now()

	// 已完成的呼叫
	completedCall := &entity.SIPRecordCall{
		ID:        "completed-id",
		SIPCallID: "completed-call-id",
		NodeIP:    "192.168.1.1",
	}
	createTime1 := now.Add(-30 * time.Minute)
	endTime := now.Add(-20 * time.Minute)
	completedCall.CreateTime = &createTime1
	completedCall.EndTime = &endTime

	// 过期的呼叫（开始超过15分钟但未结束）
	oldCall := &entity.SIPRecordCall{
		ID:        "old-id",
		SIPCallID: "old-call-id",
		NodeIP:    "192.168.1.1",
	}
	createTime2 := now.Add(-20 * time.Minute)
	oldCall.CreateTime = &createTime2

	// 活跃的呼叫
	activeCall := &entity.SIPRecordCall{
		ID:        "active-id",
		SIPCallID: "active-call-id",
		NodeIP:    "192.168.1.1",
	}
	createTime3 := now.Add(-5 * time.Minute)
	activeCall.CreateTime = &createTime3

	// 添加到缓存
	callRecordCache["completed-call-id"] = completedCall
	callRecordCache["old-call-id"] = oldCall
	callRecordCache["active-call-id"] = activeCall

	// 检查初始状态
	assert.Equal(t, 3, len(callRecordCache))

	// 调用刷新函数
	FlushCacheToDB()

	// 验证结果
	assert.Equal(t, 2, testUpdates)          // 应该有2个记录被更新（已完成和过期的）
	assert.Equal(t, 1, len(callRecordCache)) // 只有活跃的记录应该保留
	_, exists := callRecordCache["active-call-id"]
	assert.True(t, exists) // 活跃记录应该存在
}

// 测试从INVITE到BYE的完整呼叫生命周期
func TestCompleteSIPCall(t *testing.T) {
	// 准备测试环境
	resetTestEnv()

	// 临时替换数据库更新函数
	origUpdateOne := internalUpdateOne
	internalUpdateOne = func(ctx context.Context, filter interface{}, update interface{}, opts ...interface{}) (interface{}, error) {
		testUpdates++
		return nil, nil
	}

	// 临时替换数据库插入函数
	origInsertOne := internalInsertOne
	internalInsertOne = func(ctx context.Context, document interface{}, opts ...interface{}) (interface{}, error) {
		testInsertions++
		return nil, nil
	}

	// 测试结束后恢复原始函数
	defer func() {
		internalUpdateOne = origUpdateOne
		internalInsertOne = origInsertOne
	}()

	// 准备测试数据
	callID := "complete-flow-call-id"
	nodeIP := "192.168.1.1"
	rawMsg := "SIP/2.0 INVITE"
	rawPtr := &rawMsg

	// 1. INVITE请求
	inviteTime := time.Now()
	invite := entity.SIP{
		NodeIP:       nodeIP,
		Title:        "INVITE",
		CSeqMethod:   "INVITE",
		CallID:       callID,
		ToUsername:   "alice",
		FromUsername: "bob",
		UserAgent:    "test-ua",
		SrcHost:      "src-host",
		SrcPort:      5060,
		SrcAddr:      "src-addr",
		DstHost:      "dst-host",
		DstPort:      5060,
		DstAddr:      "dst-addr",
		CreateAt:     inviteTime,

		Raw:                    rawPtr,
		TimestampMicroWithDate: 123456789,
	}

	// 处理INVITE
	SaveOptimized(invite)
	time.Sleep(10 * time.Millisecond) // 给goroutine一点时间执行

	// 验证记录插入
	assert.Equal(t, 1, testInsertions)

	// 2. 180 Ringing
	ringingTime := inviteTime.Add(1 * time.Second)
	ringing := entity.SIP{
		NodeIP:     nodeIP,
		Title:      "180",
		CSeqMethod: "INVITE",
		CallID:     callID,
		CreateAt:   ringingTime,
		Raw:        rawPtr,
	}

	// 处理Ringing
	SaveOptimized(ringing)

	// 3. 200 OK (应答)
	answerTime := ringingTime.Add(2 * time.Second)
	answer := entity.SIP{
		NodeIP:     nodeIP,
		Title:      "200",
		CSeqMethod: "INVITE",
		CallID:     callID,
		CreateAt:   answerTime,
		Raw:        rawPtr,
	}

	// 处理应答
	SaveOptimized(answer)

	// 4. BYE (挂断)
	endTime := answerTime.Add(30 * time.Second)
	bye := entity.SIP{
		NodeIP:     nodeIP,
		Title:      "200",
		CSeqMethod: "BYE",
		CallID:     callID,
		CreateAt:   endTime,
		Raw:        rawPtr,
	}

	// 处理挂断
	SaveOptimized(bye)

	// 验证呼叫记录已被保存并从缓存中移除
	assert.Equal(t, 4, testInsertions)       // 应该有4个原始SIP消息被保存
	assert.Equal(t, 1, testUpdates)          // 应该有1个呼叫记录被更新
	assert.Equal(t, 0, len(callRecordCache)) // 缓存应该为空
}
