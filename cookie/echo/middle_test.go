package jarmiddle

import "testing"

func TestGet(t *testing.T) {
	// 리퀘스트에 쿠키가 있으면
	// Get 으로 가져올 수 있다
	// response 에는 없다
}

func TestSet(t *testing.T) {
	// Set으로 쿠키값을 지정하면
	// response 에 해당값이 있다
}

func TestRemove(t *testing.T) {
	// Request에 쿠키가 있을 때
	// Remove로 쿠키를 삭제하면
	// response 에 해당값이 빈 상태로 온다

	// Request에 쿠키가 없을 때
	// Remove로 쿠키를 삭제하면
	// response 에 해당값이 없다
}
