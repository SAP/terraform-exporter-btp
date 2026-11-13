package tfutils

import (
	"net/url"
	"strings"
)

// Assignment inventory values are scoped to the selected subaccount. Preserve
// the provider identity's case and escape each component so separators in names
// cannot merge distinct assignments. A blank result makes query parsing reject
// an incomplete or ambiguous identity instead of exporting the wrong assignment.
func roleCollectionAssignmentInventoryValue(identity map[string]any) string {
	fields := make(map[string]string)

	for _, key := range []string{
		"subaccount_id", "role_collection_name", "origin",
		"user_name", "group_name", "attribute_name", "attribute_value",
	} {

		if value, exists := identity[key]; exists && value != nil {
			text, ok := value.(string)
			if !ok || (text == "" && key != "origin") {
				return ""
			}
			fields[key] = text
		}

	}

	if fields["subaccount_id"] == "" || fields["role_collection_name"] == "" {
		return ""
	}

	if _, hasOrigin := fields["origin"]; !hasOrigin {
		return ""
	}

	parts := []string{fields["role_collection_name"], fields["origin"]}
	user, group := fields["user_name"], fields["group_name"]
	attribute, value := fields["attribute_name"], fields["attribute_value"]

	switch {
	case user != "" && group == "" && attribute == "" && value == "":
		parts = append(parts, "user", user)
	case group != "" && user == "" && attribute == "" && value == "":
		parts = append(parts, "group", group)
	case attribute != "" && value != "" && user == "" && group == "":
		parts = append(parts, "attribute", attribute, value)
	default:
		return ""
	}

	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}

	return strings.Join(parts, "|")
}
