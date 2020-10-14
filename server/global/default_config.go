package global

// EnvPrefix for env
const EnvPrefix = "ex"

// DefaultConfig TODO
var DefaultConfig = map[string]interface{}{
	"worker.number":   4,
	"queue.number":    4000,
	"http.addr":       ":7890",
	"http.log.enable": true,
	"http.api.path":   "",
	"limit.http.size": 1 << 20,
}
