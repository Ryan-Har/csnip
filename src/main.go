package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Ryan-Har/csnip/database"
	"github.com/Ryan-Har/csnip/options"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/google/uuid"
)

func main() {
	db, err := database.NewSQLiteHandler()
	if err != nil {
		log.Fatal(err)
	}

	opt, err := options.GetOptions()
	if err != nil {
		log.Fatal("error getting options", err)
	}

	switch opt.RunType {
	case options.RunTypeCli:
		opt.CliOpts.Run(db)
	case options.RunTypeTui:
		fmt.Println("running as tui")
	case options.RunTypeDaemon:
		fmt.Println("running as daemon")
	}

	//exampleUseOfChroma()
}

func exampleUseOfChroma() {
	db, err := database.NewSQLiteHandler()
	if err != nil {
		log.Fatal(err)
	}

	snip, err := db.GetSnippetByUUID(uuid.MustParse("cde66214-d303-4f5c-8bf5-ca9ae60dc36f"))
	if err != nil {
		log.Fatal(err)
	}
	lexer := lexers.Get(snip.Language.String())
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("dracula")

	formatter := formatters.Get("ANSI")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	iterator, err := lexer.Tokenise(nil, snip.Code)
	if err != nil {
		log.Fatal(err)
	}

	err = formatter.Format(os.Stdout, style, iterator)
	if err != nil {
		log.Fatal(err)
	}
}
