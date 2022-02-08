package ignore

import (
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
)

var defaultIgnore Ignore
var matchIgnoreDeletedPath bool

// Init init default ignore component
// ignoreConf the config file path of the ignore component
// ignoreDeletedPath whether ignore the deleted path
func Init(ignoreConf string, ignoreDeletedPath bool) error {
	if !util.IsEmpty(ignoreConf) {
		ig, err := New(ignoreConf)
		if err != nil {
			return err
		}
		defaultIgnore = ig
	}
	matchIgnoreDeletedPath = ignoreDeletedPath
	return nil
}

func Match(s string) bool {
	if defaultIgnore != nil {
		return defaultIgnore.Match(s)
	}
	return false
}

func MatchPath(path, caller, desc string) bool {
	matched := false
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
