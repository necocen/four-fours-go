package fourfours

type BinaryOp struct {
	cost  uint8
	token OperatorToken
	apply func(lhs, rhs float64) (float64, bool)
}

func NewBinaryOp(token OperatorToken, cost uint8, apply func(lhs, rhs float64) (float64, bool)) BinaryOp {
	return BinaryOp{
		cost:  cost,
		token: token,
		apply: apply,
	}
}

func (op *BinaryOp) Apply(lhs, rhs *Equation) *Equation {
	res, ok := op.apply(lhs.value, rhs.value)
	if !ok {
		return nil
	}

	tokens := make([]OperatorToken, len(lhs.tokens)+len(rhs.tokens)+1)
	copy(tokens, lhs.tokens)
	copy(tokens[len(lhs.tokens):], rhs.tokens)
	tokens[len(lhs.tokens)+len(rhs.tokens)] = op.token

	return &Equation{
		tokens: tokens,
		cost:   lhs.cost + rhs.cost + op.cost,
		value:  res,
	}
}

type BinaryOpPrinter struct {
	token            OperatorToken
	prefix           string
	infix            string
	suffix           string
	precedence       int
	leftAssociative  bool
	rightAssociative bool
	parenLeft        bool
	parenRight       bool
}

func NewBinaryOpPrinter(token OperatorToken, prefix, infix, suffix string, precedence int, leftAssociative, rightAssociative, parenLeft, parenRight bool) BinaryOpPrinter {
	return BinaryOpPrinter{
		token:            token,
		prefix:           prefix,
		infix:            infix,
		suffix:           suffix,
		precedence:       precedence,
		leftAssociative:  leftAssociative,
		rightAssociative: rightAssociative,
		parenLeft:        parenLeft,
		parenRight:       parenRight,
	}
}
