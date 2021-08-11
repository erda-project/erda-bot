package events

const TypePR = "pull_request"

const (
	PRActionOpened   = "opened"
	PRActionEdited   = "edited"
	PRActionReopened = "reopened"
)

var SupportedPREventActions = []string{PRActionOpened, PRActionEdited, PRActionReopened}
