package ignore

import (
	"os"
	"strings"
)

type Ignore interface {
	Match(s string) bool
}

type ignore struct {
	rules []*rule
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

func parseIgnoreFile(ignoreFile string) ([]*rule, error) {
	conf, err := os.ReadFile(ignoreFile)
	if err != nil {
		return nil, err
	}
	return parse(conf)
}

func parse(data []byte) (rs []*rule, err error) {
	conf := string(data)
	lines := strings.Split(conf, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 && line[0] != '#' {
			r, err := newRule(line)
			if err != nil {
				return nil, err
			}
			rs = append(rs, r)
		}
	}
	return rs, nil
}
