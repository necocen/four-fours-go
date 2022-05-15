package fourfours

import (
	"fmt"
)

const parenLeft string = "("
const parenRight string = ")"

func Print(unaryOps []UnaryOpPrinter, binaryOps []BinaryOpPrinter, e Equation) string {
	// 部分式文字列と、最後に追加された演算子の優先度のペアのスタック
	stack := []struct {
		string
		precedence int
	}{}

T:
	for _, token := range e.tokens {
		if token >= 0xf0 {
			// 数値2桁目以降
			stack[len(stack)-1].string = fmt.Sprintf("%s%d", stack[len(stack)-1].string, uint(token)-0xf0)
			stack[len(stack)-1].precedence = 0
		} else if token >= 0xe0 {
			// 数値1桁目
			stack = append(stack, struct {
				string
				precedence int
			}{
				string:     fmt.Sprintf("%d", uint8(token)-0xe0),
				precedence: 0,
			})
		} else {
			for _, op := range unaryOps {
				if op.token != token {
					continue
				}

				pop := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				var expr string
				if op.paren && pop.precedence >= op.precedence {
					expr = parenLeft + pop.string + parenRight
				} else {
					expr = pop.string
				}

				stack = append(stack, struct {
					string
					precedence int
				}{
					string:     op.prefix + expr + op.suffix,
					precedence: op.precedence,
				})

				continue T
			}
			for _, op := range binaryOps {
				if op.token != token {
					continue
				}

				rhs := stack[len(stack)-1]
				lhs := stack[len(stack)-2]
				stack = stack[:len(stack)-2]

				var expr1, expr2 string
				if op.parenLeft && (lhs.precedence > op.precedence || lhs.precedence == op.precedence && !op.leftAssociative) {
					expr1 = parenLeft + lhs.string + parenRight
				} else {
					expr1 = lhs.string
				}
				if op.parenRight && (rhs.precedence > op.precedence || rhs.precedence == op.precedence && !op.rightAssociative) {
					expr2 = parenLeft + rhs.string + parenRight
				} else {
					expr2 = rhs.string
				}

				stack = append(stack, struct {
					string
					precedence int
				}{
					string:     op.prefix + expr1 + op.infix + expr2 + op.suffix,
					precedence: op.precedence,
				})

				continue T
			}

			panic("Unexpected token")
		}
	}

	if len(stack) != 1 {
		panic("Invalid equation")
	}

	return stack[0].string
}
