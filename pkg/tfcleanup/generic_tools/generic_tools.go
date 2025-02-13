package generictools

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
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
		rootBody.AppendNewline()
	}

	variablesSrc := hclwrite.Format(f.Bytes())

	path := filepath.Join(directory, "variables.tf")
	err := os.WriteFile(path, variablesSrc, 0644)
	if err != nil {
		log.Printf("Failed to write file %q: %s", path, err)
		return
	}
}

func ReplaceStringTokenVar(tokens hclwrite.Tokens, identifier string) (replacedTokens hclwrite.Tokens, valueForVariable string) {
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
				Bytes: []byte(dependencyAddress + ".id"),
			},
		}
	}

	return tokens
}

func GetStringToken(tokens hclwrite.Tokens) (value string) {
	oQuote := tokens[0]
	strTok := tokens[1]
	cQuote := tokens[2]
	if oQuote.Type == hclsyntax.TokenOQuote && strTok.Type == hclsyntax.TokenQuotedLit && cQuote.Type == hclsyntax.TokenCQuote {
		value = string(strTok.Bytes)
	}

	return value
}

func ExtractBlockInformation(inBlocks []string) (blockType string, blockIdentifier string, resourceAddress string) {
	blockType = strings.Split(inBlocks[0], ",")[0]
	blockIdentifier = strings.Split(inBlocks[0], ",")[1]
	blockAddress := strings.Split(inBlocks[0], ",")[2]
	resourceAddress = blockIdentifier + "." + blockAddress

	return blockType, blockIdentifier, resourceAddress
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

func IsGlobalAccountParent(btpClient *btpcli.ClientFacade, parentId string) (isParent bool) {
	globalAccountId, _ := btpcli.GetGlobalAccountId(btpClient)

	if parentId == globalAccountId {
		isParent = true
	}
	return
}

func RemoveConfigBlock(body *hclwrite.Body, blockIdentifier string, resourceAddress string) {
	for _, block := range body.Blocks() {
		if block.Labels()[0] == blockIdentifier && block.Labels()[1] == resourceAddress {
			body.RemoveBlock(block)
			break
		}
	}
}

func RemoveEmptyAttributes(body *hclwrite.Body) {
	for name, attr := range body.Attributes() {
		tokens := attr.Expr().BuildTokens(nil)

		// Check for a NULL value
		if len(tokens) == 1 && string(tokens[0].Bytes) == EmptyString {
			body.RemoveAttribute(name)
		}

		// Check for an empty JSON encoded string or an empty Map
		var combinedString string
		if len(tokens) == 5 || len(tokens) == 2 {
			for _, token := range tokens {
				combinedString += string(token.Bytes)
			}
		}

		if combinedString == EmptyJson || combinedString == EmptyMap {
			body.RemoveAttribute(name)
		}
	}
}

func ReplaceMainDependency(body *hclwrite.Body, mainIdentifier string, mainAddress string) {
	if mainAddress == "" {
		return
	}

	for name, attr := range body.Attributes() {
		tokens := attr.Expr().BuildTokens(nil)

		if name == mainIdentifier && len(tokens) == 3 {
			replacedTokens := ReplaceDependency(tokens, mainAddress)
			body.SetAttributeRaw(name, replacedTokens)
		}
	}
}

func ProcessParentAttribute(body *hclwrite.Body, description string, btpClient *btpcli.ClientFacade, variables *VariableContent) {
	parentAttr := body.GetAttribute(ParentIdentifier)
	if parentAttr == nil {
		return
	}

	tokens := parentAttr.Expr().BuildTokens(nil)
	if len(tokens) == 3 {

		parentId := GetStringToken(tokens)

		if IsGlobalAccountParent(btpClient, parentId) {
			body.RemoveAttribute(ParentIdentifier)
		} else {
			replacedTokens, parentValue := ReplaceStringTokenVar(tokens, ParentIdentifier)
			if parentValue != "" {
				(*variables)[ParentIdentifier] = VariableInfo{
					Description: description,
					Value:       parentValue,
				}
			}
			body.SetAttributeRaw(ParentIdentifier, replacedTokens)
		}
	}
}

func ReplaceAttribute(body *hclwrite.Body, identifier string, description string, variables *VariableContent) {
	attribute := body.GetAttribute(identifier)

	if attribute != nil {
		tokens := attribute.Expr().BuildTokens(nil)

		if len(tokens) == 3 {
			replacedTokens, attrValue := ReplaceStringTokenVar(tokens, identifier)
			(*variables)[identifier] = VariableInfo{
				Description: description,
				Value:       attrValue,
			}
			body.SetAttributeRaw(identifier, replacedTokens)
		}
	}
}
