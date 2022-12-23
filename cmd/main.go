// Copyright (c) 2022, chuckldev.  All rights reserved

package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"calculator/pkg/operorder"

	exprtree "github.com/chuckldev/goexprtree"
	stack "github.com/chuckldev/gostack"
)

// operator map
var opm map[string]func(float64, float64) float64
var rxOpenParens = regexp.MustCompile(`\(`)
var rxCloseParens = regexp.MustCompile(`\)`)
var rxNum = regexp.MustCompile(`(([0-9]*[.])?[0-9]+)|[0-9]`)
var rxPerCom = regexp.MustCompile(`[,.]`)
var rxOper = regexp.MustCompile(`[\*\/\+\-\^]`)

func init() {
	opm = make(map[string]func(float64, float64) float64)
	opm[`^`] = Raise
	opm[`*`] = Mult
	opm[`/`] = Div
	opm[`+`] = Add
	opm[`-`] = Subt
}

func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

func IsOpenParen(s string) bool {
	return rxOpenParens.MatchString(s)
}

func IsCloseParen(s string) bool {
	return rxCloseParens.MatchString(s)
}

func IsNum(s string) bool {
	return rxNum.MatchString(s)
}

func IsPeriodComma(s string) bool {
	return rxPerCom.MatchString(s)
}

func IsOperator(s string) bool {
	return rxOper.MatchString(s)
}

func SplitString(expr string) []string {
	s := strings.Split(strings.ReplaceAll(expr, " ", ""), "")
	var res []string
	for i := 0; i < len(s); i++ {
		if !IsNum(s[i]) && !IsPeriodComma(s[i]) {
			res = append(res, s[i])
			continue
		} else {
			c := ""
			for (i < len(s) && IsNum(s[i])) || (i < len(s) && IsPeriodComma(s[i])) {
				if s[i] == "," {
					i++
					continue
				}
				c += s[i]
				i++
			}
			res = append(res, c)
			if i < len(s) {
				res = append(res, s[i])
			}
		}
	}
	return res
}

func Add(op1, op2 float64) float64 {
	return op1 + op2
}

func Subt(op1, op2 float64) float64 {
	return op1 - op2
}

func Mult(op1, op2 float64) float64 {
	return op1 * op2
}

func Div(op1, op2 float64) float64 {
	var res float64
	if op2 == 0 {
		panic("Cannot divide by Zero ( 0 )")
	}
	res = op1 / op2
	return res
}

func Raise(op1, op2 float64) float64 {
	return math.Pow(op1, op2)
}

func Calculate(t *exprtree.ExprTree) float64 {
	if t == nil {
		return 0.0
	}

	if t.Left == nil && t.Right == nil {
		res, _ := strconv.ParseFloat(t.Value, 64)
		return res
	}

	left := Calculate(t.Left)
	right := Calculate(t.Right)

	return opm[t.Value](left, right)
}

func ToString(t *exprtree.ExprTree) string {
	if t == nil {
		return ""
	}

	ToString(t.Left)
	fmt.Print(t.Value)
	ToString(t.Right)

	return ""
}

func BuildExprTree(s *stack.Stack) *exprtree.ExprTree {
	exprStack := stack.New()
	for !s.IsEmpty() {
		val := s.Pop().(string)
		if IsNum(val) {
			t := exprtree.New(val)
			exprStack.Push(t)
		}

		if IsOperator(val) {
			t1 := exprStack.Pop().(*exprtree.ExprTree)
			t2 := exprStack.Pop().(*exprtree.ExprTree)
			t := &exprtree.ExprTree{Value: val, Left: t2, Right: t1}
			exprStack.Push(t)
		}
	}

	temp := exprStack.Pop().(*exprtree.ExprTree)
	ToString(temp)
	fmt.Println()
	return temp
}

func EvaluateExpression(chars []string) float64 {
	opStack := stack.Stack{}
	outStack := stack.Stack{}

	for _, v := range chars {
		if IsNum(v) {
			outStack.Push(v)
			continue
		}

		if IsOperator(v) {
			newOper := operorder.New(v)
			if !opStack.IsEmpty() {
				topOper := opStack.Peek().(*operorder.Operator)
				for {
					if opStack.IsEmpty() {
						break
					}

					if topOper.Value == "(" {
						break
					}

					precValue := topOper.Compare(newOper)

					if precValue > 0 || (precValue == 0 && newOper.Associativity == "left") {
						topOper = opStack.Pop().(*operorder.Operator)
						outStack.Push(topOper.Value)
						break
					}
					break
				}
			}
			opStack.Push(&newOper)
			continue
		}

		if IsOpenParen(v) {
			newOper := operorder.Operator{Value: "("}
			opStack.Push(&newOper)
			continue
		}

		if IsCloseParen(v) {
			topOper := opStack.Peek().(*operorder.Operator)
			for !IsOpenParen(topOper.Value) {
				if opStack.IsEmpty() {
					break
				}
				oper := opStack.Pop().(*operorder.Operator)
				outStack.Push(oper.Value)
				topOper = opStack.Peek().(*operorder.Operator)
			}
			if IsOpenParen(topOper.Value) {
				opStack.Pop()
			}
			continue
		}
	}

	for !opStack.IsEmpty() {
		oper := opStack.Pop().(*operorder.Operator)
		outStack.Push(oper.Value)
	}

	revStack := outStack.Reverse()
	exprTree := BuildExprTree(&revStack)

	res := Calculate(exprTree)
	return res
}

func main() {
	if len(os.Args) < 1 {
		panic("Calculator requires at least one argument ( string ).")
	}

	// the problem inputted from cmd line
	expr := os.Args[1]

	chars := SplitString(expr)
	result := EvaluateExpression(chars)
	fmt.Println(result)
}
