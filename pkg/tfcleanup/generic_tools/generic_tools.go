package generictools

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
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
	if len(tokens) != 3 {
		return tokens, ""
	}

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
	if len(tokens) != 3 {
		return tokens
	}

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
	if len(tokens) != 3 {
		return ""
	}

	oQuote := tokens[0]
	strTok := tokens[1]
	cQuote := tokens[2]
	if oQuote.Type == hclsyntax.TokenOQuote && strTok.Type == hclsyntax.TokenQuotedLit && cQuote.Type == hclsyntax.TokenCQuote {
		value = string(strTok.Bytes)
	}

	return value
}

func ExtractBlockInformation(inBlocks []string) (blockType string, blockIdentifier string, resourceAddress string) {

	if len(strings.Split(inBlocks[0], ",")) < 3 {
		return
	}

	blockType = strings.Split(inBlocks[0], ",")[0]
	blockIdentifier = strings.Split(inBlocks[0], ",")[1]
	blockAddress := strings.Split(inBlocks[0], ",")[2]

	if blockIdentifier != "" && blockAddress != "" {
		resourceAddress = blockIdentifier + "." + blockAddress
	}

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
	// Execute Mock in case of testing
	if testing.Testing() {
		return testutils.GetGlobalAccountMockParentData(parentId)
	}

	globalAccountId, _ := btpcli.GetGlobalAccountId(btpClient)

	if parentId == globalAccountId {
		isParent = true
	}
	return
}

func RemoveConfigBlock(body *hclwrite.Body, resourceAddress string) {
	for _, block := range body.Blocks() {
		address := block.Labels()[0] + "." + block.Labels()[1]
		if address == resourceAddress {
			body.RemoveBlock(block)
		}
	}
}

func RemoveImportBlock(body *hclwrite.Body, resourceAddress string, resultStore *map[string]int) {

	taintedBlocks := []*hclwrite.Block{}

	for _, block := range body.Blocks() {

		importTargetAttr := block.Body().GetAttribute("to")

		if importTargetAttr == nil {
			return
		}

		tokens := importTargetAttr.Expr().BuildTokens(nil)
		address := string(tokens[0].Bytes) + string(tokens[1].Bytes) + string(tokens[2].Bytes)

		if address == resourceAddress {
			taintedBlocks = append(taintedBlocks, block)
		}
	}

	for _, block := range taintedBlocks {
		body.RemoveBlock(block)
		(*resultStore)[strings.Split(resourceAddress, ".")[0]] -= 1
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

func ReplaceSpaceDependency(body *hclwrite.Body, spaceIdentifier string, spaceAddress string, spaceId string) {
	if spaceAddress == "" {
		return
	}

	for name, attr := range body.Attributes() {
		tokens := attr.Expr().BuildTokens(nil)

		if name == spaceIdentifier && len(tokens) == 3 && spaceId == GetStringToken(tokens) {
			replacedTokens := ReplaceDependency(tokens, spaceAddress)
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

		if len(tokens) != 3 {
			return
		}
		replacedTokens, attrValue := ReplaceStringTokenVar(tokens, identifier)
		if attrValue == "" {
			return
		}

		(*variables)[identifier] = VariableInfo{
			Description: description,
			Value:       attrValue,
		}
		body.SetAttributeRaw(identifier, replacedTokens)
	}
}

func RemoveUnusedImports(directory string, blocksToRemove *[]BlockSpecifier, resultStore *map[string]int) {
	for _, block := range *blocksToRemove {
		filePath := filepath.Join(directory, block.BlockIdentifier+"_import.tf")
		f := GetHclFile(filePath)
		body := f.Body()
		RemoveImportBlock(body, block.ResourceAddress, resultStore)
		ProcessChanges(f, filePath)
	}
}

func RemoveEmptyFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Failed to read directory %q: %s", dir, err)
		return err
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		info, _ := os.Lstat(path)
		if info.Size() == 0 {
			os.Remove(path)
		}
	}
	return nil
}
