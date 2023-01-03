package pipeline

import (
	"context"
	"log"
	"os"
)

type StageFunc func(ctx context.Context, inp <-chan int64) (out chan int64)

type Pipeline struct {
	lgr    *log.Logger
	lgrErr *log.Logger
	stages []StageFunc
}

func New(stages ...StageFunc) *Pipeline {
	return &Pipeline{
		lgr:    log.New(os.Stdout, "", 0),
		lgrErr: log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Llongfile),
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
