package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const fontMatterFormatString = `
# Terraform exporter for SAP BTP

This document explains the syntax and parameters for the various Terraform exporter for SAP BTP commands.

`

const directoryMode fs.FileMode = 0755

func generateMarkdown(rootCmd *cobra.Command) {
	color.NoColor = true
	fmt.Println("Generating markdown documentation")

	basename := strings.ReplaceAll(rootCmd.CommandPath(), " ", "_") + ".md"
	filename := filepath.Join("./docs", basename)

	if err := os.MkdirAll(filepath.Dir(filename), directoryMode); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	docFile, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	defer docFile.Close()

	// Write front-matter to the file:
	if _, err := docFile.WriteString(fontMatterFormatString); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	if err := genMarkdownFile(docFile, rootCmd); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Generated documentation to %v", filename)
}

// addCodeFences adds Markdown code fences (i.e. ```) to example commands listed in help
// text. An example command is a line which begins with a tab character and a dollar sign
// (which signifies the terminal prompt). Blocks of example commands are preceded and terminated
// by whitespace only lines.
func addCodeFencesToSampleCommands(s string) string {
	lines := strings.Split(s, "\n")
	newLines := []string{}

	inBlock := false
	for idx, line := range lines {
		// blank lines cause possible state changes...
		if strings.TrimSpace(line) == "" {
			if inBlock {
				inBlock = false
				newLines = append(newLines, "```")
				newLines = append(newLines, line)
			} else if !inBlock && idx+1 < len(lines) && strings.HasPrefix(lines[idx+1], "\t$") {
				inBlock = true
				newLines = append(newLines, line)
				newLines = append(newLines, "```bash")
			} else {
				newLines = append(newLines, line)
			}
		} else {
			if inBlock && strings.HasPrefix(line, "\t$") {
				line = formatCommandLine(line)
			}
			newLines = append(newLines, line)
		}
	}
	if inBlock {
		newLines = append(newLines, "```")
	}

	return strings.Join(newLines, "\n")
}

var precedingDollarRegexp = regexp.MustCompile(`^([\s]*)\$ (.*)$`)

func formatCommandLine(line string) string {
	return precedingDollarRegexp.ReplaceAllString(line, "$1$2")
}

// Adjusted GenMarkdownTree from spf13/cobra/docs@v1.3.0 package:
//
//   - Emit one help text for all commands into the
//     same unified writer (so they all appear in the same file)
//
//   - Fix the markdown links to refer to anchors
//     in the current file instead of separate files on disk.
func genMarkdownFile(w io.Writer, cmd *cobra.Command) error {
	linkMapper := func(s string) string {
		commandName := strings.TrimSuffix(s, ".md")
		return "#" + strings.ReplaceAll(commandName, "_", "-")
	}

	if err := genMarkdownCustom(cmd, w, linkMapper); err != nil {
		return err
	}

	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		if err := genMarkdownFile(w, c); err != nil {
			return err
		}
	}

	return nil
}

var linkRegexp = regexp.MustCompile(`(https://[^ ]*)`)

func convertLinksToMarkdown(s string) string {
	return linkRegexp.ReplaceAllStringFunc(s, func(link string) string {
		if strings.HasSuffix(link, ".") {
			link = link[:len(link)-1]
			return fmt.Sprintf("[%s](%s).", link, link)
		} else {
			return fmt.Sprintf("[%s](%s)", link, link)
		}
	})
}

// Adjusted `GetMarkdownCustom` from the spf13/cobra/docs@v1.3.0 package:
//
//   - No link to the parent command in the "See also" section when the parent command
//     is itself the root command
//
//   - Add a "Back to top" link at the end of every "See also" section that links back to the root
//     command.
//
//   - Use addCodeFencesToSampleCommands to add code fences to the long help where needed.
//
//   - Format URLs as markdown links (the text of the link is the URL).
func genMarkdownCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString("## " + name + "\n\n")
	buf.WriteString(cmd.Short + "\n\n")
	if len(cmd.Long) > 0 {
		buf.WriteString("### Synopsis\n\n")
		buf.WriteString(convertLinksToMarkdown(addCodeFencesToSampleCommands(cmd.Long)) + "\n\n")
	}

	if cmd.Runnable() {
		fmt.Fprintf(buf, "```bash\n%s\n```\n\n", cmd.UseLine())
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("### Examples\n\n```bash\n")
		lines := strings.Split(cmd.Example, "\n")
		for _, line := range lines {
			buf.WriteString(formatCommandLine(line) + "\n")
		}
		buf.WriteString("```\n\n")
	}

	if err := printOptions(buf, cmd, name); err != nil {
		return err
	}
	if hasSeeAlso(cmd) {
		buf.WriteString("### See also\n\n")

		if cmd.HasParent() {
			parent := cmd.Parent()

			// Write a link to the parent, assuming that it is not the root command
			if parent != cmd.Root() {
				pname := parent.CommandPath()
				link := pname + ".md"
				link = strings.ReplaceAll(link, " ", "_")
				fmt.Fprintf(buf, "* [%s](%s): %s\n", pname, linkHandler(link), parent.Short)
			}
			cmd.VisitParents(func(c *cobra.Command) {
				if c.DisableAutoGenTag {
					cmd.DisableAutoGenTag = c.DisableAutoGenTag
				}
			})
		}

		children := cmd.Commands()
		sort.Sort(byName(children))

		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			cname := name + " " + child.Name()
			link := cname + ".md"
			link = strings.ReplaceAll(link, " ", "_")
			fmt.Fprintf(buf, "* [%s](%s): %s\n", cname, linkHandler(link), child.Short)
		}

		// for child commands, write a link back to the root command with the text "Back to top".
		if cmd.HasParent() {
			root := cmd.Root()
			cname := root.Name()
			link := cname + ".md"
			link = strings.ReplaceAll(link, " ", "_")
			fmt.Fprintf(buf, "* [Back to top](%s)\n", linkHandler(link))
		}

		buf.WriteString("\n")
	}
	if !cmd.DisableAutoGenTag {
		buf.WriteString("###### Auto generated by spf13/cobra on " + time.Now().Format("2-Jan-2006") + "\n")
	}
	_, err := buf.WriteTo(w)
	return err
}

// No changes ot printOptions from spf13/cobra/docs@v1.3.0
func printOptions(buf *bytes.Buffer, cmd *cobra.Command, _ string) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	if flags.HasAvailableFlags() {
		buf.WriteString("### Options\n\n```azdeveloper\n")
		flags.PrintDefaults()
		buf.WriteString("```\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)
	if parentFlags.HasAvailableFlags() {
		buf.WriteString("### Options inherited from parent commands\n\n```azdeveloper\n")
		parentFlags.PrintDefaults()
		buf.WriteString("```\n\n")
	}
	return nil
}

// No changes to hasSeeAlso from spf13/cobra/docs@v1.3.0
func hasSeeAlso(cmd *cobra.Command) bool {
	if cmd.HasParent() {
		return true
	}
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		return true
	}
	return false
}

// No changes to byName from spf13/cobra/docs@v1.3.0
type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
