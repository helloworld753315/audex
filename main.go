package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "audex",
		Usage: "Convert MP4 to M4A with metadata preservation",
		Action: func(c *cli.Context) error {
			if c.NArg() < 2 {
				return fmt.Errorf("requires input and output files\nUsage: audex <input.mp4> <output.m4a>")
			}

			input := c.Args().Get(0)
			output := c.Args().Get(1)

			if err := convert(input, output); err != nil {
				return err
			}

			fmt.Printf("Successfully converted %s to %s\n", input, output)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func convert(input, output string) error {
	return fmt.Errorf("not implemented yet")
}
