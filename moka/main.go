package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
)

func main() {
	Converter()
}

func Converter() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: moka <URL> <NAME.md>")
		os.Exit(1)
	}
	input := os.Args[1]
	cmdCurl := exec.Command("curl", "--no-progress-meter", input)
	output, err := cmdCurl.CombinedOutput()
	if err != nil {
		log.Fatalf("Error during the CURL process:\n%s", err)
	}

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(
				commonmark.WithStrongDelimiter("**"),
				commonmark.WithBulletListMarker("-"),
				commonmark.WithCodeBlockFence("```"),
				commonmark.WithHeadingStyle(commonmark.HeadingStyleATX),
				commonmark.WithLinkEmptyContentBehavior(commonmark.LinkBehaviorSkip),
				commonmark.WithLinkEmptyHrefBehavior(commonmark.LinkBehaviorSkip),
				// ...additional configurations for the plugin
			),
			table.NewTablePlugin(
				table.WithHeaderPromotion(true),
				table.WithSkipEmptyRows(true),
				table.WithSpanCellBehavior(table.SpanBehaviorEmpty),
			// ...additional plugins (e.g. table)
			),
		),
	)

	markdown, err := conv.ConvertString(string(output))
	if err != nil {
		log.Fatal(err)
	}
	mdFile := os.Args[2]
	err = os.WriteFile(mdFile, []byte(markdown), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done")
}
