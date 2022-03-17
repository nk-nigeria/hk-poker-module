package api

import (
	"context"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	"github.com/heroiclabs/nakama-common/runtime"
)

type ProcessorPackager struct {
	state      *entity.MatchState
	processor  *MatchProcessor
	logger     runtime.Logger
	dispatcher runtime.MatchDispatcher
	messages   []runtime.MatchData
}

func NewProcessorPackage(state *entity.MatchState, processor *MatchProcessor, logger runtime.Logger, dispatcher runtime.MatchDispatcher, messages []runtime.MatchData) *ProcessorPackager {
	return &ProcessorPackager{
		state:      state,
		processor:  processor,
		logger:     logger,
		dispatcher: dispatcher,
		messages:   messages,
	}
}

func (p ProcessorPackager) GetState() *entity.MatchState {
	return p.state
}

func (p ProcessorPackager) GetProcessor() *MatchProcessor {
	return p.processor
}

func (p ProcessorPackager) GetLogger() runtime.Logger {
	return p.logger
}

func (p ProcessorPackager) GetDispatcher() runtime.MatchDispatcher {
	return p.dispatcher
}

func (p ProcessorPackager) GetMessages() []runtime.MatchData {
	return p.messages
}

func GetProcessorPackagerFromContext(ctx context.Context) *ProcessorPackager {
	return ctx.Value(processorKey).(*ProcessorPackager)
}

func GetContextWithProcessorPackager(procPkg *ProcessorPackager) context.Context {
	return context.WithValue(context.TODO(), processorKey, procPkg)
}
