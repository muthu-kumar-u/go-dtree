package main

import (
	"fmt"
	"os"

	dtree "github.com/muthu-kumar-u/go-dtree/tree"
)

func main() {
	// bytes, err := os.ReadFile("test.json")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// tree, err := dtree.CreateTreeFromjson(bytes)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	tree := &dtree.Tree{}
	input := map[string]interface{}{
		"predict" : "CheckBudget",
		"value" : false,
	}
	result, err := tree.Decide(input)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	
	fmt.Println(result)
}