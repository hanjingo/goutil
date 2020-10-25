package steper

import (
	"container/list"
	"context"
	"errors"
)

type Round struct {
	stepList  *list.List
	status    uint32
	closeFunc context.CancelFunc
	currElem  *list.Element
	loopTimes int64
}

func NewRound() *Round {
	back := &Round{
		stepList:  list.New(),
		loopTimes: 0,
	}
	back.stepList.PushFront(NewStep(ROUND_START, startFunc)) //所有round的开头
	return back
}

func (r *Round) Run() {
	if r.stepList.Len() == 0 {
		return
	}
	r.currElem = r.stepList.Front()
	ctx, cancel := context.WithCancel(context.Background())
	r.closeFunc = cancel
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case r.status = <-r.currElem.Value.(*Step).Do():
				if r.status == ROUND_DONE {
					return
				}
				if r.status == ROUND_START {
					r.loopTimes++
				}
				r.currElem = r.currElem.Next()
				if r.currElem == nil || r.currElem.Value == nil {
					return
				}
				continue
			}
		}
	}()
}

func (r *Round) Destroy() {
	if r.closeFunc != nil {
		r.closeFunc()
	}
	if r.stepList.Len() > 0 {
		r.Range(r.DelStep)
	}
}

func (r *Round) Reset() {
	r.currElem = r.stepList.Front()
}

func (r *Round) PushBack(step *Step) {
	r.stepList.PushBack(step)
}

func (r *Round) PushFront(step *Step) {
	r.insertStep(r.stepList.Front().Value.(*Step), step)
}

func (r *Round) InsertStep(prevStep *Step, step *Step) {
	r.insertStep(prevStep, step)
}

func (r *Round) insertStep(prevStep *Step, step *Step) (*list.Element, error) {
	prevE := r.getElemByStep(prevStep)
	if prevE == nil {
		return nil, errors.New("上一步为空")
	}
	return r.stepList.InsertAfter(step, prevE), nil
}

func (r *Round) DelStep(step *Step) {
	e := r.getElemByStep(step)
	if e == nil {
		return
	}
	step.Destroy()
	r.stepList.Remove(e)
}

func (r *Round) Range(f func(step *Step)) {
	if r.stepList.Len() == 0 {
		return
	}
	i := r.stepList.Len()
	for i > 1 {
		e := r.stepList.Front().Next() //跳过start函数
		if e == nil {
			return
		}
		f(e.Value.(*Step))
		i--
	}
}

func (r *Round) getElemByStep(step *Step) *list.Element {
	if r.stepList.Len() == 0 {
		return nil
	}
	e := r.stepList.Front()
	for {
		if e == nil {
			return nil
		}
		if e.Value.(*Step) == step {
			return e
		}
		e = e.Next()
	}
}

func (r *Round) GetCurrStatus() uint32 {
	return r.status
}

func (r *Round) GetCurrStep() *Step {
	if r.currElem == nil {
		return nil
	}
	if r.currElem.Value == nil {
		return nil
	}
	return r.currElem.Value.(*Step)
}

func startFunc(ctx context.Context) uint32 {
	return ROUND_START
}
