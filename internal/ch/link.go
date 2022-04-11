//go:build chimera
// +build chimera

package ch

/*
#cgo pkg-config: --static libch
#cgo linux LDFLAGS: -lm -lstdc++ -lpcre
#cgo darwin LDFLAGS: -lpcre
*/
import "C"
