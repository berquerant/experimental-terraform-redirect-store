package api

type (
	ScanRequest  struct{}
	ScanResponse struct {
		Records []*Record `json:"records,omitempty"`
		Error   string    `json:"error,omitempty"`
	}

	GetRequest struct {
		Name string `json:"name"`
	}
	GetResponse struct {
		Record *Record `json:"record,omitempty"`
		Error  string  `json:"error,omitempty"`
	}

	PutRequest struct {
		Record *Record `json:"record"`
	}
	PutResponse struct {
		Record *Record `json:"record,omitempty"`
		Error  string  `json:"error,omitempty"`
	}

	DeleteRequest struct {
		Name string `json:"name"`
	}
	DeleteResponse struct {
		Error string `json:"error,omitempty"`
	}

	RedirectRequest struct {
		Name string `json:"name"`
	}
	RedirectResponse struct {
		To    string `json:"to"`
		Error string `json:"error,omitempty"`
	}
)
