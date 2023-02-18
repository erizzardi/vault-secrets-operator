package utils

import (
	"os"
	"testing"
)

func TestGetEnvThatExists(t *testing.T) {
	envName := "TESTENV"
	envValue := "TESTVALUE"
	os.Setenv(envName, envValue)
	res := GetEnvOrFallback(envName, "")

	if res != envValue {
		t.Errorf("%s != %s", res, envValue)
	}
}

func TestGetEnvThatDoesNotExist(t *testing.T) {
	envName := "TESTENV"
	fallback := "FALLBACK"
	res := GetEnvOrFallback(envName, fallback)
	if res != fallback {
		t.Errorf("%s != %s", res, fallback)
	}
}

func TestFlagToEnv(t *testing.T) {
	in := "as5d&f-asd*f_asd^f"
	expected := "AS5D&F_ASD*F_ASD^F"
	res := FlagToEnv(in)
	if res != expected {
		t.Errorf("%s != %s", res, expected)
	}
}
