package httparser

type state uint8

const (
	dead state = iota + 1
	reqOrRspStart
)
