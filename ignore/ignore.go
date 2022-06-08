package ignore

import (
	"os"
	"strings"

	"github.com/no-src/log"
)

// Ignore support to check the string matches the ignore rule or not
type Ignore interface {
	Match(s string) bool
}

type ignore struct {
	rules []Rule
}

// New get a default Ignore instance
func New(ignoreFile string) (Ignore, error) {
	rules, err := parseIgnoreFile(ignoreFile)
	if err != nil {
		return nil, err
	}
	return &ignore{
		rules: rules,
	}, nil
}

func (ig *ignore) Match(s string) bool {
	for _, rule := range ig.rules {
		if rule.Match(s) {
			return true
		}
	}
	return false
}

func parseIgnoreFile(ignoreFile string) ([]Rule, error) {
	conf, err := os.ReadFile(ignoreFile)
	if err != nil {
		return nil, err
	}
	return parse(conf)
}

func parse(data []byte) (rs []Rule, err error) {
	conf := string(data)
	lines := strings.Split(conf, "\n")
	switchName := filePathSwitch
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 && line[0] != '#' {
			if line == filePathSwitch {
				switchName = filePathSwitch
			} else if line == regexpSwitch {
				switchName = regexpSwitch
			} else {
				r, err := newRule(line, switchName)
				if err != nil {
					return nil, err
				}
				log.Debug("register %s rule, expression=%s", r.SwitchName(), r.Expression())
				rs = append(rs, r)
			}
		}
	}
	return rs, nil
}
