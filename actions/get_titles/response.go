package get_titles

type Response struct {
	Info []Info `json:"info"`
}

type Info struct {
	Url   string `json:"url"`
	Err   string `json:"error,omitempty"`
	Title string `json:"page_info,omitempty"`
}
