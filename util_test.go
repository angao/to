package main

import (
	"testing"
)

func TestOctetStream(t *testing.T) {
	if OctetStream("application/octet-stream") {
		t.Log("ok")
	} else {
		t.Error("fail")
	}
}

func TestImage(t *testing.T) {
	if Image("image/png") {
		t.Log("ok")
	} else {
		t.Error("fail")
	}
}

func TestText(t *testing.T) {
	if Text("application/json") {
		t.Log("ok")
	} else if Text("application/xml") {
		t.Log("ok")
	} else if Text("text/html; charset=utf-8") {
		t.Log("ok")
	} else {
		t.Error("fail")
	}
}
