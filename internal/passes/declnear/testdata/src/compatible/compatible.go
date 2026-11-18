package compatible

func Run() { _ = parseOptions() }

func report() {}

type options struct{ value int }

func parseOptions() options { return options{} }

func consume(settings options) { _ = settings }

// A conventional constructor and its type may form a later block without
// forcing exported API functions below that block.
type configured struct{}

func NewConfigured() *configured { return &configured{} }

func configuredHelper(value configured) { _ = value }
