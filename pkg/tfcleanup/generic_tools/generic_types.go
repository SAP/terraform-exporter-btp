package generictools

const EmptyString = "null"
const EmptyJson = "jsonencode({})"

type VariableInfo struct {
	Description string
	Value       string
}

type VariableContent map[string]VariableInfo
