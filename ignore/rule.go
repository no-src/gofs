package ignore

import (
	"regexp"
)

type rule struct {
	reg *regexp.Regexp
}

func newRule(expr string) (*rule, error) {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &rule{
		reg: reg,
	}, nil
}

func (r *rule) Match(s string) bool {
	return r.reg.MatchString(s)
}
