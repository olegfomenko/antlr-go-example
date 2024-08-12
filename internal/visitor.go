package internal

import (
	expr "antrl-go-example/internal/antlr"
	"errors"
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"math/big"
)

type Visitor struct {
	*expr.BaseExprVisitor
	storage map[string]*big.Int
}

type ReturnValue struct {
	Value *big.Int
	Error error
}

func NewCalcVisitor() expr.ExprVisitor {
	return &Visitor{storage: make(map[string]*big.Int)}
}

var _ expr.ExprVisitor = &Visitor{}

var ErrorBlank = errors.New("blank Line")

func (v Visitor) Visit(tree antlr.ParseTree) interface{} {
	switch t := tree.(type) {
	case *antlr.ErrorNodeImpl:
		return ReturnValue{nil, fmt.Errorf("syntax error near '%s'", t.GetText())}
	default:
		if cr, ok := tree.Accept(v).(ReturnValue); ok {
			return cr
		}
	}

	return ReturnValue{nil, fmt.Errorf("visit result not of type CalcReturn")}
}

func (v Visitor) VisitChildren(node antlr.RuleNode) interface{} {
	for _, n := range node.GetChildren() {
		cr := v.Visit(n.(antlr.ParseTree)).(ReturnValue)
		if cr.Error != nil {
			if errors.Is(cr.Error, ErrorBlank) {
				continue
			}
			return cr
		}
	}
	return ReturnValue{nil, nil}
}

func (v Visitor) VisitProg(ctx *expr.ProgContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v Visitor) VisitPrintExpr(ctx *expr.PrintExprContext) interface{} {
	res := v.Visit(ctx.Expr()).(ReturnValue)
	if res.Error == nil {
		fmt.Println(res.Value)
	}

	return res
}

func (v Visitor) VisitAssign(ctx *expr.AssignContext) interface{} {
	id := ctx.ID().GetText()
	cr := v.Visit(ctx.Expr()).(ReturnValue)

	if cr.Error != nil {
		return ReturnValue{nil, fmt.Errorf("error with assignment '%s': %w", ctx.GetText(), cr.Error)}
	}

	v.storage[id] = cr.Value
	return cr
}

func (v Visitor) VisitBlank(ctx *expr.BlankContext) interface{} {
	return ReturnValue{nil, ErrorBlank}
}

func (v Visitor) VisitParens(ctx *expr.ParensContext) interface{} {
	return v.Visit(ctx.Expr()).(ReturnValue)
}

func (v Visitor) VisitMulDiv(ctx *expr.MulDivContext) interface{} {
	crLeft := v.Visit(ctx.Expr(0)).(ReturnValue)
	if crLeft.Error != nil {
		return ReturnValue{nil, fmt.Errorf("error with left side visit '%s': %w", ctx.Expr(0).GetText(), crLeft.Error)}
	}
	crRight := v.Visit(ctx.Expr(1)).(ReturnValue)
	if crRight.Error != nil {
		return ReturnValue{nil, fmt.Errorf("error with right side visit '%s': %w", ctx.Expr(1).GetText(), crRight.Error)}
	}

	operator := ctx.GetOp().GetTokenType()

	if operator == expr.ExprLexerMUL {
		return ReturnValue{new(big.Int).Mul(crLeft.Value, crRight.Value), nil}
	}

	if operator == expr.ExprLexerDIV {
		return ReturnValue{new(big.Int).Div(crLeft.Value, crRight.Value), nil}
	}

	return ReturnValue{nil, fmt.Errorf("wrong operator '%v'", operator)}
}

func (v Visitor) VisitAddSub(ctx *expr.AddSubContext) interface{} {
	crLeft := v.Visit(ctx.Expr(0)).(ReturnValue)
	if crLeft.Error != nil {
		return ReturnValue{nil, fmt.Errorf("error with left side visit '%s': %w", ctx.Expr(0).GetText(), crLeft.Error)}
	}
	crRight := v.Visit(ctx.Expr(1)).(ReturnValue)
	if crRight.Error != nil {
		return ReturnValue{nil, fmt.Errorf("error with right side visit '%s': %w", ctx.Expr(1).GetText(), crRight.Error)}
	}

	operator := ctx.GetOp().GetTokenType()

	if operator == expr.ExprLexerADD {
		return ReturnValue{new(big.Int).Add(crLeft.Value, crRight.Value), nil}
	}

	if operator == expr.ExprLexerSUB {
		return ReturnValue{new(big.Int).Sub(crLeft.Value, crRight.Value), nil}
	}

	return ReturnValue{nil, fmt.Errorf("wrong operator '%v'", operator)}
}

func (v Visitor) VisitId(ctx *expr.IdContext) interface{} {
	id := ctx.ID().GetText()
	if value, ok := v.storage[id]; ok {
		return ReturnValue{value, nil}
	}
	return ReturnValue{nil, fmt.Errorf("undefined ID '%s'", id)}
}

func (v Visitor) VisitInt(ctx *expr.IntContext) interface{} {
	i, ok := new(big.Int).SetString(ctx.GetText(), 10)
	if !ok {
		return ReturnValue{nil, fmt.Errorf("couldn't parse integer: '%s'", ctx.GetText())}
	}
	return ReturnValue{i, nil}
}
