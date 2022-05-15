package fourfours

import "fmt"

type Searcher struct {
	memo      map[string]map[float64]*Equation
	unaryOps  []UnaryOp
	binaryOps []BinaryOp
}

func NewSearcher(unaryOps []UnaryOp, binaryOps []BinaryOp) Searcher {
	return Searcher{
		memo:      make(map[string]map[float64]*Equation),
		unaryOps:  unaryOps,
		binaryOps: binaryOps,
	}
}

func (s *Searcher) Search(key string) map[float64]*Equation {
	if knowledge, ok := s.memo[key]; ok {
		return knowledge
	}

	fmt.Printf("Start calculating %s\n", key)

	knowledge := make(map[float64]*Equation)
	e := NewEquationWithNumber(key)
	if e == nil {
		panic("Invalid number")
	}
	knowledge[e.value] = e

	for i := 1; i < len(key); i++ {
		keyLeft := key[:i]
		keyRight := key[i:]
		fmt.Printf("Combining %s and %s\n", keyLeft, keyRight)
		knowledgeLeft := s.Search(keyLeft)
		knowledgeRight := s.Search(keyRight)
		for _, op := range s.binaryOps {
			for _, el := range knowledgeLeft {
				for _, er := range knowledgeRight {
					newEq := op.Apply(el, er)
					if newEq == nil {
						continue
					}
					if oldEq, ok := knowledge[newEq.value]; !ok || oldEq.cost > newEq.cost {
						knowledge[newEq.value] = newEq
					}
				}
			}
		}
	}

	for i := 0; i < 3; i++ {
		fmt.Printf("Applying Unary Ops to %s (%d/3)\n", key, i+1)
		for _, op := range s.unaryOps {
			for _, e := range knowledge {
				newEq := op.Apply(e)
				if newEq == nil {
					continue
				}
				if oldEq, ok := knowledge[newEq.value]; !ok || oldEq.cost > newEq.cost {
					knowledge[newEq.value] = newEq
				}
			}
		}
	}

	s.memo[key] = knowledge

	fmt.Printf("End calculating %s\n", key)
	return knowledge
}
