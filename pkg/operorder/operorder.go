package operorder

var precedence map[string]Operator

func init() {
	precedence = make(map[string]Operator)
	precedence["**"] = Operator{Value: "**", Precedence: 1, Associativity: "right"}
	precedence["^"] = Operator{Value: "^", Precedence: 1, Associativity: "right"}
	precedence["*"] = Operator{Value: "*", Precedence: 2, Associativity: "left"}
	precedence["/"] = Operator{Value: "/", Precedence: 2, Associativity: "left"}
	precedence["+"] = Operator{Value: "+", Precedence: 3, Associativity: "left"}
	precedence["-"] = Operator{Value: "-", Precedence: 3, Associativity: "left"}
}

type Comparator interface {
	Compare(o Operator) int
}

type Operator struct {
	Value         string
	Precedence    int
	Associativity string
}

func New(op string) Operator {
	return precedence[op]
}

func compare(o1, o2 Operator) int {
	if o1.Precedence < o2.Precedence {
		return 1
	} else if o1.Precedence > o2.Precedence {
		return -1
	}
	return 0
}

func (op Operator) Compare(o Operator) int {
	o1 := precedence[op.Value]
	o2 := precedence[o.Value]
	return compare(o1, o2)
}
