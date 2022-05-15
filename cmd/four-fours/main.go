package main

import (
	"fmt"
	"fourfours"
	"math"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()
	negate := fourfours.NewUnaryOp(fourfours.OperatorToken(0x00), 2, func(v float64) (float64, bool) { return -v, true })
	sqrt := fourfours.NewUnaryOp(0x01, 4, func(v float64) (float64, bool) {
		if v > 0 && v != 1 {
			return math.Sqrt(v), true
		} else {
			return 0, false
		}
	})
	fact := fourfours.NewUnaryOp(0x02, 6, func(v float64) (float64, bool) {
		table := [9]float64{1, 1, 2, 6, 24, 120, 720, 5040, 40320}
		integer, frac := math.Modf(v)
		if integer < 0 || int(integer) >= len(table) || frac != 0 {
			return 0, false
		} else {
			return float64(table[int(integer)]), true
		}
	})
	add := fourfours.NewBinaryOp(0x10, 1, func(lhs, rhs float64) (float64, bool) { return lhs + rhs, true })
	sub := fourfours.NewBinaryOp(0x11, 2, func(lhs, rhs float64) (float64, bool) { return lhs - rhs, true })
	mul := fourfours.NewBinaryOp(0x12, 3, func(lhs, rhs float64) (float64, bool) { return lhs * rhs, true })
	div := fourfours.NewBinaryOp(0x13, 4, func(lhs, rhs float64) (float64, bool) {
		if rhs == 0 {
			return 0, false
		} else {
			return lhs / rhs, true
		}
	})
	pow := fourfours.NewBinaryOp(0x14, 6, func(lhs, rhs float64) (float64, bool) { return math.Pow(lhs, rhs), true })

	searcher := fourfours.NewSearcher([]fourfours.UnaryOp{negate, sqrt, fact}, []fourfours.BinaryOp{add, sub, mul, div, pow})

	knowledge := searcher.Search("4444")

	negateP := fourfours.NewUnaryOpPrinter(0x00, "-", "", 3, true)
	sqrtP := fourfours.NewUnaryOpPrinter(0x01, "√", "", 1, true)
	factP := fourfours.NewUnaryOpPrinter(0x02, "", "!", 2, true)
	addP := fourfours.NewBinaryOpPrinter(0x10, "", "+", "", 6, true, true, true, true)
	subP := fourfours.NewBinaryOpPrinter(0x11, "", "-", "", 6, true, false, true, true)
	mulP := fourfours.NewBinaryOpPrinter(0x12, "", "*", "", 5, true, true, true, true)
	divP := fourfours.NewBinaryOpPrinter(0x13, "", "/", "", 5, true, false, true, true)
	powP := fourfours.NewBinaryOpPrinter(0x14, "", "^", "", 3, false, true, true, true)

	result := make(map[int]*fourfours.Equation)
	for v, eq := range knowledge {
		_, frac := math.Modf(v)
		if frac != 0 || v < 0 || v > 1000 {
			continue
		}
		result[int(v)] = eq
	}

	for i := 0; i <= 1000; i++ {
		eq, ok := result[i]
		if !ok {
			continue
		}
		expr := fourfours.Print([]fourfours.UnaryOpPrinter{negateP, sqrtP, factP}, []fourfours.BinaryOpPrinter{addP, subP, mulP, divP, powP}, eq)
		fmt.Printf("%d = %s\n", i, expr)
	}
}
