package callbuffer

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"sip-monitor/src/entity"
)

func TestNewCallBuffer(t *testing.T) {
	logger := logrus.New()
	saveCall := func(item *entity.SIPRecordCall) {}
	
	cb := NewCallBuffer(logger, saveCall)
	
	assert.NotNil(t, cb)
	assert.Equal(t, StageInit, cb.stage)
	assert.NotNil(t, cb.stageChange)
	assert.NotNil(t, cb.timer)
}

func TestCallBuffer_listener(t *testing.T) {
	logger := logrus.New()
	saveCalled := false
	saveCall := func(item *entity.SIPRecordCall) {
		saveCalled = true
	}
	
	cb := NewCallBuffer(logger, saveCall)
	
	// 测试超时情况
	cb.timer.Reset(100 * time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	assert.True(t, saveCalled)
	
	// 测试状态变化
	cb = NewCallBuffer(logger, saveCall)
	go func() {
		cb.stageChange <- StageAnswer
		time.Sleep(50 * time.Millisecond)
		cb.stageChange <- StageEnd
	}()
	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, StageEnd, cb.stage)
}

func TestCallBuffer_Add(t *testing.T) {
	logger := logrus.New()
	saveCall := func(item *entity.SIPRecordCall) {}
	
	cb := NewCallBuffer(logger, saveCall)
	
	testCases := []struct {
		name     string
		sip      *entity.SIP
		expected CallStage
	}{
		{
			name: "INVITE",
			sip: &entity.SIP{
				CSeqMethod: "INVITE",
				Title:      "INVITE",
			},
			expected: StageCreate,
		},
		{
			name: "BYE",
			sip: &entity.SIP{
				CSeqMethod: "BYE",
				Title:      "BYE",
			},
			expected: StageEnd,
		},
		// 添加更多测试用例...
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cb.Add(tc.sip)
			assert.Equal(t, tc.expected, cb.stage)
		})
	}
}