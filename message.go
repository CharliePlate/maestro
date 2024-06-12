package maestro

type ActionType string

const (
	ActionTypeAcknowledge ActionType = "acknowledge"
	ActionTypeSubscribe   ActionType = "subscribe"
)

type Message struct {
	Content    interface{}
	ConnID     string
	ActionType ActionType
}
