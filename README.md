
Go load testing tool for testing line server

1. How to run: 

go build
./goloadtest ip:port.

2. How it works

The test result will be written into perf_test.txt file, and console.

The tool simulate multiple clients. Each client creates connection, and then generate random line numbers, and issue 10 GET requests to line server continously.

Once a set of clients finish, it will increment # of clients, and do the same tests again until # of clients reachs a limit.




