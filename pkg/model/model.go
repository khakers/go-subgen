package model

//go:generate go-enum -type=Model -all=false -string=true -new=true -string=true -text=true -json=true -yaml=false

type Model uint8

const (
	Tiny_en Model = iota
	Tiny
	Base_en
	Base
	Small_en
	Small
	Medium_en
	Medium
	Large_v1
	Large
)
