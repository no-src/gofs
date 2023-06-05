package integration

import (
	"os"
	"testing"
	"time"

	"github.com/no-src/fsctl/command"
	"github.com/no-src/gofs/cmd"
	"github.com/no-src/gofs/result"
)

func getRunConf(conf string) string {
	return "./testdata/conf/" + conf
}

func getTestConf(conf string) string {
	return "./testdata/test/" + conf
}

func runWithConfigFile(path string) result.Result {
	return cmd.RunWithConfigFile(path)
}

func testIntegrationClientServer(t *testing.T, runServerConf string, runClientConf string, testConf string) {
	if len(runServerConf) > 0 {
		runServerConf = getRunConf(runServerConf)
	}
	runClientConf = getRunConf(runClientConf)
	testConf = getTestConf(testConf)

	commands, err := command.ParseConfigFile(testConf)
	if err != nil {
		t.Errorf("parse confile file error, err=%v", err)
		return
	}

	if err = commands.ExecInit(); err != nil {
		t.Errorf("execute init commands error, err=%v", err)
		return
	}

	var sr result.Result
	if len(runServerConf) == 0 {
		sr = noServer()
	} else {
		sr = runWithConfigFile(runServerConf)
	}
	if err = sr.WaitInit(); err != nil {
		t.Errorf("wait gofs server init error, err=%v", err)
		return
	}

	cr := runWithConfigFile(runClientConf)
	if err = cr.WaitInit(); err != nil {
		t.Errorf("wait gofs client init error, err=%v", err)
		// shutdown the server
		if err = sr.Shutdown(); err != nil {
			t.Errorf("gofs server shutdown error, %v", err)
		}
		return
	}

	time.Sleep(time.Second)

	if err = commands.ExecActions(); err != nil {
		t.Errorf("execute actions commands error, err=%v", err)
	}

	if err = cr.Shutdown(); err != nil {
		t.Errorf("gofs client shutdown error, %v", err)
	}

	if err = sr.Shutdown(); err != nil {
		t.Errorf("gofs server shutdown error, %v", err)
	}

	if err = cr.Wait(); err != nil {
		t.Errorf("wait for the gofs client exit error, %v", err)
	}

	if err = sr.Wait(); err != nil {
		t.Errorf("wait for the gofs server exit error, %v", err)
	}

	if err = commands.ExecClear(); err != nil {
		t.Errorf("execute clear commands error, err=%v", err)
	}
}

func noServer() result.Result {
	r := result.New()
	r.InitDone()
	r.RegisterNotifyHandler(func(s os.Signal, timeout ...time.Duration) error {
		r.Done()
		return nil
	})
	return r
}
