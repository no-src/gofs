package ignore

import (
	"fmt"
	"regexp"
)

type regexpRule struct {
	expr string
	reg  *regexp.Regexp
}

func newRegexpRule(expr string) (Rule, error) {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("parse %s rule failed, expression=%s, %w", regexpSwitch, expr, err)
	}
	return &regexpRule{
		expr: expr,
		reg:  reg,
	}, nil
}

func (r *regexpRule) Match(s string) bool {
	return r.reg.MatchString(s)
}

func (r *regexpRule) SwitchName() string {
	return regexpSwitch
}

func (r *regexpRule) Expression() string {
	return r.expr
}
