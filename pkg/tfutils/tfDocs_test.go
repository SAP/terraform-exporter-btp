package tfutils

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseImport(t *testing.T) {
	var tests = []struct {
		input    []string
		expected string
	}{
		{
			input:    readlines(t, "testdata/testImport.md"),
			expected: "import {\n\t\t\t\tto = btp_subaccount.<resource_name>\n\t\t\t\tid = \"<subaccount_id>\"\n\t\t\t  }",
		},
		{
			input:    readlines(t, "testdata/testImportUsingIdentity.md"),
			expected: "import {\n\t\t\t\t\tto = btp_subaccount_subscription.<resource_name>\n\t\t\t\t\tidentity = {\n\t\t\t\t\tsubaccount_id = \"<subaccount_id>\"\napp_name = \"<app_name>\"\nplan_name = \"<plan_name>\"\n\n\t\t\t\t\t}\n\t\t\t\t  }",
		},
	}

	for _, tt := range tests {
		parser := tfMarkdownParser{}
		parser.parseImport(tt.input)
		actual := parser.ret.Import
		assert.Equal(t, tt.expected, actual)
	}
}

func TestParseArgumentFromMarkdownLine(t *testing.T) {
	//nolint:lll
	tests := []struct {
		input         string
		expectedName  string
		expectedDesc  string
		expectedFound bool
	}{
		{"* `name` (String) A descriptive name of the subaccount.", "name", "A descriptive name of the subaccount.", true},
		{"* `region` (String) The region of the subaccount.", "region", "The region of the subaccount.", true},
		{"  * `subaccount_id` - The id of the subaccount", "subaccount_id", "The id of the subaccount", true},
		{"  * id - The id of the subaccount", "", "", false},
	}

	for _, test := range tests {
		name, desc, isFound := parseArgumentFromMarkdownLine(test.input)
		assert.Equal(t, test.expectedName, name)
		assert.Equal(t, test.expectedDesc, desc)
		assert.Equal(t, test.expectedFound, isFound)
	}
}

func TestParseAttrReferenceSection(t *testing.T) {
	ret := EntityDocs{
		Arguments:  make(map[string]*argumentDocs),
		Attributes: make(map[string]string),
	}
	parseAttrReferenceSection([]string{
		"The following attributes are exported:",
		"",
		"* `id` - The ID of the subaccount",
		"* `name`- The name of the subaccount",
		"* `region` - The region of the subaccount",
	}, &ret)
	assert.Len(t, ret.Attributes, 3)
}

func TestParseTFMarkdownResource(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		kind     DocKind
		rawName  string
		fileName string
	}

	test := func(name string, configure ...func(tc *testCase)) testCase {
		tc := testCase{
			name:     name,
			kind:     "resouces",
			rawName:  "btp_subaccount",
			fileName: "subaccount.md",
		}
		for _, c := range configure {
			c(&tc)
		}
		return tc
	}

	tests := []testCase{
		test("subaccount"),
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.NotZero(t, tt.name)
			input := testFilePath(tt.name, "input_subaccount.md")
			expected := filepath.Join(tt.name, "expected_subaccount.json")
			p := &tfMarkdownParser{
				kind:             tt.kind,
				markdownFileName: tt.fileName,
				rawname:          tt.rawName,
			}

			inputBytes, err := os.ReadFile(input)
			require.NoError(t, err)

			actual, err := p.parse(inputBytes)
			require.NoError(t, err)

			actualBytes, err := json.MarshalIndent(actual, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			compareTestFile(t, expected, string(actualBytes), assert.JSONEq)
		})
	}
}

func TestParseTFMarkdownDataSource(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		kind     DocKind
		rawName  string
		fileName string
	}

	test := func(name string, configure ...func(tc *testCase)) testCase {
		tc := testCase{
			name:     name,
			kind:     "data-sources",
			rawName:  "btp_subaccount",
			fileName: "subaccount.md",
		}
		for _, c := range configure {
			c(&tc)
		}
		return tc
	}

	tests := []testCase{
		test("subaccount"),
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.NotZero(t, tt.name)
			input := testFilePath(tt.name, "input_subaccount_data.md")
			expected := filepath.Join(tt.name, "expected_subaccount_data.json")
			p := &tfMarkdownParser{
				kind:             tt.kind,
				markdownFileName: tt.fileName,
				rawname:          tt.rawName,
			}

			inputBytes, err := os.ReadFile(input)
			require.NoError(t, err)

			actual, err := p.parse(inputBytes)
			require.NoError(t, err)

			actualBytes, err := json.MarshalIndent(actual, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			compareTestFile(t, expected, string(actualBytes), assert.JSONEq)
		})
	}
}

func compareTestFile(
	t *testing.T, path, actual string,
	comp func(t assert.TestingT, expected string, actual string, msgAndArgs ...interface{}) bool,
) {
	comp(t, readTestFile(t, path), actual)
}

func readTestFile(t *testing.T, name string) string {
	bytes, err := os.ReadFile(testFilePath(name))
	if err != nil {
		t.Fatal(err)
	}
	return strings.ReplaceAll(string(bytes), "\r\n", "\n")
}

func readlines(t *testing.T, file string) []string {
	t.Helper()
	f, err := os.Open(file)
	require.NoError(t, err)

	defer func() {
		// Ignore error in test function
		_ = f.Close()
	}()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func testFilePath(path ...string) string {
	return filepath.Join(append([]string{"testdata"}, path...)...)
}

func TestReorganizeText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Creates a role in a directory.",
			expected: "Creates a role in a directory.",
		},
		{
			input:    "Terraform string",
			expected: "",
		},
		{
			input:    "\n(Optional)\nThe role description.",
			expected: "The role description.",
		},
		{
			input:    "id -> The combined unique ID of the role",
			expected: "id > The combined unique ID of the role",
		},
		{
			input:    "\n(Required)\nThe ID of the xsuaa application.",
			expected: "The ID of the xsuaa application.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, elided := reorgenizeText(tt.input)

			assert.Equal(t, tt.expected, got)
			assert.Equalf(t, got == "", elided,
				"We should only see an empty result for non-empty inputs if we have elided text")
		})
	}
}

func TestArgumentRegularExp(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string]*argumentDocs
	}{
		{
			name: "Discovers * bullet descriptions",
			input: []string{
				"* `id` - (Optional) A description of the directory",
				"* `ipv6_address_count`- (Optional) A number of IPv6 addresses to associate with the primary network interface.",
				"* `ipv6_addresses` - (Optional) Specify one or more IPv6 addresses from the range of the subnet to associate with the primary network interface",
				"* `tags` - (Optional) A mapping of tags to assign to the resource.",
			},
			expected: map[string]*argumentDocs{
				"id": {
					description: "A description of the directory",
				},
				"ipv6_address_count": {
					description: "A number of IPv6 addresses to associate with the primary network interface.",
				},
				"ipv6_addresses": {
					description: "Specify one or more IPv6 addresses from the range of the subnet to associate with the primary network interface",
				},
				"tags": {
					description: "A mapping of tags to assign to the resource.",
				},
			},
		},
		{
			name: "Cleans up tabs",
			input: []string{
				"* `node_pool_config` (Input only) The configuration for the node pool. ",
				"       If specified, Dataproc attempts to create a node pool with the specified shape. ",
				"       If one with the same name already exists, it is verified against all specified fields. ",
				"       If a field differs, the virtual cluster creation will fail.",
			},
			expected: map[string]*argumentDocs{
				"node_pool_config": {description: "The configuration for the node pool. \nIf specified, " +
					"Dataproc attempts to create a node pool with the specified shape.\nIf one with the same name " +
					"already exists, it is verified against all specified fields.\nIf a field differs, the virtual " +
					"cluster creation will fail.",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ret := EntityDocs{
				Arguments: make(map[string]*argumentDocs),
			}
			parseArgumentReferenceSection(tt.input, &ret)

			assert.Equal(t, tt.expected, ret.Arguments)
		})
	}
}

func TestGetNestedBlockName(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"The `features` object supports the following:", "features"},
		{"#### result_configuration Argument Reference", "result_configuration"},
		{"### advanced_security_options", "advanced_security_options"},
		{"### `server_side_encryption`", "server_side_encryption"},
		{"### Failover Routing Policy", "failover_routing_policy"},
		{"##### `log_configuration`", "log_configuration"},
		{"## Import", ""},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, getNestedBlockNames(tt.input))
	}
}

func TestParseTextSeq(t *testing.T) {
	turnaround := func(src string) {
		res, err := parseTextSequenceFromNode(parseNode(src).FirstChild, true)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, src, res)
	}

	turnaround("plain")
	turnaround("`code`")
	turnaround("*emph*")
	turnaround("**strong**")
	turnaround("[link](http://pulumi.com)")
	turnaround("plain `code` *emph* **strong** [link](http://pulumi.com)")
	turnaround(`(Block List, Max: 1) The definition for a Change  widget. (see [below for nested schema]` +
		`(#nestedblock--widget--group_definition--widget--change_definition))`)

	res, err := parseTextSequenceFromNode(parseNode("request_max_bytes").FirstChild, false)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "request_max_bytes", res)
}

func TestGetDocsForResource(t *testing.T) {

	doc, err := GetDocsForResource("SAP", "btp", "btp", "resources", "btp_subaccount_environment_instance", "v1.4.0", "github.com")

	if err != nil {
		t.Errorf("error is not expected")
	}

	assert.Equal(t, len(doc.Attributes), 6)
	assert.Equal(t, doc.Import, "import {\n\t\t\t\tto = btp_subaccount_environment_instance.<resource_name>\n\t\t\t\tid = \"<subaccount_id>,<environment_instance_id>\"\n\t\t\t  }")
}

func TestGetDocsForResource_MissingDoc(t *testing.T) {
	org := "testOrg"
	provider := "testProvider"
	resourcePrefix := "testPrefix"
	kind := DocKind("testKind") // Replace with actual DocKind value
	rawname := "testRawname"
	providerModuleVersion := "v1.0.0"
	githost := "github.com"

	_, err := GetDocsForResource(org, provider, resourcePrefix, kind, rawname, providerModuleVersion, githost)

	if err == nil {
		t.Errorf("expected an error but got nil")
	}

	expectedErrorMsg := "could not find docs for testKind testRawname"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("expected error message %q but got %q", expectedErrorMsg, err.Error())
	}
}
