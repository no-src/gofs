package ignore

import (
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/logger"
	"github.com/no-src/gofs/util/stringutil"
)

// PathIgnore check the ignore rules of the specified file path
type PathIgnore interface {
	// MatchPath the current string matches the rule or not, if enable the matchIgnoreDeletedPath, check the deleted file rule is matched or not first
	MatchPath(path, caller, desc string) bool
}

type pathIgnore struct {
	ig                Ignore
	ignoreDeletedPath bool
	logger            *logger.Logger
}

// NewPathIgnore create an instance of the PathIgnore component
// ignoreConf the config file path of the ignore component
// ignoreDeletedPath whether ignore the deleted path
func NewPathIgnore(ignoreConf string, ignoreDeletedPath bool, logger *logger.Logger) (PathIgnore, error) {
	pi := &pathIgnore{
		logger: logger,
	}
	if !stringutil.IsEmpty(ignoreConf) {
		ig, err := New(ignoreConf, logger)
		if err != nil {
			return nil, err
		}
		pi.ig = ig
	}
	pi.ignoreDeletedPath = ignoreDeletedPath
	return pi, nil
}

// match the current string matches the rule or not
func (pi *pathIgnore) match(s string) bool {
	if pi.ig != nil {
		return pi.ig.Match(s)
	}
	return false
}

func (pi *pathIgnore) MatchPath(path, caller, desc string) bool {
	var matched bool
	if pi.ignoreDeletedPath {
		matched = fs.IsDeleted(path)
		if matched {
			pi.logger.Debug("[ignored] [%s] a deleted path is matched [%s] => [%s]", caller, desc, path)
			return true
		}
	}
	matched = pi.match(path)
	if matched {
		pi.logger.Debug("[ignored] [%s] an ignore path is matched [%s] => [%s]", caller, desc, path)
	}
	return matched
}
