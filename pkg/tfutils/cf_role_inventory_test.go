package tfutils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCfRoleInventoryDistinguishesUserOrigins(t *testing.T) {
	// CF assigns different GUIDs to users with the same username in different origins.
	origins := map[string]string{"user-ldap": "ldap", "user-sap": "sap.ids"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/":
			_, _ = fmt.Fprintf(w, "{\"links\":{\"uaa\":{\"href\":%q}}}", "http://"+r.Host)
		case "/v3/spaces":
			_, _ = w.Write([]byte("{\"pagination\":{\"next\":null},\"resources\":[{\"guid\":\"space-id\",\"name\":\"dev\"}]}"))
		case "/v3/users":
			guid := r.URL.Query().Get("guids")
			origin, ok := origins[guid]
			assert.True(t, ok, "unexpected user GUID %q", guid)
			_, _ = fmt.Fprintf(w, "{\"pagination\":{\"next\":null},\"resources\":[{\"guid\":%q,\"username\":\"user@example.com\",\"origin\":%q}]}", guid, origin)
		default:
			t.Errorf("unexpected CF request: %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	for _, key := range []string{"CF_USER", "CF_PASSWORD", "CF_CLIENT_ID", "CF_CLIENT_SECRET", "CF_REFRESH_TOKEN"} {
		t.Setenv(key, "")
	}
	t.Setenv("CF_API_URL", server.URL)
	// Synthetic JWT with an expiry in 2100; no live credentials are used.
	t.Setenv("CF_ACCESS_TOKEN", "eyJhbGciOiJub25lIn0."+base64.RawURLEncoding.EncodeToString([]byte("{\"exp\":4102444800}"))+".")

	for _, command := range []string{CmdCfSpaceRoleParameter, CmdCfOrgRoleParameter} {
		t.Run(command, func(t *testing.T) {
			spec, err := GetListResourceSpec(command, OrganizationLevel)
			require.NoError(t, err)
			var events []string
			for index, user := range []string{"user-ldap", "user-sap"} {
				roleType := "space_developer"
				if command == CmdCfOrgRoleParameter {
					roleType = "organization_user"
				}
				config := fmt.Sprintf("resource %q %q {\n user = %q\n type = %q\n}\n", spec.ResourceType, fmt.Sprintf("export_%d", index), user, roleType)
				importConfig := fmt.Sprintf("import {\n to = %s.export_%d\n identity = { role_guid = %q }\n}\n", spec.ResourceType, index, "role-"+user)
				event, err := json.Marshal(map[string]any{
					"type": "list_resource_found",
					"list_resource_found": map[string]any{
						"identity":        map[string]any{"role_guid": "role-" + user},
						"resource_object": map[string]any{"type": roleType, "user": user, "space": "space-id", "org": "org-id"},
						"config":          config, "import_config": importConfig,
					},
				})
				require.NoError(t, err)
				events = append(events, string(event))
			}
			results, diagnostics, err := parseQueryOutput([]byte(strings.Join(events, "\n")), spec)
			require.NoError(t, err)
			require.Empty(t, diagnostics)
			require.Len(t, results, 2)
			require.NotEqual(t, results[0].InventoryValue, results[1].InventoryValue)
			if command == CmdCfSpaceRoleParameter {
				assert.Equal(t, "dev_space-developer_user@example.com_space-id_user-ldap", results[0].InventoryValue)
			}
			all, err := filterQueryResults(spec, results, "", "", nil)
			require.NoError(t, err)
			require.Len(t, all, 2, "unfiltered export must retain both assignments")
			for _, result := range results {
				selected, err := filterQueryResults(spec, results, "", "", []string{result.InventoryValue})
				require.NoError(t, err)
				require.Len(t, selected, 1, "JSON selection must select only one origin's assignment")
				assert.Equal(t, result.Identity, selected[0].Identity)
				directory := t.TempDir()
				cwd, err := os.Getwd()
				require.NoError(t, err)
				relative, err := filepath.Rel(cwd, directory)
				require.NoError(t, err)
				require.NoError(t, WriteQueryResults(relative, "roles.tf", spec.ResourceType, selected))
				imports, err := os.ReadFile(filepath.Join(directory, spec.ResourceType+"_import.tf"))
				require.NoError(t, err)
				assert.Equal(t, result.ImportConfig+"\n\n", string(imports))
			}
		})
	}
}
