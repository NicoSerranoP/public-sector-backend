package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run ./cmd/data <input.sav> <output.csv>")
		return
	}

	if err := convert(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func convert(inputPath, outputPath string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer input.Close()

	data, err := readSav(input)
	if err != nil {
		return fmt.Errorf("reading %s: %w", inputPath, err)
	}

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer output.Close()

	writer := csv.NewWriter(output)

	if err := writer.Write(data.header()); err != nil {
		return err
	}

	if err := writer.WriteAll(data.rows); err != nil {
		return err
	}

	return output.Close()
}
