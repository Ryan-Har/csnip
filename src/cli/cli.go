package cli

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Ryan-Har/csnip/common"
	"github.com/Ryan-Har/csnip/common/models"
	"github.com/Ryan-Har/csnip/database"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/atotto/clipboard"
	"github.com/google/uuid"
)

type CLIOpts struct {
	OptType     OptType
	FlagOptions map[FlagOption]string
	Theme       string
}

type OptType string

const (
	OptTypeGet    OptType = "GET"
	OptTypeUpdate OptType = "UPDATE"
	OptTypeAdd    OptType = "ADD"
	OptTypeDelete OptType = "DELETE"
)

func (o OptType) String() string {
	return string(o)
}

type FlagOption string

const (
	FlagOptionUUID        FlagOption = "UUID"
	FlagOptionLanguage    FlagOption = "Language"
	FlagOptionTag         FlagOption = "Tag"
	FlagOptionAll         FlagOption = "All"
	FlagOptionCode        FlagOption = "Code"
	FlagOptionName        FlagOption = "Name"
	FlagOptionDescription FlagOption = "Description"
)

func (c *CLIOpts) Run(db database.DatabaseInteractions) {
	switch c.OptType {
	case OptTypeGet:
		snippets, err := c.handleGetOptType(db)
		if err != nil {
			if errors.Is(err, database.ErrNoSnippetsFound) {

			}
			fmt.Println("Error occured retrieving code snippets: ", err)
			os.Exit(1)
		}
		if len(snippets) < 1 {
			fmt.Println("No code snippets found with the provided filters")
		}
		// only display when searching with uuid
		if len(snippets) == 1 && c.FlagOptions[FlagOptionUUID] != "" {
			displaySingleSnippet(snippets[0], c.Theme)
		} else {
			displaySnippetList(snippets)
		}
		os.Exit(0)
	case OptTypeAdd:
		err := c.handleAddOptType(db)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Code snippet added to database")
		os.Exit(0)
	case OptTypeUpdate:
		err := c.handleUpdateOptType(db)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Code snippet updated")
		os.Exit(0)
	case OptTypeDelete:
		err := c.handleDeleteOptType(db)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Code snippet deleted")
		os.Exit(0)
	default:
		fmt.Println("Unknown operation")
		os.Exit(1)
	}

}

func (c *CLIOpts) handleAddOptType(db database.DatabaseInteractions) error {
	fOpts := c.FlagOptions
	var cs models.CodeSnippet

	cs.Code = fOpts[FlagOptionCode]

	if common.ValidateLanguage(fOpts[FlagOptionLanguage]) {
		cs.Language = fOpts[FlagOptionLanguage]
	} else {
		return fmt.Errorf("invalid language supplied: %v", fOpts[FlagOptionLanguage])
	}

	if name, ok := fOpts[FlagOptionName]; ok {
		cs.Name = name
	}
	if tags, ok := fOpts[FlagOptionTag]; ok {
		cs.Tags = tags
	}
	if description, ok := fOpts[FlagOptionDescription]; ok {
		cs.Description = description
	}

	err := db.AddNewSnippet(cs)
	if err != nil {
		return fmt.Errorf("unable to handle ADD with the provided options %w", err)
	}

	return nil
}

func (c *CLIOpts) handleGetOptType(db database.DatabaseInteractions) ([]models.CodeSnippet, error) {
	fOpts := c.FlagOptions
	var snippets []models.CodeSnippet

	if len(fOpts) == 2 && fOpts[FlagOptionLanguage] != "" && fOpts[FlagOptionTag] != "" {
		return db.GetSnippetsByLanguageAndTag(fOpts[FlagOptionLanguage], fOpts[FlagOptionTag])
	}
	if fOpts[FlagOptionAll] != "" {
		return db.GetSnippets(1, 100)
	}
	if fOpts[FlagOptionLanguage] != "" {
		return db.GetSnippetsByLanguage(fOpts[FlagOptionLanguage])
	}
	if fOpts[FlagOptionTag] != "" {
		return db.GetSnippetsByTag(fOpts[FlagOptionTag])
	}
	if fOpts[FlagOptionUUID] != "" {
		id, err := uuid.Parse(fOpts[FlagOptionUUID])
		if err != nil {
			return snippets, fmt.Errorf("unable to parse provided UUID")
		}
		idSnip, err := db.GetSnippetByUUID(id)
		if err != nil {
			return snippets, err
		}
		snippets = append(snippets, idSnip)
		return snippets, nil
	}
	return snippets, fmt.Errorf("unable to handle GET with the provided options")
}

func (c *CLIOpts) handleUpdateOptType(db database.DatabaseInteractions) error {
	var cs models.CodeSnippet
	cs.Code = c.FlagOptions[FlagOptionCode]

	id, err := uuid.Parse(c.FlagOptions[FlagOptionUUID])
	if err != nil {
		return fmt.Errorf("unable to parse provided UUID")
	}

	_, err = db.UpdateSnippet(id, cs)
	if err != nil {
		return fmt.Errorf("unable to handle ADD with the provided options %w", err)
	}

	return nil
}

func (c *CLIOpts) handleDeleteOptType(db database.DatabaseInteractions) error {

	id, err := uuid.Parse(c.FlagOptions[FlagOptionUUID])
	if err != nil {
		return fmt.Errorf("unable to parse provided UUID")
	}

	err = db.DeleteSnippetByUUID(id)
	if err != nil {
		return fmt.Errorf("unable to handle DELETE with the provided options %w", err)
	}

	return nil
}

func displaySnippetList(snippets []models.CodeSnippet) {
	fmt.Printf("%-36s	%-25s	%-10s	%-20s	%-30s	%-20s\n", "Uuid", "Name", "Language", "Tags", "Description", "Source")
	for _, s := range snippets {
		fmt.Printf("%-36s	%-25s	%-10s	%-20s	%-30s	%-20s\n",
			truncate(s.Uuid.String(), 36),
			truncate(s.Name, 25),
			truncate(s.Language, 10),
			truncate(s.Tags, 20),
			truncate(s.Description, 30),
			truncate(s.Source, 20),
		)
	}
}

func displaySingleSnippet(snippet models.CodeSnippet, theme string) {
	lexer := lexers.Get(snippet.Language)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get(theme)

	formatter := formatters.Get("terminal")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	iterator, err := lexer.Tokenise(nil, snippet.Code)
	if err != nil {
		log.Fatal(err)
	}

	err = formatter.Format(os.Stdout, style, iterator)
	if err != nil {
		log.Fatal(err)
	}

	// add new line to the end otherwise it doesn't display properly
	fmt.Println()

	_ = clipboard.WriteAll(snippet.Code)
}

func truncate(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength]
	}
	return s
}
