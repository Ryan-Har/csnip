package models

import (
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Language Enum
type Language string

const (
	LanguageOther      Language = "Other"
	LanguageGo         Language = "Go"
	LanguagePython     Language = "Python"
	LanguageJava       Language = "Java"
	LanguageJavascript Language = "Javascript"
	LanguageCsharp     Language = "C#"
	LanguageC          Language = "C"
	LanguageRuby       Language = "Ruby"
	LanguageRust       Language = "Rust"
	LanguagePhp        Language = "Php"
	LanguageSwift      Language = "Swift"
	LanguageSql        Language = "Sql"
	LanguageKotlin     Language = "Kotlin"
	LanguageBash       Language = "Bash"
	LanguagePowershell Language = "Powershell"
)

// Language Validation Map
var validLanguages = map[string]Language{
	"Go":         LanguageGo,
	"Python":     LanguagePython,
	"Java":       LanguageJava,
	"Javascript": LanguageJavascript,
	"C#":         LanguageCsharp,
	"C":          LanguageC,
	"Ruby":       LanguageRuby,
	"Rust":       LanguageRust,
	"Php":        LanguagePhp,
	"Swift":      LanguageSwift,
	"Sql":        LanguageSql,
	"Kotlin":     LanguageKotlin,
	"Bash":       LanguageBash,
	"Powershell": LanguagePowershell,
}

// Stringer implementation for Language enum
func (l Language) String() string {
	return string(l)
}

func ValidateLanguage(lang string) Language {
	// normalise string so that it always titles
	normalisedLangString := cases.Title(language.BritishEnglish).String(strings.ToLower(lang))

	if validLang, exists := validLanguages[normalisedLangString]; exists {
		return validLang
	}
	return LanguageOther // Default to "Other" if not recognized
}

func ListValidLanguages() []string {
	return slices.Collect(maps.Keys(validLanguages))
}

type CodeSnippet struct {
	ID           int64
	Uuid         uuid.UUID
	Name         string
	Code         string
	Language     Language
	Tags         string
	Description  string
	Source       string
	DateAdded    time.Time
	Version      int64
	SupersededBy int64
}

// map containing the hello code for various languages
var helloWorldMap = map[Language]string{
	LanguageGo:         "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}",
	LanguagePython:     "print(\"Hello, World!\")",
	LanguageJavascript: "console.log(\"Hello, World!\");",
	LanguageJava:       "public class Main {\n    public static void main(String[] args) {\n        System.out.println(\"Hello, World!\");\n    }\n}",
	LanguageC:          "#include <stdio.h>\n\nint main() {\n    printf(\"Hello, World!\n\");\n    return 0;\n}",
	LanguageRuby:       "puts \"Hello, World!\"",
	LanguagePhp:        "<?php\n echo \"Hello, World!\";\n?>",
	LanguageSwift:      "import Foundation\n\nprint(\"Hello, World!\")",
	LanguageKotlin:     "fun main() {\n    println(\"Hello, World!\")\n}",
}

func GetHelloWorldExamples() map[Language]string {
	return helloWorldMap
}
