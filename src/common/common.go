package common

import (
	"fmt"
	"github.com/alecthomas/chroma/v2/lexers"
	"os"
)

func ReadFromFile(s string) (string, error) {
	data, err := os.ReadFile(s)
	if err != nil {
		return "", fmt.Errorf("unable to read file: %w", err)
	}

	return string(data), nil
}

func ValidateLanguage(lang string) bool {
	// normalise string so that it always titles
	validLangs := lexers.Names(true)

	for _, item := range validLangs {
		if lang == item {
			return true
		}
	}
	return false
}

func ListValidLanguages() []string {
	return lexers.Names(true)
}

// map containing the hello code for various languages
var helloWorldMap = map[string]string{
	"go":     "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}",
	"python": "print(\"Hello, World!\")",
	"js":     "console.log(\"Hello, World!\");",
	"java":   "public class Main {\n    public static void main(String[] args) {\n        System.out.println(\"Hello, World!\");\n    }\n}",
	"c":      "#include <stdio.h>\n\nint main() {\n    printf(\"Hello, World!\n\");\n    return 0;\n}",
	"ruby":   "puts \"Hello, World!\"",
	"php":    "<?php\n echo \"Hello, World!\";\n?>",
	"swift":  "import Foundation\n\nprint(\"Hello, World!\")",
	"kotlin": "fun main() {\n    println(\"Hello, World!\")\n}",
}

func GetHelloWorldExamples() map[string]string {
	return helloWorldMap
}
