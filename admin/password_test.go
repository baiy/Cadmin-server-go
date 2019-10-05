package admin

import "testing"

func TestPasswrod(t *testing.T) {
	password := new(PasswrodDefault)

	result := "YTc1NzQwYzcwYTM3Yzk1NDk5OGNmZDZlMWYyYjcyOWY3NzU5ODdkZDFmNzIwYWI0YTRmNjJiOTM2NDlmYTE5MXw1NzYzNjg1Mw=="
	if !password.Verify([]byte("123456"), []byte(result)) {
		t.Error("error1")
	}

	if !password.Verify([]byte("123456"), password.Hash([]byte("123456"))) {
		t.Error("error2")
	}
}
