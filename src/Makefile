
all: ex00 ex01 ex02

ex00:
	@cd ex00 && go build
	@mv ex00/server_0 test

ex01:
	@cd ex01/client && go build
	@cd ex01/server && go build
	@mv ex01/client/candy-client test
	@mv ex01/server/server_1 test

ex02:
	@cd ex02 && go build
	@mv ex02/server_2 test

.PHONY: ex00 ex01 ex02