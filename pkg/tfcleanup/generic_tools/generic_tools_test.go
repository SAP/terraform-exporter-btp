package generictools

/*
The following functions won't be tested as they just orchestrate external functions
or are implictly tested:
- GetHclFile
- ProcessChanges
- CreateVariablesFile
- checkForChanges
- IsGlobalAccountParent
- RemoveUnusedImports
- RemoveEmptyFiles
*/

import (
	"testing"

	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestProcessParent(t *testing.T) {

	srcFileWithGaParent, trgtFileWithGaParent := testutils.GetHclFilesById("sa_with_ga_parent")
	srcFileWoGaParent, trgtFileWoGaParent := testutils.GetHclFilesById("sa_without_ga_parent")

	emptyTestContent := make(VariableContent)
	targetVariables := make(VariableContent)
	targetVariables[ParentIdentifier] = VariableInfo{
		Description: "Some Text",
		Value:       "directory",
	}

	tests := []struct {
		name          string
		src           *hclwrite.File
		trgt          *hclwrite.File
		description   string
		trgtVariables *VariableContent
	}{
		{
			name:          "Test removal of global account as parent",
			src:           srcFileWithGaParent,
			trgt:          trgtFileWithGaParent,
			description:   "Some Text",
			trgtVariables: &emptyTestContent,
		},
		{
			name:          "Test keep parent",
			src:           srcFileWoGaParent,
			trgt:          trgtFileWoGaParent,
			description:   "Some Text",
			trgtVariables: &targetVariables,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			contentToCreate := make(VariableContent)
			blocks := tt.src.Body().Blocks()
			// we assume one rrsource entry in the blocks file
			ProcessParentAttribute(blocks[0].Body(), tt.description, nil, &contentToCreate)
			assert.NoError(t, testutils.AreHclFilesEqual(tt.trgt, tt.src))
			assert.Equal(t, tt.trgtVariables, &contentToCreate)
		})
	}
}

/*
	TODO
RemoveConfigBlock
RemoveImportBlock
RemoveEmptyAttributes
ReplaceDependency
ReplaceAttribute
*/

func TestReplaceStringTokenVar(t *testing.T) {
	tests := []struct {
		name             string
		tokens           hclwrite.Tokens
		identifier       string
		expectedTokens   hclwrite.Tokens
		expectedVariable string
	}{
		{
			name: "Valid quoted string",
			tokens: hclwrite.Tokens{
				{Type: hclsyntax.TokenOQuote, Bytes: []byte("\"")},
				{Type: hclsyntax.TokenQuotedLit, Bytes: []byte("some_value")},
				{Type: hclsyntax.TokenCQuote, Bytes: []byte("\"")},
			},
			identifier: "some_identifier",
			expectedTokens: hclwrite.Tokens{
				{Type: hclsyntax.TokenIdent, Bytes: []byte("var.some_identifier")},
			},
			expectedVariable: "some_value",
		},
		{
			name: "Invalid tokens",
			tokens: hclwrite.Tokens{
				{Type: hclsyntax.TokenIdent, Bytes: []byte("some_value")},
			},
			identifier:       "some_identifier",
			expectedTokens:   hclwrite.Tokens{{Type: hclsyntax.TokenIdent, Bytes: []byte("some_value")}},
			expectedVariable: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replacedTokens, valueForVariable := ReplaceStringTokenVar(tt.tokens, tt.identifier)
			assert.Equal(t, tt.expectedTokens, replacedTokens)
			assert.Equal(t, tt.expectedVariable, valueForVariable)
		})
	}
}

func TestReplaceDependency(t *testing.T) {
	tests := []struct {
		name              string
		tokens            hclwrite.Tokens
		dependencyAddress string
		expectedTokens    hclwrite.Tokens
	}{
		{
			name: "Valid quoted string",
			tokens: hclwrite.Tokens{
				{Type: hclsyntax.TokenOQuote, Bytes: []byte("\"")},
				{Type: hclsyntax.TokenQuotedLit, Bytes: []byte("some_value")},
				{Type: hclsyntax.TokenCQuote, Bytes: []byte("\"")},
			},
			dependencyAddress: "module.some_module",
			expectedTokens: hclwrite.Tokens{
				{Type: hclsyntax.TokenIdent, Bytes: []byte("module.some_module.id")},
			},
		},
		{
			name: "Invalid tokens",
			tokens: hclwrite.Tokens{
				{Type: hclsyntax.TokenIdent, Bytes: []byte("some_value")},
			},
			dependencyAddress: "module.some_module",
			expectedTokens: hclwrite.Tokens{
				{Type: hclsyntax.TokenIdent, Bytes: []byte("some_value")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replacedTokens := ReplaceDependency(tt.tokens, tt.dependencyAddress)
			assert.Equal(t, tt.expectedTokens, replacedTokens)
		})
	}
}

func TestGetStringToken(t *testing.T) {
	tests := []struct {
		name           string
		tokens         hclwrite.Tokens
		expectedResult string
	}{
		{
			name: "Valid quoted string",
			tokens: hclwrite.Tokens{
				{Type: hclsyntax.TokenOQuote, Bytes: []byte("\"")},
				{Type: hclsyntax.TokenQuotedLit, Bytes: []byte("some_value")},
				{Type: hclsyntax.TokenCQuote, Bytes: []byte("\"")},
			},
			expectedResult: "some_value",
		},
		{
			name: "Invalid tokens length",
			tokens: hclwrite.Tokens{
				{Type: hclsyntax.TokenOQuote, Bytes: []byte("\"")},
				{Type: hclsyntax.TokenQuotedLit, Bytes: []byte("some_value")},
			},
			expectedResult: "",
		},
		{
			name: "Invalid token types",
			tokens: hclwrite.Tokens{
				{Type: hclsyntax.TokenIdent, Bytes: []byte("some_value")},
				{Type: hclsyntax.TokenQuotedLit, Bytes: []byte("some_value")},
				{Type: hclsyntax.TokenIdent, Bytes: []byte("some_value")},
			},
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStringToken(tt.tokens)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
func TestExtractBlockInformation(t *testing.T) {
	tests := []struct {
		name              string
		inBlocks          []string
		expectedBlockType string
		expectedBlockId   string
		expectedResAddr   string
	}{
		{
			name:              "Valid block information",
			inBlocks:          []string{"resource_type,resource_id,resource_address"},
			expectedBlockType: "resource_type",
			expectedBlockId:   "resource_id",
			expectedResAddr:   "resource_id.resource_address",
		},
		{
			name:              "Empty block information",
			inBlocks:          []string{"resource_type,,resource_address"},
			expectedBlockType: "resource_type",
			expectedBlockId:   "",
			expectedResAddr:   "",
		},
		{
			name:              "Invalid block information",
			inBlocks:          []string{"resource_type,resource_id"},
			expectedBlockType: "",
			expectedBlockId:   "",
			expectedResAddr:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockType, blockIdentifier, resourceAddress := ExtractBlockInformation(tt.inBlocks)
			assert.Equal(t, tt.expectedBlockType, blockType)
			assert.Equal(t, tt.expectedBlockId, blockIdentifier)
			assert.Equal(t, tt.expectedResAddr, resourceAddress)
		})
	}
}
