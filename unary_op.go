package fourfours

type UnaryOp struct {
	cost  uint8
	token OperatorToken
	apply func(v float64) (float64, bool)
}

func NewUnaryOp(token OperatorToken, cost uint8, apply func(float64) (float64, bool)) UnaryOp {
	return UnaryOp{
		cost:  cost,
		token: token,
		apply: apply,
	}
}

func (op *UnaryOp) Apply(eq *Equation) *Equation {
	res, ok := op.apply(eq.value)
	if !ok {
		return nil
	}

	tokens := make([]OperatorToken, len(eq.tokens)+1)
	copy(tokens, eq.tokens)
	tokens[len(eq.tokens)] = op.token

	return &Equation{
		tokens: tokens,
		cost:   eq.cost + op.cost,
		value:  res,
	}
}

type UnaryOpPrinter struct {
	token      OperatorToken
	prefix     string
	suffix     string
	precedence int
	paren      bool
}

func NewUnaryOpPrinter(token OperatorToken, prefix, suffix string, precedence int, paren bool) UnaryOpPrinter {
	return UnaryOpPrinter{
		token:      token,
		prefix:     prefix,
		suffix:     suffix,
		precedence: precedence,
		paren:      paren,
	}
}
