package scode

// AllowFunc allow code to be allocated
type AllowFunc func(code string) bool

func AllowAll(code string) bool { return true }
