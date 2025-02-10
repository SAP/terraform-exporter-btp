package generictools

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func GetHclFile(entryName string) (f *hclwrite.File) {
	src, err := os.ReadFile(entryName)
	if err != nil {
		log.Printf("Failed to read file %q: %s", entryName, err)
		return
	}

	f, diags := hclwrite.ParseConfig(src, entryName, hcl.Pos{Line: 1, Column: 1})

	if diags.HasErrors() {
		for _, diag := range diags {
			if diag.Subject != nil {
				log.Printf("[%s:%d] %s: %s", diag.Subject.Filename, diag.Subject.Start.Line, diag.Summary, diag.Detail)
			} else {
				log.Printf("%s: %s", diag.Summary, diag.Detail)
			}
		}
		return
	}
	return f
}

func ProcessChanges(f *hclwrite.File, path string) {
	changed := checkForChanges(f, path)

	if changed {
		info, _ := os.Lstat(path)
		err := os.WriteFile(path, f.Bytes(), info.Mode())
		if err != nil {
			log.Printf("Failed to write file %q: %s", path, err)
			return
		}
	}
}

func CreateVariablesFile(contentToCreate VariableContent, directory string) {
	f := hclwrite.NewEmptyFile()

	rootBody := f.Body()

	for key, value := range contentToCreate {
		varBlock := rootBody.AppendNewBlock("variable", []string{key})
		varBody := varBlock.Body()

		varBody.SetAttributeRaw("type", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenStringLit,
				Bytes: []byte("string"),
			},
		})

		varBody.SetAttributeRaw("description", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenStringLit,
				Bytes: []byte("\"" + value.Description + "\""),
			},
		})

		varBody.SetAttributeRaw("default", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenStringLit,
				Bytes: []byte("\"" + value.Value + "\""),
			},
		})
	}

	variablesSrc := hclwrite.Format(f.Bytes())

	path := filepath.Join(directory, "variables.tf")
	err := os.WriteFile(path, variablesSrc, 0644)
	if err != nil {
		log.Printf("Failed to write file %q: %s", path, err)
		return
	}
}

func ReplaceStringToken(tokens hclwrite.Tokens, identifier string) (replacedTokens hclwrite.Tokens, valueForVariable string) {

	oQuote := tokens[0]
	strTok := tokens[1]
	cQuote := tokens[2]
	if oQuote.Type == hclsyntax.TokenOQuote && strTok.Type == hclsyntax.TokenQuotedLit && cQuote.Type == hclsyntax.TokenCQuote {
		valueForVariable = string(strTok.Bytes)
		return hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte("var." + identifier),
			},
		}, valueForVariable
	}

	return tokens, ""
}

func ReplaceDependency(tokens hclwrite.Tokens, dependencyAddress string) (replacedTokens hclwrite.Tokens) {

	oQuote := tokens[0]
	strTok := tokens[1]
	cQuote := tokens[2]
	if oQuote.Type == hclsyntax.TokenOQuote && strTok.Type == hclsyntax.TokenQuotedLit && cQuote.Type == hclsyntax.TokenCQuote {
		return hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte(dependencyAddress),
			},
		}
	}

	return tokens
}

func checkForChanges(f *hclwrite.File, path string) (changed bool) {

	changed = false

	originalContent, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read file %q: %s", path, err)
		return
	}

	updatedContent := f.Bytes()
	if !bytes.Equal(originalContent, updatedContent) {
		changed = true
	}
	return changed
}
