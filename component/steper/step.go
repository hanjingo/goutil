package steper

import (
	"context"
)

type StepFunc func(ctx context.Context) uint32

type Step struct {
	bValid     bool
	f          StepFunc
	cancelFunc context.CancelFunc
	doneChan   chan uint32
	id         uint32
}

func NewStep(id uint32, f func(ctx context.Context) uint32) *Step {
	back := &Step{
		bValid: true,
		f:      f,
		id:     id,
	}
	back.doneChan = make(chan uint32, 1)
	return back
}

func (s *Step) GetId() uint32 {
	return s.id
}

func (s *Step) Do() chan uint32 {
	defer close(s.doneChan)
	if !s.bValid {
		return s.StepOver()
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	s.doneChan <- s.f(ctx)
	return s.doneChan
}

//让当前
func (s *Step) Stop() {
	s.bValid = false
}

//中断
func (s *Step) Suspend() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}

//跳过
func (s *Step) StepOver() chan uint32 {
	s.doneChan <- STEP_DONE
	return s.doneChan
}

//摧毁
func (s *Step) Destroy() {
	if !s.bValid {
		return
	}
	s.bValid = false
	s.Suspend()
}
