GO=go
TEST_FLAGS=-timeout 5s -shuffle on -cover

test:
	${GO} test ${TEST_FLAGS} ./...