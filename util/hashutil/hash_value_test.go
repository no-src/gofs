package hashutil

import "testing"

func TestHashValues(t *testing.T) {
	var hvs HashValues
	expect := "461d19e03559ff8a1284951bab8327e1"
	if hvs.Last() != nil {
		t.Errorf("test TestHashValues.Last error, expect get nil")
	}
	hvs = append(hvs, NewHashValue(1, "21cc28409729565fc1a4d2dd92db269f"))
	hvs = append(hvs, NewHashValue(2, expect))

	if hvs.Last() == nil {
		t.Errorf("test TestHashValues.Last error, expect:%s, actual get nil", expect)
		return
	}

	actual := hvs.Last().Hash
	if actual != expect {
		t.Errorf("test TestHashValues.Last error, expect:%s, actual:%s", expect, actual)
	}
}
