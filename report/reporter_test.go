package report

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/eventlog"
	"testing"
	"time"
)

func TestReporterWithEnable(t *testing.T) {
	initGlobalReporter()
	testReporter(t, true)
}

func TestReporterWithDisable(t *testing.T) {
	initGlobalReporter()
	testReporter(t, false)
}

func testReporter(t *testing.T, enabled bool) {
	user := &auth.HashUser{
		UserNameHash: "698d51a19d8a121c",
		Perm:         "rwx",
	}
	GlobalReporter.Enable(enabled)
	addrOnline := "127.0.0.1:12345"
	addrOffline := "127.0.0.1:54321"
	GlobalReporter.PutConnection(addrOffline)
	GlobalReporter.PutConnection(addrOnline)
	time.Sleep(time.Millisecond * 100)
	GlobalReporter.DeleteConnection(addrOffline)
	GlobalReporter.PutAuth(addrOnline, user)
	GlobalReporter.PutAuth("", user)
	GlobalReporter.PutAuth(addrOnline, nil)
	GlobalReporter.PutEvent(eventlog.NewEvent("./reporter_test.go", "WRITE"))
	GlobalReporter.PutApiStat("127.0.0.1")
	GlobalReporter.PutApiStat("127.0.0.1")
	GlobalReporter.PutApiStat("192.168.1.1")
	time.Sleep(time.Millisecond * 100)

	if enabled {
		testGetReporterWithEnable(t, addrOnline, addrOffline)
	} else {
		testGetReporterWithDisable(t, addrOnline)
	}
}

func testGetReporterWithDisable(t *testing.T, addrOnline string) {

	r := GlobalReporter.GetReport()
	online := r.Online[addrOnline]
	if online != nil {
		t.Errorf("[disabled] test PutConnection [Online] error, expect get a nil connection")
		return
	}

	expectConnCount := 0
	actualConnCount := len(r.Online)
	if actualConnCount != expectConnCount {
		t.Errorf("[disabled] test PutConnection [Online] error, expect get %d connection, actual:%d", expectConnCount, actualConnCount)
	}

	if len(r.Offline) != 0 {
		t.Errorf("[disabled] test PutConnection [Offline] error, expect get 0 connection")
	}

	expectEventCount := 0
	actualEventCount := r.Events.Len()
	if actualEventCount != expectEventCount {
		t.Errorf("[disabled] test PutEvent error, expect get %d event, actual:%d", expectEventCount, actualEventCount)
	}

	expectVisitorStat := 0
	actualVisitorStat := len(r.ApiStat.VisitorStat)
	if expectVisitorStat != actualVisitorStat {
		t.Errorf("[disabled] test PutApiStat error, expect get %d visitor, actual:%d", expectVisitorStat, actualVisitorStat)
	}

	var expectAccessCount uint64
	actualAccessCount := r.ApiStat.AccessCount
	if expectAccessCount != actualAccessCount {
		t.Errorf("[disabled] test PutApiStat error, expect get %d access count, actual:%d", expectAccessCount, actualAccessCount)
	}
}

func testGetReporterWithEnable(t *testing.T, addrOnline, addrOffline string) {
	r := GlobalReporter.GetReport()
	online := r.Online[addrOnline]
	if online == nil {
		t.Errorf("[enabled] test PutConnection [Online] error, get a nil connection")
		return
	}

	expectConnCount := 1
	actualConnCount := len(r.Online)
	if actualConnCount != expectConnCount {
		t.Errorf("[enabled] test PutConnection [Online] error, expect get %d connection, actual:%d", expectConnCount, actualConnCount)
	}

	if len(r.Offline) != 1 && r.Offline[0].Addr == addrOffline {
		t.Errorf("[enabled] test PutConnection [Offline] error,expect get connection %s, actual get a nil connection", addrOffline)
	}

	if !online.IsAuth {
		t.Errorf("[enabled] test PutAuth error, expect get an authorized connection")
	}

	expectEventCount := 1
	actualEventCount := r.Events.Len()
	if actualEventCount != expectEventCount {
		t.Errorf("[enabled] test PutEvent error, expect get %d event, actual:%d", expectEventCount, actualEventCount)
	}

	expectVisitorStat := 2
	actualVisitorStat := len(r.ApiStat.VisitorStat)
	if expectVisitorStat != actualVisitorStat {
		t.Errorf("[enabled] test PutApiStat error, expect get %d visitor, actual:%d", expectVisitorStat, actualVisitorStat)
	}

	var expectAccessCount uint64 = 3
	actualAccessCount := r.ApiStat.AccessCount
	if expectAccessCount != actualAccessCount {
		t.Errorf("[enabled] test PutApiStat error, expect get %d access count, actual:%d", expectAccessCount, actualAccessCount)
	}
}
