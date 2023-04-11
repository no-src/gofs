package report

import (
	"sync"
	"testing"
	"time"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/eventlog"
)

func TestReporter_WithEnable(t *testing.T) {
	testReporter(t, true)
}

func TestReporter_WithDisable(t *testing.T) {
	testReporter(t, false)
}

func TestReporter_WithEnable_Concurrent(t *testing.T) {
	testReporterConcurrent(t, true)
}

func TestReporter_WithDisable_Concurrent(t *testing.T) {
	testReporterConcurrent(t, false)
}

func getTestReporter(enabled bool) (reporter Reporter, addrOnline, addrOffline string) {
	user := &auth.SessionUser{
		UserName: "698d51a19d8a121c",
		Perm:     "rwx",
	}
	reporter = NewReporter()
	reporter.Enable(enabled)
	addrOnline = "127.0.0.1:12345"
	addrOffline = "127.0.0.1:54321"
	reporter.PutConnection(addrOffline, user)
	reporter.PutConnection(addrOnline, user)
	time.Sleep(time.Millisecond * 100)
	reporter.DeleteConnection(addrOffline)
	reporter.PutConnection(addrOnline, nil)
	reporter.PutEvent(eventlog.NewEvent("./reporter_test.go", "WRITE"))
	reporter.PutApiStat("127.0.0.1")
	reporter.PutApiStat("127.0.0.1")
	reporter.PutApiStat("192.168.1.1")
	time.Sleep(time.Millisecond * 100)
	return
}

func testReporterConcurrent(t *testing.T, enabled bool) {
	reporter, addrOnline, addrOffline := getTestReporter(enabled)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			if enabled {
				testGetReporterWithEnable(t, reporter, addrOnline, addrOffline)
			} else {
				testGetReporterWithDisable(t, reporter, addrOnline)
			}
			reporter.Enable(enabled)
			wg.Done()
		}()
	}
	wg.Wait()
}

func testReporter(t *testing.T, enabled bool) {
	reporter, addrOnline, addrOffline := getTestReporter(enabled)
	if enabled {
		testGetReporterWithEnable(t, reporter, addrOnline, addrOffline)
	} else {
		testGetReporterWithDisable(t, reporter, addrOnline)
	}
}

func testGetReporterWithDisable(t *testing.T, reporter Reporter, addrOnline string) {

	r := reporter.GetReport()
	online := r.Online[addrOnline]
	if online != nil {
		t.Errorf("[disabled] test PutConnection [Online] error, expect to get a nil connection")
		return
	}

	expectConnCount := 0
	actualConnCount := len(r.Online)
	if actualConnCount != expectConnCount {
		t.Errorf("[disabled] test PutConnection [Online] error, expect to get %d connection, actual:%d", expectConnCount, actualConnCount)
	}

	if len(r.Offline) != 0 {
		t.Errorf("[disabled] test PutConnection [Offline] error, expect to get 0 connection")
	}

	expectEventCount := 0
	actualEventCount := r.Events.Len()
	if actualEventCount != expectEventCount {
		t.Errorf("[disabled] test PutEvent error, expect to get %d event, actual:%d", expectEventCount, actualEventCount)
	}

	expectVisitorStat := 0
	actualVisitorStat := len(r.ApiStat.VisitorStat)
	if expectVisitorStat != actualVisitorStat {
		t.Errorf("[disabled] test PutApiStat error, expect to get %d visitor, actual:%d", expectVisitorStat, actualVisitorStat)
	}

	var expectAccessCount uint64
	actualAccessCount := r.ApiStat.AccessCount
	if expectAccessCount != actualAccessCount {
		t.Errorf("[disabled] test PutApiStat error, expect to get %d access count, actual:%d", expectAccessCount, actualAccessCount)
	}
}

func testGetReporterWithEnable(t *testing.T, reporter Reporter, addrOnline, addrOffline string) {
	r := reporter.GetReport()
	online := r.Online[addrOnline]
	if online == nil {
		t.Errorf("[enabled] test PutConnection [Online] error, get a nil connection")
		return
	}

	expectConnCount := 1
	actualConnCount := len(r.Online)
	if actualConnCount != expectConnCount {
		t.Errorf("[enabled] test PutConnection [Online] error, expect to get %d connection, actual:%d", expectConnCount, actualConnCount)
	}

	if len(r.Offline) != 1 && r.Offline[0].Addr == addrOffline {
		t.Errorf("[enabled] test PutConnection [Offline] error,expect to get connection %s, actual get a nil connection", addrOffline)
	}

	expectEventCount := 1
	actualEventCount := r.Events.Len()
	if actualEventCount != expectEventCount {
		t.Errorf("[enabled] test PutEvent error, expect to get %d event, actual:%d", expectEventCount, actualEventCount)
	}

	expectVisitorStat := 2
	actualVisitorStat := len(r.ApiStat.VisitorStat)
	if expectVisitorStat != actualVisitorStat {
		t.Errorf("[enabled] test PutApiStat error, expect to get %d visitor, actual:%d", expectVisitorStat, actualVisitorStat)
	}

	var expectAccessCount uint64 = 3
	actualAccessCount := r.ApiStat.AccessCount
	if expectAccessCount != actualAccessCount {
		t.Errorf("[enabled] test PutApiStat error, expect to get %d access count, actual:%d", expectAccessCount, actualAccessCount)
	}
}
