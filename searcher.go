package fourfours

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

type Searcher struct {
	memo      map[string]*map[float64]Equation
	unaryOps  []UnaryOp
	binaryOps []BinaryOp
}

func NewSearcher(unaryOps []UnaryOp, binaryOps []BinaryOp) Searcher {
	return Searcher{
		memo:      make(map[string]*map[float64]Equation),
		unaryOps:  unaryOps,
		binaryOps: binaryOps,
	}
}

type unaryOpTask struct {
	UnaryOp
	Equation
}

type binaryOpTask struct {
	BinaryOp
	lhs Equation
	rhs Equation
}

type taskResult struct {
	Equation
	ok bool
}

func (s *Searcher) Search(key string) *map[float64]Equation {
	if knowledge, ok := s.memo[key]; ok {
		return knowledge
	}

	fmt.Printf("Start calculating %s\n", key)

	knowledge := make(map[float64]Equation)
	e, ok := NewEquationWithNumber(key)
	if !ok {
		panic("Invalid number")
	}
	knowledge[e.value] = e

	for i := 1; i < len(key); i++ {
		keyLeft := key[:i]
		keyRight := key[i:]
		knowledgeLeft := *s.Search(keyLeft)
		knowledgeRight := *s.Search(keyRight)
		fmt.Printf("Combining %s and %s\n", keyLeft, keyRight)
		tasks := make(chan binaryOpTask)
		results := make(chan taskResult)
		done := make(chan struct{})
		var counter int64

		// spawn workers
		for i := 0; i < runtime.NumCPU(); i++ {
			go func(tasks <-chan binaryOpTask, results chan<- taskResult) {
				for task := range tasks {
					newEq, ok := task.BinaryOp.Apply(&task.lhs, &task.rhs)
					results <- taskResult{
						Equation: newEq,
						ok:       ok,
					}
				}
			}(tasks, results)
		}

		// put jobs
		go func() {
			for _, op := range s.binaryOps {
				for _, el := range knowledgeLeft {
					for _, er := range knowledgeRight {
						atomic.AddInt64(&counter, 1)
						tasks <- binaryOpTask{BinaryOp: op, lhs: el, rhs: er}
					}
				}
			}
			close(done)
			close(tasks)
		}()

	WaitBinaryOp:
		// wait for results
		for {
			select {
			case result := <-results:
				if result.ok {
					if oldEq, ok := knowledge[result.value]; !ok || oldEq.cost > result.cost {
						knowledge[result.value] = result.Equation
					}
				}
				atomic.AddInt64(&counter, -1)
			case <-done:
				if counter == 0 {
					break WaitBinaryOp
				}
			}
		}
	}

	for i := 0; i < 3; i++ {
		oldKnowledge := make(map[float64]Equation, len(knowledge))
		for k, v := range knowledge {
			oldKnowledge[k] = v
		}
		fmt.Printf("Applying Unary Ops to %s (%d/3)\n", key, i+1)
		tasks := make(chan unaryOpTask)
		results := make(chan taskResult)
		done := make(chan struct{})
		var counter int64

		// spawn workers
		for i := 0; i < runtime.NumCPU(); i++ {
			go func(tasks <-chan unaryOpTask, results chan<- taskResult) {
				for task := range tasks {
					newEq, ok := task.UnaryOp.Apply(&task.Equation)
					results <- taskResult{
						Equation: newEq,
						ok:       ok,
					}
				}
			}(tasks, results)
		}

		// put jobs
		go func() {
			for _, op := range s.unaryOps {
				for _, e := range oldKnowledge {
					atomic.AddInt64(&counter, 1)
					tasks <- unaryOpTask{UnaryOp: op, Equation: e}
				}
			}
			close(done)
			close(tasks)
		}()

	WaitUnaryOp:
		// wait for results
		for {
			select {
			case result := <-results:
				if result.ok {
					if oldEq, ok := knowledge[result.value]; !ok || oldEq.cost > result.cost {
						knowledge[result.value] = result.Equation
					}
				}
				atomic.AddInt64(&counter, -1)
			case <-done:
				if counter == 0 {
					break WaitUnaryOp
				}
			}
		}
	}

	s.memo[key] = &knowledge

	fmt.Printf("End calculating %s\n", key)
	return &knowledge
}
