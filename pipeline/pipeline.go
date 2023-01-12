package pipeline

import (
	"context"
)

type StageFunc func(ctx context.Context, inp <-chan int64) (out chan int64)

type Pipeline struct {
	stages []StageFunc
}

func New(stages ...StageFunc) *Pipeline {
	return &Pipeline{
		stages: stages,
	}
}

func (p *Pipeline) Run(ctx context.Context, inp <-chan int64) (out <-chan int64) {
	out = inp
	for _, stage := range p.stages {
		//out = p.runStageInt(p.stages[index], out)
		out = stage(ctx, out)
	}
	return out
}
