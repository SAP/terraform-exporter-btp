package btpcli

const (
	ActionGet  Action = "get"
	ActionList Action = "list"
)

func NewGetRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionGet, command, args)
}

func NewListRequest(command string, args any) *CommandRequest {
	return NewCommandRequest(ActionList, command, args)
}
