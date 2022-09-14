package main

import (
	"os"
	"testing"
)

func TestGetEnvVar(t *testing.T) {
	os.Setenv("TEST", "test")

	if GetEnvVar("TEST", "fallback") != "test" {
		t.Error("GetEnvVar() failed")
	}

	if GetEnvVar("TEST2", "fallback") != "fallback" {
		t.Error("GetEnvVar() failed")
	}
}

func TestGetBucketAndPrefix(t *testing.T) {
	path := "test/in1"
	bucket, prefix := GetBucketAndPrefix(path)

	if bucket != "test" && prefix != "in1" {
		t.Error("Bucket and prefix not correctly parsed")
	}

	path = "test/in1/in2"
	bucket, prefix = GetBucketAndPrefix(path)

	if bucket != "test" && prefix != "in1/in2" {
		t.Error("Bucket and prefix not correctly parsed")
	}
}
