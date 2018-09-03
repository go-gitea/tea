// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import "fmt"

var (
	showLog bool
)

// Println println content according the flag
func Println(a ...interface{}) {
	if showLog {
		fmt.Println(a...)
	}
}

// Printf printf content according the flag
func Printf(format string, a ...interface{}) {
	if showLog {
		fmt.Printf(format, a...)
	}
}

// Error println content as an error information
func Error(a ...interface{}) {
	fmt.Println(a...)
}

// Errorf printf content as an error information
func Errorf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}
