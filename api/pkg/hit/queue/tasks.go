package queue

const TypeHitExecute = "worker:hit:execute"

type ExecuteHitPayload struct {
	TaskID string `json:"task_id"`
	URL    string `json:"url"`
	Method string `json:"method"`
}
