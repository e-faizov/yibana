package internal

import (
	"testing"
)

func TestHash(t *testing.T) {
	hash := calcHash("teststr_||13123", "testkey")
	if hash != "a2dd0b714f9d42f16073de346a7780fdbc56d3db4828611e31f9416a90070929" {
		t.Error("wrong hash")
	}

	hash = calcHash("teststr_||13123", "")
	if hash != "f49245637e8030bc41855623913fe971ab2bf17634acbc75198e8fb7cc447dd1" {
		t.Error("wrong hash")
	}
}
