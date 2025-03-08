package main

import (
	"fmt"
	"os"

	dtree "github.com/muthu-kumar-u/go-dtree/tree"
)

func main() {
	bytes, err := os.ReadFile(os.Getenv("TREE_FILE"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tree, err := dtree.CreateTreeFromjson(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Sample decision input
	input := map[string]interface{}{
		"ShouldBuy":       true,
		"CheckBudget":     false,
		"NecessityCheck":  true,
		"AlternativeCheck": false,
	}

	result, err  := tree.Decide(input)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Decision Result:", result)
}
