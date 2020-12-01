
start:
	nodemon --exec go run main.go --signal SIGTERM

ginkgo:
	cd test && ginkgo watch
