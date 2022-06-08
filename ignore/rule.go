package ignore

import "strings"

// Rule the match rule provider
type Rule interface {
	// Match reports whether the string s is matched of this rule
	Match(s string) bool
	// SwitchName return the rule switch name
	SwitchName() string
	// Expression return the rule expression
	Expression() string
}

const (
	filePathSwitch = "[filepath]"
	regexpSwitch   = "[regexp]"
)

func newRule(expr string, switchName string) (Rule, error) {
	switchName = strings.TrimSpace(switchName)
	if switchName == regexpSwitch {
		return newRegexpRule(expr)
	}
	return newFilePathRule(expr)
}
