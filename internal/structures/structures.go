package structures

type WorkerResponse struct {
	Success  bool   `json:"success"`
	Location string `json:"location"`
}

type BackendPublishPayload struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}
