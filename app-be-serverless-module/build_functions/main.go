package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// functionList := getLambdaFunctionList()
	if len(os.Args) < 2 {
		fmt.Println("please supply function name")
		os.Exit(1)
	}
	functionName := os.Args[1]

	cmd := exec.Command("go", "build", "-o", fmt.Sprintf("dist/functions/%s", functionName), fmt.Sprintf("functions/%s/main.go", functionName))
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "GOARCH=amd64")
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Finished compiling %s function", functionName)
}
