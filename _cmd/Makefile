all: gen_unhex gen_tokens

gen_unhex:
	go run gen_unhex.go >unhex.go
	mv unhex.go ../

gen_tokens:
	go run gen_token.go >tokens.go
	mv tokens.go ../
