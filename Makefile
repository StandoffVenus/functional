GO=go
TEST_FLAGS=-timeout 5s -shuffle on -cover

test:
	${GO} test ${TEST_FLAGS} ./...

view-coverage: TEST_FLAGS=-coverprofile coverage.out
view-coverage: test
	go tool cover -html=coverage.out