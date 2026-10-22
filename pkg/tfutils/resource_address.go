package tfutils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

var generatedResourceLabelPattern = regexp.MustCompile(`^export_([0-9]+)$`)

type terraformFile struct {
	path string
	mode fs.FileMode
	file *hclwrite.File
}

// normalizeGeneratedResourceAddresses replaces Terraform's generic export_N
// resource labels with labels derived from the resource type. It updates all
// traversals in the configuration from the same address map, including import
// targets and dependencies created during code cleanup.
func normalizeGeneratedResourceAddresses(directory string) error {
	terraformFiles, err := readTerraformFiles(directory)
	if err != nil {
		return err
	}

	renames := make(map[string]string)
	declaredAddresses := make(map[string]struct{})
	for _, terraformFile := range terraformFiles {
		for _, block := range terraformFile.file.Body().Blocks() {
			labels := block.Labels()
			if block.Type() != "resource" || len(labels) != 2 {
				continue
			}

			resourceType, resourceLabel := labels[0], labels[1]
			declaredAddresses[resourceType+"."+resourceLabel] = struct{}{}
			if normalizedLabel, ok := normalizedResourceLabel(resourceType, resourceLabel); ok {
				renames[resourceType+"."+resourceLabel] = normalizedLabel
			}
		}
	}

	for oldAddress, newLabel := range renames {
		resourceType, _, _ := strings.Cut(oldAddress, ".")
		newAddress := resourceType + "." + newLabel
		if _, exists := declaredAddresses[newAddress]; exists && newAddress != oldAddress {
			return fmt.Errorf("cannot rename %s to %s: target resource address already exists", oldAddress, newAddress)
		}
	}

	for _, terraformFile := range terraformFiles {
		changed := renameResourceBlocks(terraformFile.file.Body(), renames)
		if rewriteResourceReferences(terraformFile.file.Body(), renames) {
			changed = true
		}
		if changed {
			if err := os.WriteFile(terraformFile.path, terraformFile.file.Bytes(), terraformFile.mode.Perm()); err != nil {
				return fmt.Errorf("write Terraform configuration %q: %w", terraformFile.path, err)
			}
		}
	}

	return nil
}

func readTerraformFiles(directory string) ([]terraformFile, error) {
	var terraformFiles []terraformFile
	err := filepath.WalkDir(directory, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || filepath.Ext(path) != ".tf" {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		source, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		parsed, diagnostics := hclwrite.ParseConfig(source, path, hcl.InitialPos)
		if diagnostics.HasErrors() {
			return fmt.Errorf("parse Terraform configuration %q: %s", path, diagnostics.Error())
		}
		terraformFiles = append(terraformFiles, terraformFile{path: path, mode: info.Mode(), file: parsed})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read Terraform configuration: %w", err)
	}
	return terraformFiles, nil
}

func normalizedResourceLabel(resourceType, resourceLabel string) (string, bool) {
	match := generatedResourceLabelPattern.FindStringSubmatch(resourceLabel)
	if match == nil {
		return "", false
	}

	var name string
	switch {
	case strings.HasPrefix(resourceType, "btp_"):
		name = strings.TrimPrefix(resourceType, "btp_")
	case strings.HasPrefix(resourceType, "cloudfoundry_"):
		name = strings.TrimPrefix(resourceType, "cloudfoundry_")
	default:
		return "", false
	}
	if name == "" {
		return "", false
	}
	return name + "_" + match[1], true
}

func renameResourceBlocks(body *hclwrite.Body, renames map[string]string) bool {
	changed := false
	for _, block := range body.Blocks() {
		labels := block.Labels()
		if block.Type() == "resource" && len(labels) == 2 {
			if newLabel, ok := renames[labels[0]+"."+labels[1]]; ok {
				block.SetLabels([]string{labels[0], newLabel})
				changed = true
			}
		}
	}
	return changed
}

func rewriteResourceReferences(body *hclwrite.Body, renames map[string]string) bool {
	changed := false
	for name, attribute := range body.Attributes() {
		tokens := attribute.Expr().BuildTokens(nil)
		attributeChanged := false
		for i := 0; i+2 < len(tokens); i++ {
			if tokens[i].Type != hclsyntax.TokenIdent || tokens[i+1].Type != hclsyntax.TokenDot || tokens[i+2].Type != hclsyntax.TokenIdent {
				continue
			}
			address := string(tokens[i].Bytes) + "." + string(tokens[i+2].Bytes)
			if newLabel, ok := renames[address]; ok {
				tokens[i+2].Bytes = []byte(newLabel)
				attributeChanged = true
			}
		}
		if attributeChanged {
			body.SetAttributeRaw(name, tokens)
			changed = true
		}
	}

	for _, block := range body.Blocks() {
		if rewriteResourceReferences(block.Body(), renames) {
			changed = true
		}
	}
	return changed
}
