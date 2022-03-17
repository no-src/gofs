package ignore

import (
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/util/stringutil"
	"github.com/no-src/log"
)

var defaultIgnore Ignore
var matchIgnoreDeletedPath bool

// Init init default ignore component
// ignoreConf the config file path of the ignore component
// ignoreDeletedPath whether ignore the deleted path
func Init(ignoreConf string, ignoreDeletedPath bool) error {
	if !stringutil.IsEmpty(ignoreConf) {
		ig, err := New(ignoreConf)
		if err != nil {
			return err
		}
		defaultIgnore = ig
	}
	matchIgnoreDeletedPath = ignoreDeletedPath
	return nil
}

// Match the current string matches the rule or not
func Match(s string) bool {
	if defaultIgnore != nil {
		return defaultIgnore.Match(s)
	}
	return false
}

// MatchPath the current string matches the rule or not, if enable the matchIgnoreDeletedPath, check the deleted file rule is matched or not first
func MatchPath(path, caller, desc string) bool {
	var matched bool
	if matchIgnoreDeletedPath {
		matched = fs.IsDeleted(path)
		if matched {
			log.Debug("[ignored] [%s] a deleted path is matched [%s] => [%s]", caller, desc, path)
			return true
		}
	}

	matched = Match(path)
	if matched {
		log.Debug("[ignored] [%s] an ignore path is matched [%s] => [%s]", caller, desc, path)
	}
	return matched
}
