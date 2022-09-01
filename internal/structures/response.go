package structures

type WorkerResponse struct {
	Success  bool   `json:"success"`
	Location string `json:"location"`
}
