package main

import (
	"fmt"
	"log"
	"os"

	"github.com/helloworld753315/audex/converter"
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
	// メタデータを読み取る
	metadata, err := converter.ReadMetadata(input)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	// デバッグ: メタデータを表示
	fmt.Printf("Title: %s\n", metadata.Title)
	fmt.Printf("Artist: %s\n", metadata.Artist)
	fmt.Printf("Album: %s\n", metadata.Album)
	if metadata.Picture != nil {
		fmt.Printf("Artwork: %s (%d bytes)\n", metadata.Picture.MIMEType, len(metadata.Picture.Data))
	}

	// TODO: 音声ストリーム抽出
	// TODO: メタデータ書き込み

	return fmt.Errorf("conversion not fully implemented yet")
}
