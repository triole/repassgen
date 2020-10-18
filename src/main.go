package main

import (
	"fmt"
)

func main() {
	argparse()
	var entropy float64

	fmt.Printf("\n")
	for i := 1; i <= *argsNo; i++ {

		out, err := Generate(GenerateInput{
			Pattern: *argsPattern,
		})
		if err != nil {
			panic(err)
		}

		if i == 1 {
			entropy = out.PatternEntropy
		}

		fmt.Printf("%s\n", string(out.Password))
		if i == *argsNo {
			fmt.Printf("\n")
		}
	}

	if *argsEntropy == true {
		fmt.Printf("  pattern %s, p. entropy %.2f bits\n\n", *argsPattern, entropy)
	}
}
