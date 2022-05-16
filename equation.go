package fourfours

import "strconv"

type Equation struct {
	tokens []OperatorToken
	cost   uint8
	value  float64
}

func (e *Equation) IsBetterThan(other *Equation) bool {
	return e.cost < other.cost
}

func NewEquationWithNumber(numbers string) (Equation, bool) {
	num, err := strconv.Atoi(numbers)
	if err != nil {
		return Equation{}, false
	}

	tokens := make([]OperatorToken, len(numbers))
	for i, c := range []byte(numbers) {
		if i == 0 {
			tokens[i] = OperatorToken(0xe0 + c - '0')
		} else {
			tokens[i] = OperatorToken(0xf0 + c - '0')
		}
	}

	return Equation{
		tokens: tokens,
		cost:   0,
		value:  float64(num),
	}, true
}
