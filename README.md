# fuzzy

## What will it do?
`fuzzy` is a codegen tool for fuzz tests in Go. It uses the [gofuzz](www.github.com/google/gofuzz)
package to fuzz receivers and function arguments, and then executes the functions.
This takes place in a `*_test.go` file, so that the new tests hook into the existing runs.
