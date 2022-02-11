# executor

[![Tests](https://github.com/solar-jsoc/executor/workflows/Tests/badge.svg)](https://github.com/solar-jsoc/executor/actions)
[![Coverage Status](https://coveralls.io/repos/github/solar-jsoc/executor/badge.svg?branch=master)](https://coveralls.io/github/solar-jsoc/executor?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/solar-jsoc/executor)](https://goreportcard.com/report/github.com/solar-jsoc/executor)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/solar-jsoc/executor)](https://pkg.go.dev/github.com/solar-jsoc/executor)


## Examples
```golang
package main

import (
	"fmt"

 	"github.com/solar-jsoc/executor"
)

func main() {
	stdOut, stdErr, err := executor.Exec("nmap -sV 1.1.1.1 -Pn -oX out.xml > /dev/null && cat out.xml", "", nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println(stdErr)
		return
	}
	
	fmt.Println(string(stdOut))
}
```
