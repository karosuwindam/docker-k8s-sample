package webserver

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	wantPort := 8080
	t.Setenv("PORT", fmt.Sprintf("%d", wantPort))

	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}
	if got.Port != strconv.Itoa(wantPort) {
		t.Errorf("want %v, but %s", wantPort, got.Port)
	}
	wantProtcal := "tcp"
	if got.Protocol != wantProtcal {
		t.Errorf("want %s, but %s", wantProtcal, got.Protocol)
	}
	t.Logf("Port %v, Protocol %v, IP %v", got.Port, got.Protocol, got.IP)

}
