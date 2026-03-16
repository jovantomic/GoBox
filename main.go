package main

import (
	"fmt"
	"os"
	"os/exec"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	switch os.Args[1] {
	case "run":
		cmd := exec.Command(os.Args[2], os.Args[3:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

	default:
		panic("wooooo")
	}
}

func run() {
	fmt.Println("Running the application...")

}
