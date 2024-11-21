package tfutils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/go-git/go-git/v5"
	bf "github.com/russross/blackfriday/v2"
)

const (
	sectionOther               = 0
	sectionExampleUsage        = 1
	sectionArgsReference       = 2
	sectionAttributesReference = 3
	sectionFrontMatter         = 4
	sectionImports             = 5
)

// DocKind indicates what kind of entity's documentation is being requested.
type DocKind string

// argumentDocs contains the documentation metadata for an argument of the resource.
type argumentDocs struct {
	// The description for this argument.
	description string

	// (Optional) The names and descriptions for each argument of this argument.
	arguments map[string]string

	isNested bool
}

// EntityDocs represents the documentation for a resource or datasource as extracted from TF markdown.
type EntityDocs struct {
	// Description is the description of the resource
	Description string
	Arguments   map[string]*argumentDocs
	Attributes  map[string]string
	Import      string
}

var repoPaths sync.Map

func getRepositoryPath(githubHost, organization, provider, version string) (string, error) {
	relativePath := fmt.Sprintf("%s/%s/terraform-provider-%s", githubHost, organization, provider)
	if version != "" {
		relativePath = fmt.Sprintf("%s@%s", relativePath, version)
	}

	if path, ok := repoPaths.Load(relativePath); ok {
		return path.(string), nil
	}

	gitRepoDownloadPath := TmpFolder + "/" + relativePath

	_, er := git.PlainClone(gitRepoDownloadPath, false, &git.CloneOptions{
		URL: "https://github.com/SAP/terraform-provider-btp.git",
	})

	if er != nil {
		return "", fmt.Errorf("error cloning git repo: %w", er)
	}

	repoPaths.Store(relativePath, gitRepoDownloadPath)

	return gitRepoDownloadPath, nil
}

func getDocsPath(repo string, kind DocKind) string {
	newDocsExist := checkIfNewDocsExist(repo)

	if !newDocsExist {
		kindString := string([]rune(kind)[0])
		return filepath.Join(repo, "website", "docs", kindString)
	}

	kindString := string(kind)
	return filepath.Join(repo, "docs", kindString)
}

func checkIfNewDocsExist(repo string) bool {
	newDocsPath := filepath.Join(repo, "docs", "resources")
	_, err := os.Stat(newDocsPath)
	return !os.IsNotExist(err)
}

// readMarkdown searches all possible locations for the markdown content
func readMarkdown(repository string, kind DocKind, markdownName string) ([]byte, string, bool, error) {
	locationPrefix := getDocsPath(repository, kind)

	location := filepath.Join(locationPrefix, markdownName)
	markdownBytes, err := os.ReadFile(location)
	if err == nil {
		return markdownBytes, markdownName, true, nil
	}
	return nil, "", false, err
}

func getMarkdownDetails(org string, provider string, resourcePrefix string,
	kind DocKind, rawName string, providerVersion string, githost string,
) ([]byte, string, bool, error) {

	repoPath, err := getRepositoryPath(githost, org, provider, providerVersion)
	if err != nil {
		return nil, "", false, fmt.Errorf("error in getting repository path : %v", err)
	}

	markdownName := strings.TrimPrefix(rawName, resourcePrefix+"_") + ".md"

	markdownBytes, markdownFileName, found, err := readMarkdown(repoPath, kind, markdownName)
	if !found {
		return nil, "", false, err
	}

	return markdownBytes, markdownFileName, true, nil
}

type tfMarkdownParser struct {
	kind             DocKind
	markdownFileName string
	rawname          string
	ret              EntityDocs
}

// splitGroupLines splits and groups a string, s, by a given separator, sep.
func splitGroupLines(s, sep string) [][]string {
	return grpLines(strings.Split(s, "\n"), sep)
}

// grpLines take a slice of strings, lines, and returns a nested slice of strings. When groupLines encounters a line
// that in the input that starts with the supplied string sep, it will begin a new entry in the outer slice.
func grpLines(lines []string, sep string) [][]string {
	var buffer []string
	var sections [][]string
	for _, line := range lines {
		if strings.Index(line, sep) == 0 {
			sections = append(sections, buffer)
			buffer = []string{}
		}
		buffer = append(buffer, line)
	}
	if len(buffer) > 0 {
		sections = append(sections, buffer)
	}
	return sections
}

func parseDoc(text string) *bf.Node {
	mdProc := bf.New(bf.WithExtensions(bf.FencedCode))
	return mdProc.Parse([]byte(text))
}
func parseNode(text string) *bf.Node {
	return parseDoc(text).FirstChild
}

type paramFlags int

type parameter struct {
	name     string
	desc     string
	typeDecl string
}

type nestedSchema struct {
	longName string
	linkID   *string
	optional []parameter
	required []parameter
	readonly []parameter
}

type topLevelSchema struct {
	optional       []parameter
	required       []parameter
	readonly       []parameter
	nestedSchemata []nestedSchema
}

const (
	required paramFlags = iota
	optional
	readonly
)

var optionalPattern = regexp.MustCompile("(?i)^optional[:]?$")
var requiredPattern = regexp.MustCompile("(?i)^required[:]?$")
var readonlyPattern = regexp.MustCompile("(?i)^read-only[:]?$")

func parseParameterFlagLiteral(text string) *paramFlags {
	if optionalPattern.MatchString(text) {
		return paramFlagsPointr(optional)
	}
	if requiredPattern.MatchString(text) {
		return paramFlagsPointr(required)
	}
	if readonlyPattern.MatchString(text) {
		return paramFlagsPointr(readonly)
	}
	return nil
}

func paramFlagsPointr(flags paramFlags) *paramFlags {
	result := new(paramFlags)
	*result = flags
	return result
}

func parseParameterList(node *bf.Node, ingestNode func(node *bf.Node)) (*[]parameter, error) {
	var out []parameter
	if node == nil || node.Type != bf.List {
		return &out, nil
	}
	item := node.FirstChild
	for item != nil {
		if item.Type != bf.Item {
			return nil, fmt.Errorf("expected an Item")
		}
		param, err := parseParam(item)
		if err != nil {
			return nil, err
		}
		if param == nil {
			return nil, fmt.Errorf("expected a parameter, got %v", prettyPrint(item))
		}
		out = append(out, *param)
		item = item.Next
	}
	defer ingestNode(node)
	return &out, nil
}

// Used for debugging blackfriday parse trees by visualizing them.
func prettyPrint(node *bf.Node) string {
	if node == nil {
		return "nil"
	}
	bytes, err := json.MarshalIndent(treeify(node), "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func treeify(node *bf.Node) interface{} {
	if node == nil {
		return "nil"
	}
	if node.Type == bf.Text {
		return string(node.Literal)
	}
	var result []interface{}
	result = append(result, fmt.Sprintf("[%s]", node.Type))

	c := node.FirstChild
	for c != nil {
		result = append(result, treeify(c))
		c = c.Next
	}

	if node.Literal != nil {
		result = append(result, string(node.Literal))
	}
	return result
}

var seeBelowNestedSchemaPattern = regexp.MustCompile(`[(]see \[below for nested schema\][(][^)]*[)][)]`)

func parseParam(node *bf.Node) (*parameter, error) {
	if node == nil || node.Type != bf.Item {
		return nil, nil
	}
	paragraph := node.FirstChild
	if paragraph == nil || paragraph.Type != bf.Paragraph || paragraph.Next != nil {
		return nil, nil
	}
	emptyTextNode := paragraph.FirstChild
	if emptyTextNode == nil || emptyTextNode.Type != bf.Text || len(emptyTextNode.Literal) > 0 {
		return nil, nil
	}
	strongOrCode := emptyTextNode.Next
	isStrong := strongOrCode != nil && strongOrCode.Type == bf.Strong
	isCode := strongOrCode != nil && strongOrCode.Type == bf.Code
	if !isStrong && !isCode {
		return nil, nil
	}
	var parameterName string
	if isStrong {
		parsed, err := parseTextSequenceFromNode(strongOrCode.FirstChild, false)
		if err != nil {
			return nil, err
		}
		parameterName = parsed
	} else {
		parameterName = string(strongOrCode.Literal)
	}

	paramDescription, err := parseTextSequenceFromNode(strongOrCode.Next, true)
	if err != nil {
		return nil, err
	}
	return parseParameterFromDesc(parameterName, cleanDescription(paramDescription)), nil
}

var typeDeclarationPattern = regexp.MustCompile(`^\s*[(]([^[)]+)[)]\s*`)

func parseParameterFromDesc(paramName string, desc string) *parameter {
	if typeDeclarationPattern.MatchString(desc) {
		typeDecl := typeDeclarationPattern.FindStringSubmatch(desc)[1]
		desc = typeDeclarationPattern.ReplaceAllString(desc, "")

		return &parameter{
			name:     paramName,
			desc:     desc,
			typeDecl: typeDecl,
		}
	}
	return &parameter{
		name: paramName,
		desc: desc,
	}
}

func cleanDescription(description string) string {
	description = seeBelowNestedSchemaPattern.ReplaceAllString(description, "")
	return strings.TrimSpace(description)
}

func parseTextSequenceFromNode(startNode *bf.Node, useStarsForStrongAndEmph bool) (string, error) {
	var parseError error
	textBuilder := strings.Builder{}
	currentNode := startNode
	for currentNode != nil {
		currentNode.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
			switch node.Type {
			case bf.Text:
				textBuilder.WriteString(string(node.Literal))
			case bf.Code:
				textBuilder.WriteString("`")
				textBuilder.WriteString(string(node.Literal))
				textBuilder.WriteString("`")
			case bf.Link:
				if entering {
					textBuilder.WriteString("[")
				} else {
					textBuilder.WriteString("](")
					textBuilder.WriteString(string(node.Destination))
					textBuilder.WriteString(")")
				}
			case bf.Strong:
				if useStarsForStrongAndEmph {
					textBuilder.WriteString("**")
				} else {
					textBuilder.WriteString("__")
				}
			case bf.Emph:
				if useStarsForStrongAndEmph {
					textBuilder.WriteString("*")
				} else {
					textBuilder.WriteString("_")
				}
			case bf.HTMLSpan:
				textBuilder.WriteString(`\n\n`)
			default:
				parseError = fmt.Errorf("found a tag it cannot yet render back to Markdown: %s",
					prettyPrint(node))
				return bf.Terminate
			}
			return bf.GoToNext
		})
		currentNode = currentNode.Next
	}
	return textBuilder.String(), parseError
}
func parseParamSec(node *bf.Node, ingestNode func(node *bf.Node)) (paramFlags, *[]parameter, *bf.Node, error) {
	if node != nil && (node.Type == bf.Paragraph || node.Type == bf.Heading) {
		sectionLabel := node.FirstChild
		if sectionLabel != nil && sectionLabel.Type == bf.Text && sectionLabel.Next == nil {
			flags := parseParameterFlagLiteral(string(sectionLabel.Literal))
			if flags == nil {
				return -1, nil, nil, nil
			}
			parameterList, err := parseParameterList(node.Next, ingestNode)
			if err != nil {
				return -1, nil, nil, err
			}

			if parameterList == nil || node.Next == nil {
				return -1, nil, nil, fmt.Errorf("expected a parameter list, got %s", prettyPrint(node.Next))
			}

			defer ingestNode(node)
			return *flags, parameterList, node.Next.Next, nil
		}
	}
	return -1, nil, nil, nil
}

func parseTopLevelSchema(node *bf.Node, ingestNode func(node *bf.Node)) (*topLevelSchema, error) {
	if ingestNode == nil {
		ingestNode = func(node *bf.Node) {}
	}
	if node == nil || node.Type != bf.Heading {
		return nil, nil
	}
	label := node.FirstChild
	if label == nil || label.Type != bf.Text || string(label.Literal) != "Schema" {
		return nil, nil
	}
	tls := &topLevelSchema{}
	currentNode := node.Next
	for currentNode != nil {
		flags, par, next, err := parseParamSec(currentNode, ingestNode)
		if err != nil {
			return nil, err
		}
		if par != nil {
			switch flags {
			case optional:
				tls.optional = *par
			case required:
				tls.required = *par
			case readonly:
				tls.readonly = *par
			}
			currentNode = next
		} else {
			break
		}
	}

	var nested []nestedSchema
	currentNode = node.Next
	for currentNode != nil {
		nestedSchema, err := parseNestedSchema(currentNode, ingestNode)
		if err != nil {
			return nil, err
		}
		if nestedSchema != nil {
			nested = append(nested, *nestedSchema)
		}
		currentNode = currentNode.Next
	}

	tls.nestedSchemata = nested

	ingestNode(node)
	return tls, nil
}

var preamblePattern = regexp.MustCompile("^[<]a id=[\"]([^\"]+)[\"][>]$")

func parsePreamble(node *bf.Node, processNode func(node *bf.Node)) *string {
	if node == nil || node.Type != bf.Paragraph {
		return nil
	}
	firstTextNode := node.FirstChild
	if firstTextNode == nil || firstTextNode.Type != bf.Text || len(firstTextNode.Literal) > 0 {
		return nil
	}
	htmlSpanNode := firstTextNode.Next
	if htmlSpanNode == nil || htmlSpanNode.Type != bf.HTMLSpan {
		return nil
	}
	secondTextNode := htmlSpanNode.Next
	if secondTextNode == nil || secondTextNode.Type != bf.Text || len(secondTextNode.Literal) > 0 {
		return nil
	}
	closingHtmlSpanNode := secondTextNode.Next
	if closingHtmlSpanNode == nil || closingHtmlSpanNode.Type != bf.HTMLSpan || string(closingHtmlSpanNode.Literal) != "</a>" {
		return nil
	}
	preambleMatches := preamblePattern.FindStringSubmatch(string(htmlSpanNode.Literal))
	if len(preambleMatches) > 1 {
		defer processNode(node)
		return &preambleMatches[1]
	}
	return nil
}

func parseNestedSchema(node *bf.Node, ingestNode func(node *bf.Node)) (*nestedSchema, error) {
	if ingestNode == nil {
		ingestNode = func(node *bf.Node) {}
	}

	if node.Prev != nil && parsePreamble(node.Prev, func(x *bf.Node) {}) != nil {
		return nil, nil
	}

	linkID := parsePreamble(node, ingestNode)
	if linkID != nil {
		node = node.Next
	}

	if node == nil || node.Type != bf.Heading {
		return nil, nil
	}

	label := node.FirstChild

	if label == nil || label.Type != bf.Text || string(label.Literal) != "Nested Schema for " {
		return nil, nil
	}

	code := label.Next
	if code == nil || code.Type != bf.Code {
		return nil, fmt.Errorf("expected a Code block, got %s", prettyPrint(code))
	}

	ns := &nestedSchema{
		longName: string(code.Literal),
		linkID:   linkID,
	}

	currentNode := node.Next
	for {
		flags, par, next, err := parseParamSec(currentNode, ingestNode)
		if err != nil {
			return nil, err
		}
		if par != nil {
			switch flags {
			case optional:
				ns.optional = *par
			case required:
				ns.required = *par
			case readonly:
				ns.readonly = *par
			}
			currentNode = next
		} else {
			break
		}
	}

	defer ingestNode(node)
	return ns, nil
}

func (ns *topLevelSchema) requiredParameters() []parameter {
	return ns.required
}

func (ns *nestedSchema) allParameters() []parameter {
	return append(append(ns.optional, ns.required...), ns.readonly...)
}

func parseTopLevelSchemaIntoDocs(
	accumulatedDocs *EntityDocs,
	topLevelSchema *topLevelSchema,
) {
	//for _, param := range topLevelSchema.allParameters() {
	for _, param := range topLevelSchema.requiredParameters() {
		oldDesc, haveOldDesc := accumulatedDocs.Attributes[param.name]
		if haveOldDesc && oldDesc != param.desc {
			fmt.Print("\r\n")
			log.Fatalf("Description conflict for top-level attribute %s; candidates are `%s` and `%s`",
				param.name,
				oldDesc,
				param.desc)
		}
		accumulatedDocs.Attributes[param.name] = param.desc
	}

	for _, ns := range topLevelSchema.nestedSchemata {
		nestedSchema := ns // this stops implicit memory addressing
		parseNestedSchemaIntoDocuments(accumulatedDocs, &nestedSchema)
	}
}

func parseNestedSchemaIntoDocuments(
	accumulatedDocs *EntityDocs,
	nestedSchema *nestedSchema,
) {

	args, _ := accumulatedDocs.getOrCreateArgumentDocs(nestedSchema.longName)
	args.isNested = true

	for _, param := range nestedSchema.allParameters() {
		oldDescription, hasAlready := args.arguments[param.name]
		if hasAlready && oldDescription != param.desc {
			fmt.Print("\r\n")
			log.Fatalf("Description conflict for param %s from %s; candidates are `%s` and `%s`",
				param.name,
				nestedSchema.longName,
				oldDescription,
				param.desc)
		}
		args.arguments[param.name] = param.desc
		fullParameterName := fmt.Sprintf("%s.%s", nestedSchema.longName, param.name)
		paramArgs, created := accumulatedDocs.getOrCreateArgumentDocs(fullParameterName)
		if !created && paramArgs.description != param.desc {
			fmt.Print("\r\n")
			log.Fatalf("Description conflict for param %s; candidates are `%s` and `%s`",
				fullParameterName,
				paramArgs.description,
				param.desc)
		}
		paramArgs.isNested = true
		paramArgs.description = param.desc
	}
}

func (ed *EntityDocs) getOrCreateArgumentDocs(argumentName string) (*argumentDocs, bool) {
	if ed.Arguments == nil {
		ed.Arguments = make(map[string]*argumentDocs)
	}
	var created bool
	arguments, has := ed.Arguments[argumentName]
	if !has {
		arguments = &argumentDocs{arguments: make(map[string]string)}
		ed.Arguments[argumentName] = arguments
		created = true
	}
	return arguments, created
}

func (p *tfMarkdownParser) parseSchemaHavingNestedSections(subsection []string) {
	node := parseNode(strings.Join(subsection, "\n"))
	topLevelSchema, err := parseTopLevelSchema(node, nil)
	if err != nil {
		log.Printf("error: Failure in parsing resource name: %s, subsection: %s", p.rawname, subsection[0])
		return
	}
	if topLevelSchema == nil {
		log.Print("Failed to parse top-level Schema section")
		return
	}
	parseTopLevelSchemaIntoDocs(&p.ret, topLevelSchema)
}

// isBlankString returns true if the line is all whitespace.
func isBlankString(line string) bool {
	return strings.TrimSpace(line) == ""
}

func (p *tfMarkdownParser) reformatSubsections(lines []string) ([]string, bool, bool) {
	var result []string
	hasExamples, isEmpty := false, true

	var inOICSButton bool // True if we are removing an "Open in Cloud Shell" button.
	for i, line := range lines {
		if inOICSButton {
			if strings.Index(lines[i], "</div>") == 0 {
				inOICSButton = false
			}
		} else {
			if strings.Index(line, "<div") == 0 && strings.Contains(line, "oics-button") {
				inOICSButton = true
			} else {
				if strings.Index(line, "```") == 0 {
					hasExamples = true
				} else if !isBlankString(line) {
					isEmpty = false
				}

				result = append(result, line)
			}
		}
	}

	return result, hasExamples, isEmpty
}

func (p *tfMarkdownParser) parseMarkdownSections(h2Section []string) error {
	if len(h2Section) == 0 {
		fmt.Print("\r\n")
		log.Fatalf("Unparseable H2 doc section for %v; consider overriding doc source location", p.rawname)
		return nil
	}

	header := trimHeader(h2Section[0])

	sectionType := sectionOther

	switch header {
	case "Timeout", "Timeouts", "User Project Override", "User Project Overrides":
		log.Printf("Ignoring section [%v] for [%v]", header, p.rawname)
		return nil
	case "Example Usage":
		sectionType = sectionExampleUsage
	case "Arguments Reference", "Argument Reference", "Argument reference", "Nested Blocks", "Nested blocks":
		sectionType = sectionArgsReference
	case "Attributes Reference", "Attribute Reference", "Attribute reference":
		sectionType = sectionAttributesReference
	case "Import", "Imports":
		sectionType = sectionImports
	case "---":
		sectionType = sectionFrontMatter
	case "Schema":
		p.parseSchemaHavingNestedSections(h2Section)
		return nil
	}

	var wrHeader bool
	for _, h3Section := range grpLines(h2Section[1:], "### ") {
		if len(h3Section) == 0 {
			log.Printf("Empty or unparseable H3 doc section for %v",
				p.rawname)
			continue
		}

		reformattedSection, _, isEmpty := p.reformatSubsections(h3Section)
		if isEmpty {
			continue
		}

		switch sectionType {
		case sectionArgsReference:
			parseArgumentReferenceSection(reformattedSection, &p.ret)
		case sectionAttributesReference:
			parseAttrReferenceSection(reformattedSection, &p.ret)
		case sectionFrontMatter:
			p.parseIntro(reformattedSection)
		case sectionImports:
			p.parseImport(reformattedSection)
		default:
			_, isArgument := p.ret.Arguments[header]
			if isArgument || strings.HasSuffix(header, "Configuration Block") {
				parseArgumentReferenceSection(reformattedSection, &p.ret)
				continue
			}

			if !wrHeader {
				p.ret.Description += fmt.Sprintf("## %s\n", header)
				wrHeader = true
				if !isBlankString(reformattedSection[0]) {
					p.ret.Description += "\n"
				}
			}
			p.ret.Description += strings.Join(reformattedSection, "\n") + "\n"
		}
	}

	return nil
}

func trimHeader(header string) string {
	if strings.HasPrefix(header, "## ") {
		return header[3:]
	}
	return header
}

func (p *tfMarkdownParser) parseIntro(subsection []string) {
	// The header of the MarkDown will have two "---"s paired up to delineate the header. Skip this.
	var endHeaderFound bool
	for len(subsection) > 0 {
		current := subsection[0]
		subsection = subsection[1:]
		if current == "---" {
			endHeaderFound = true
			break
		}
	}
	if !endHeaderFound {
		log.Print("", "Expected to pair --- begin/end for resource %v's header", p.rawname)
	}

	previousBlank := true
	var h1ResourceFound bool
	for _, line := range subsection {
		if strings.Index(line, "# ") == 0 {
			h1ResourceFound = true
			previousBlank = true
		} else if !isBlankString(line) || !previousBlank {
			p.ret.Description += line + "\n"
			previousBlank = false
		} else if isBlankString(line) {
			previousBlank = true
		}
	}
	if !h1ResourceFound {
		log.Printf("Expected an H1 in markdown for resource %v", p.rawname)
	}
}

//nolint:lll
var (
	argumentBulletRegexp = regexp.MustCompile(
		"^\\s*[*+-]\\s*`([a-zA-z0-9_]*)`\\s*(\\([a-zA-Z]*\\)\\s*)?\\s*[:–-]?\\s*(\\([^\\)]*\\)[-\\s]*)?(.*)",
	)

	attributeBulletRegexp = regexp.MustCompile(
		"^\\s*[*+-]\\s*`([a-zA-z0-9_]*)`\\s*[:–-]?\\s*(.*)",
	)
)

var nestedObjectRegexps = []*regexp.Regexp{
	regexp.MustCompile("`([a-z_]+)`.*following"),

	regexp.MustCompile("(?i)## ([a-z_]+).* argument reference"),

	regexp.MustCompile("###+ ([a-z_]+).*"),

	regexp.MustCompile("###+ `([a-z_]+).*`"),

	regexp.MustCompile("`([a-zA-Z_.\\[\\]]+)`.*supports:"),

	regexp.MustCompile("###+ ([a-zA-Z_ ]+).*"),
}

// parseArgumentFromMarkdownLine takes a line of Markdown and attempts to parse it for a Terraform argument and its
// description
func parseArgumentFromMarkdownLine(line string) (string, string, bool) {
	matches := argumentBulletRegexp.FindStringSubmatch(line)

	if len(matches) > 4 {
		return matches[1], matches[4], true
	}

	return "", "", false
}

func getNestedBlockNames(line string) string {
	nestedBlockName := ""
	for _, regex := range nestedObjectRegexps {
		matchedString := regex.FindStringSubmatch(line)
		if len(matchedString) >= 2 {
			nestedBlockName = strings.ToLower(matchedString[1])
			nestedBlockName = strings.Replace(nestedBlockName, " ", "_", -1)
			nestedBlockName = strings.TrimSuffix(nestedBlockName, "[]")
			parts := strings.Split(nestedBlockName, ".")
			nestedBlockName = parts[len(parts)-1]
			break
		}
	}

	return nestedBlockName
}

var genericNestedRegexp = regexp.MustCompile("supports? the following:")

func parseArgumentReferenceSection(subsec []string, entity *EntityDocs) {
	var lastArgument, nestedBlock string

	addHeading := func(headingName string, headingDescription string, line string) {

		initializeArgumentDocs := func(key string) {
			if entity.Arguments[key] == nil {
				entity.Arguments[key] = &argumentDocs{
					arguments: make(map[string]string),
				}
			} else if entity.Arguments[key].arguments == nil {
				entity.Arguments[key].arguments = make(map[string]string)
			}
		}

		if nestedBlock != "" {
			initializeArgumentDocs(nestedBlock)
			entity.Arguments[nestedBlock].arguments[headingName] = headingDescription

			if entity.Arguments[headingName] == nil {
				entity.Arguments[headingName] = &argumentDocs{
					description: headingDescription,
					isNested:    true,
				}
			}
		} else {
			if genericNestedRegexp.MatchString(line) {
				return
			}
			entity.Arguments[headingName] = &argumentDocs{description: headingDescription}
		}
	}

	extendAvailableHeading := func(line string) {
		line = "\n" + strings.TrimSpace(line)
		if nestedBlock != "" {
			entity.Arguments[nestedBlock].arguments[lastArgument] += line

			if entity.Arguments[lastArgument].isNested {
				entity.Arguments[lastArgument].description += line
			}
		} else {
			if genericNestedRegexp.MatchString(line) {
				lastArgument = ""
				nestedBlock = ""
				return
			}
			entity.Arguments[lastArgument].description += line
		}
	}

	var hadSpace bool
	for _, line := range subsec {
		if name, desc, matchFound := parseArgumentFromMarkdownLine(line); matchFound {
			addHeading(name, desc, line)
			lastArgument = name
		} else if strings.TrimSpace(line) == "---" {
			lastArgument = ""
		} else if nestedBlockCurrentLine := getNestedBlockNames(line); hadSpace && nestedBlockCurrentLine != "" {
			nestedBlock = nestedBlockCurrentLine
			lastArgument = ""
		} else if !isBlankString(line) && lastArgument != "" {
			extendAvailableHeading(line)
		} else if nestedBlockCurrentLine := getNestedBlockNames(line); nestedBlockCurrentLine != "" {
			nestedBlock = nestedBlockCurrentLine
			lastArgument = ""
		} else if lastArgument != "" {
			extendAvailableHeading(line)
		}
		hadSpace = isBlankString(line)
	}

	for _, v := range entity.Arguments {
		v.description = strings.TrimRightFunc(v.description, unicode.IsSpace)
		for k, d := range v.arguments {
			v.arguments[k] = strings.TrimRightFunc(d, unicode.IsSpace)
		}
	}
}

func (p *tfMarkdownParser) parseImport(importLines []string) {
	var token string

	defer func() {
		stringContainsTerraform := strings.Contains(strings.ToLower(p.ret.Import), "terraform")
		if stringContainsTerraform {
			message := fmt.Sprintf("parseImport %q should not render the string"+
				" 'terraform' in its emitted markdown.\n"+
				"**Input**:\n%s\n\n**Rendered**:\n%s\n\n",
				token, strings.Join(importLines, "\n"), p.ret.Import)
			log.Print(message)
		}
	}()

	var importString []string
	for _, line := range importLines {
		if strings.Contains(line, "**NOTE:") || strings.Contains(line, "**Please Note:") ||
			strings.Contains(line, "**Note:**") {
			continue
		}
		if strings.Contains(line, "Import is supported using the following syntax") {
			continue
		}

		line = strings.Replace(line, "```shell", "", -1)
		line = strings.Replace(line, "```sh", "", -1)
		line = strings.Replace(line, "```", "", -1)

		if strings.Contains(line, "# terraform import") {
			line = strings.Replace(line, "$ ", "", -1)
			line = strings.Replace(line, "# terraform import ", "", -1)

			parts := strings.Split(line, " ")
			importTemplate := `import {
				to = %s
				id = "%s"
			  }`
			importString = append(importString, fmt.Sprintf(importTemplate, parts[0], parts[1]))

		}
	}

	p.ret.Import = p.ret.Import + strings.Join(importString, " ")
}

func parseAttrReferenceSection(attributeLines []string, entity *EntityDocs) {
	var lastMatchedAttribute string
	for _, line := range attributeLines {
		matches := attributeBulletRegexp.FindStringSubmatch(line)
		if len(matches) >= 2 {
			entity.Attributes[matches[1]] = matches[2]
			lastMatchedAttribute = matches[1]
		} else if !isBlankString(line) && lastMatchedAttribute != "" {
			entity.Attributes[lastMatchedAttribute] += "\n" + strings.TrimSpace(line)
		} else {
			lastMatchedAttribute = ""
		}
	}
}

var markdownPageReferenceLink = regexp.MustCompile(`\[[1-9]+\]: /docs/providers(?:/[a-z1-9_]+)+\.[a-z]+`)
var (
	// Match a [markdown](link)
	//markdownLink = regexp.MustCompile(`\[([^\]]*)\]\(([^\)]*)\)`)

	// Match a ```fenced code block```.
	codeBlocks = regexp.MustCompile(`(?ms)\x60\x60\x60[^\n]*?$.*?\x60\x60\x60\s*$`)
)

func reorgenizeText(text string) (string, bool) {

	cleanupText := func(text string) (string, bool) {
		if strings.Contains(text, "Terraform") || strings.Contains(text, "terraform") {
			return "", true
		}
		text = strings.ReplaceAll(text, "-> ", "> ")
		text = strings.ReplaceAll(text, "~> ", "> ")

		text = strings.TrimPrefix(text, "\n(Required)\n")
		text = strings.TrimPrefix(text, "\n(Optional)\n")

		text = markdownPageReferenceLink.ReplaceAllStringFunc(text, func(referenceLink string) string {
			parts := strings.Split(referenceLink, " ")
			return fmt.Sprintf("%s https://www.terraform.io%s", parts[0], parts[1])
		})

		return text, false
	}

	codeBlocks := codeBlocks.FindAllStringIndex(text, -1)

	var parts []string
	start, end := 0, 0
	for _, codeBlock := range codeBlocks {
		end = codeBlock[0]

		clean, elided := cleanupText(text[start:end])
		if elided {
			return "", true
		}
		parts = append(parts, clean)

		start = codeBlock[1]
		parts = append(parts, text[end:start])
	}
	if start != len(text) {
		clean, elided := cleanupText(text[start:])
		if elided {
			return "", true
		}
		parts = append(parts, clean)
	}

	return strings.TrimSpace(strings.Join(parts, "")), false
}

func cleanupDocument(name string, doc EntityDocs) (EntityDocs, bool) {
	hasElidedDoc := false
	cleanedArguments := make(map[string]*argumentDocs, len(doc.Arguments))

	for argKey, argValue := range doc.Arguments {
		cleanedText, elided := reorgenizeText(argValue.description)
		if elided {
			log.Printf("Found <elided> in docs for argument [%v] in [%v].", argKey, name)
			hasElidedDoc = true
		}

		cleanedArguments[argKey] = &argumentDocs{
			description: cleanedText,
			arguments:   make(map[string]string, len(argValue.arguments)),
			isNested:    argValue.isNested,
		}

		for kk, vv := range argValue.arguments {
			cleanedText, elided := reorgenizeText(vv)
			if elided {
				log.Printf("Found <elided> in docs for nested argument [%v] in [%v].", kk, name)
				hasElidedDoc = true
			}
			cleanedArguments[argKey].arguments[kk] = cleanedText
		}
	}

	cleanedAttributes := make(map[string]string, len(doc.Attributes))
	for attrKey, v := range doc.Attributes {
		cleanedText, elided := reorgenizeText(v)
		if elided {
			log.Printf("Found <elided> in docs for attribute [%v] in [%v].", attrKey, name)
			hasElidedDoc = true
		}
		cleanedAttributes[attrKey] = cleanedText
	}

	cleanupText, _ := reorgenizeText(doc.Description)

	var data_source_pattern = regexp.MustCompile("## Example Usage\\n\\n.{3}terraform\n?#?.*\\ndata")
	if data_source_pattern.MatchString(cleanupText) {
		doc.Import = fmt.Sprintf("data \"%s\" \"%s\"{\n  ", name, "all")
		if len(cleanedAttributes) != 0 {
			for attrName, attrValue := range cleanedAttributes {
				doc.Import += fmt.Sprintf("%s = \"%s\" \n", attrName, attrValue)
			}
		} else if name == "cloudfoundry_users" || name == "cloudfoundry_domains" {
			doc.Import += "org = \"The ID of the organization\" \n"
		} else {
			doc.Import += "id = \"The ID of the subaccount\" \n"
		}
		doc.Import += "}\n"
	}

	return EntityDocs{
		Description: cleanupText,
		Arguments:   cleanedArguments,
		Attributes:  cleanedAttributes,
		Import:      doc.Import,
	}, hasElidedDoc
}

func (p *tfMarkdownParser) parse(tfMarkdown []byte) (EntityDocs, error) {
	p.ret = EntityDocs{
		Arguments:  make(map[string]*argumentDocs),
		Attributes: make(map[string]string),
	}
	markdown := string(tfMarkdown)

	// Replace any Windows-style newlines.
	markdown = strings.Replace(markdown, "\r\n", "\n", -1)

	// Replace redundant comment.
	markdown = strings.Replace(markdown, "<!-- schema generated by tfplugindocs -->", "", -1)

	// Split the sections by H2 topics in the Markdown file.
	sections := splitGroupLines(markdown, "## ")

	for _, section := range sections {
		if err := p.parseMarkdownSections(section); err != nil {
			return EntityDocs{}, err
		}
	}

	doc, _ := cleanupDocument(p.rawname, p.ret)

	return doc, nil
}

// parseTFMarkdown takes a TF website markdown doc and extracts a structured representation for use in
// generating doc comments
func parseTFMarkdown(kind DocKind,
	markdown []byte, markdownFileName, rawname string) (EntityDocs, error) {

	p := &tfMarkdownParser{
		kind:             kind,
		markdownFileName: markdownFileName,
		rawname:          rawname,
	}
	return p.parse(markdown)
}

// GetDocsForResource extracts documentation details for the given package from
// TF website documentation markdown content
func GetDocsForResource(org string, provider string, resourcePrefix string, kind DocKind,
	rawname string /* info tfbridge.ResourceOrDataSourceInfo, */, providerModuleVersion string,
	githost string) (EntityDocs, error) {

	markdownBytes, markdownFileName, found, err := getMarkdownDetails(org, provider,
		resourcePrefix, kind, rawname, providerModuleVersion, githost)
	if !found {
		return EntityDocs{}, fmt.Errorf("could not find docs for %v %v. err: %v", kind, rawname, err)
	}

	doc, err := parseTFMarkdown(kind, markdownBytes, markdownFileName, rawname)
	if err != nil {
		return EntityDocs{}, err
	}

	return doc, nil
}
