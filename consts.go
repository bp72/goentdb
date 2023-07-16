package goentdb

type Origin uint

const (
	OriginUnkown     Origin = 0
	OriginXvideos    Origin = 1
	OriginPorntube   Origin = 2
	OriginEporner    Origin = 3
	OriginPornone    Origin = 4
	OriginCumlouder  Origin = 5
	OriginSuperporn  Origin = 6
	OriginAlphaporno Origin = 8
)

type EntKeywordType uint

const (
	EntKeywordTag EntKeywordType = iota
	EntKeywordModel
	EntKeywordKeyword
)
