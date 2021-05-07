all: 

gen: gen_unhex gen_tokens

gen_unhex:
	go run _cmd/gen_unhex.go >unhex.go

gen_tokens:
	go run _cmd/gen_token.go >tokens.go

gen_method:
	go run _cmd/gen_method.go


example.run: example
	- ./request
	- ./response
	- ./request_or_response

example:
	- go build _example/request.go
	- go build _example/response.go
	- go build _example/request_or_response.go


clean:
	- rm request response request_or_response *~
