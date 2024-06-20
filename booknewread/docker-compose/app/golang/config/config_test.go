package config

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	os.Setenv("TRACER_ON", "true")
	Init()
	if !TraData.TracerUse {
		t.Fail()
	}
	os.Setenv("TRACER_ON", "false")
	Init()
	if TraData.TracerUse {
		t.Fail()
	}
	a, b := TracerS(context.TODO(), "aa", "bb")
	fmt.Println(a, b)
}
