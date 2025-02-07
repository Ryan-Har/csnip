package options

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Ryan-Har/csnip/cli"
	"github.com/Ryan-Har/csnip/common/models"
)

type RunType int

const (
	RunTypeDaemon RunType = iota
	RunTypeCli
	RunTypeTui
)

type Options struct {
	RunType RunType
	CliOpts cli.CLIOpts
}

func GetOptions() (Options, error) {
	var opt Options

	// global flags
	helpFlag := flag.Bool("h", false, "Show help message")
	dFlag := flag.Bool("d", false, "Run as Daemon")

	flag.Parse()

	if *helpFlag {
		fmt.Println("csnip <subcommand> flags")
		fmt.Println("  subcommands: get, add, update")
		fmt.Println("  csnip <subcommand> -h for help")
		fmt.Println()
		flag.Usage()

		os.Exit(0)
	}

	if *dFlag {
		opt.RunType = RunTypeDaemon // not yet implemented
		return opt, nil
	} else if len(os.Args) == 1 {
		opt.RunType = RunTypeTui // not yet implemented
		return opt, nil
	} else {
		opt.RunType = RunTypeCli
	}

	switch strings.ToLower(os.Args[1]) {
	case "get":
		opt.CliOpts = handleGetFlagset()
	case "add":
		opt.CliOpts = handleAddFlagset()
	case "update":
		opt.CliOpts = handleUpdateFlagset()
	case "delete":
		opt.CliOpts = handleDeleteFlagset()
	default:
		fmt.Println("Unknown command: ", os.Args[1])
		os.Exit(1)
	}

	opt.CliOpts.Theme = CodeSyntaxHighlightingTheme
	return opt, nil
}

func handleGetFlagset() cli.CLIOpts {
	var cliOpts cli.CLIOpts
	cliOpts.OptType = cli.OptTypeGet
	cliOpts.FlagOptions = map[cli.FlagOption]string{}

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	allFlag := getCmd.Bool("a", false, "Get a list of code snippets without filtering")
	langFlag := getCmd.String("l", "", "Get a list of code snippets matching the language")
	tagFlag := getCmd.String("t", "", "Get a list of code snippets matching the tag")
	idFlag := getCmd.String("i", "", "Get by uuid of the code snippet")

	getCmd.Parse(os.Args[2:])
	if getCmd.Parsed() {
		if getCmd.NFlag() == 0 || *allFlag {
			cliOpts.FlagOptions[cli.FlagOptionAll] = "all"
		}
		if *idFlag != "" {
			cliOpts.FlagOptions[cli.FlagOptionUUID] = *idFlag
			return cliOpts
		}
		if *langFlag != "" && *tagFlag != "" {
			cliOpts.FlagOptions[cli.FlagOptionLanguage] = *langFlag
			cliOpts.FlagOptions[cli.FlagOptionTag] = *tagFlag
			return cliOpts
		}
		if *langFlag != "" {
			cliOpts.FlagOptions[cli.FlagOptionLanguage] = *langFlag
			return cliOpts
		}
		if *tagFlag != "" {
			cliOpts.FlagOptions[cli.FlagOptionTag] = *tagFlag
			return cliOpts
		}
	}
	return cliOpts
}

func handleAddFlagset() cli.CLIOpts {
	var cliOpts cli.CLIOpts
	cliOpts.OptType = cli.OptTypeAdd
	cliOpts.FlagOptions = map[cli.FlagOption]string{}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	nameFlag := addCmd.String("n", "", "Optional friendly Name used to reference the snippet of code")
	codeFlag := addCmd.String("c", "", "The snippet of code being stored. \"-\" to read from stdin, provide a file or a string of code")
	langFlag := addCmd.String("l", "", "Language of the snippet of code")
	tagsFlag := addCmd.String("t", "", "Optional comma seperated list of tags to assign to the snippet of code")
	descriptionFlag := addCmd.String("d", "", "Optional description for the snippet of code")

	addCmd.Parse(os.Args[2:])
	if addCmd.Parsed() {

		if addCmd.NFlag() == 0 || *codeFlag == "" || *langFlag == "" {
			fmt.Println("Both code (-c) and language (-l) flags must be used")
			addCmd.Usage()
			os.Exit(1)
		}

		// read from stdin with - or read from file if it exists, otherwise treat it as code input.
		if *codeFlag == "-" {
			input, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println("Error reading from stdin")
				os.Exit(1)
			}
			cliOpts.FlagOptions[cli.FlagOptionCode] = string(input)
		} else if _, err := os.Stat(*codeFlag); err == nil {
			content, err := os.ReadFile(*codeFlag)
			if err != nil {
				fmt.Println("Error reading file: ", *codeFlag)
				os.Exit(1)
			}
			cliOpts.FlagOptions[cli.FlagOptionCode] = string(content)
		} else {
			cliOpts.FlagOptions[cli.FlagOptionCode] = *codeFlag
		}

		validLanguage := models.ValidateLanguage(*langFlag)
		if validLanguage == models.LanguageOther {
			// prompt for the a supported language
			fmt.Println("Language provided unknown, please enter one of the following:")
			fmt.Println(strings.Join(models.ListValidLanguages()[:], ","))
			os.Exit(1)
		} else {
			cliOpts.FlagOptions[cli.FlagOptionLanguage] = validLanguage.String()
		}

		if *nameFlag != "" {
			cliOpts.FlagOptions[cli.FlagOptionName] = *nameFlag
		}
		if *tagsFlag != "" {
			cliOpts.FlagOptions[cli.FlagOptionTag] = *tagsFlag
		}
		if *descriptionFlag != "" {
			cliOpts.FlagOptions[cli.FlagOptionDescription] = *descriptionFlag
		}

	}

	return cliOpts
}

func handleUpdateFlagset() cli.CLIOpts {
	var cliOpts cli.CLIOpts
	cliOpts.OptType = cli.OptTypeUpdate
	cliOpts.FlagOptions = map[cli.FlagOption]string{}

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	idFlag := updateCmd.String("i", "", "Get by uuid of the code snippet")
	codeFlag := updateCmd.String("c", "", "The snippet of code being stored. \"-\" to read from stdin, provide a file or a string of code")

	updateCmd.Parse(os.Args[2:])
	if updateCmd.Parsed() {
		if updateCmd.NFlag() == 0 || *codeFlag == "" || *idFlag == "" {
			fmt.Println("Both code (-c) and uuid (-i) flags must be used")
			updateCmd.Usage()
			os.Exit(1)
		}

		cliOpts.FlagOptions[cli.FlagOptionUUID] = *idFlag

		// read from stdin with - or read from file if it exists, otherwise treat it as code input.
		if *codeFlag == "-" {
			input, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println("Error reading from stdin")
				os.Exit(1)
			}
			cliOpts.FlagOptions[cli.FlagOptionCode] = string(input)
		} else if _, err := os.Stat(*codeFlag); err == nil {
			content, err := os.ReadFile(*codeFlag)
			if err != nil {
				fmt.Println("Error reading file: ", *codeFlag)
				os.Exit(1)
			}
			cliOpts.FlagOptions[cli.FlagOptionCode] = string(content)
		} else {
			cliOpts.FlagOptions[cli.FlagOptionCode] = *codeFlag
		}
	}

	return cliOpts
}

func handleDeleteFlagset() cli.CLIOpts {
	var cliOpts cli.CLIOpts
	cliOpts.OptType = cli.OptTypeDelete
	cliOpts.FlagOptions = map[cli.FlagOption]string{}

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	idFlag := deleteCmd.String("i", "", "Delete by uuid of the code snippet")

	deleteCmd.Parse(os.Args[2:])
	if deleteCmd.Parsed() {
		if deleteCmd.NFlag() == 0 || *idFlag == "" {
			fmt.Println(" uuid (-i) flags must be used")
			deleteCmd.Usage()
			os.Exit(1)
		}

		cliOpts.FlagOptions[cli.FlagOptionUUID] = *idFlag
	}

	return cliOpts
}
