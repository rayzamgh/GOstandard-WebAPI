package responses

/*
============================
			NOTE
============================
This is responses template
*/

type SingleResponse struct {
	Data interface{} `json:"data"`
	Meta SimpleMeta  `json:"meta"`
}

type VueTablePaginateResponse struct {
	CurrentPage int64       `json:"current_page"`
	From        int64       `json:"from"`
	LastPage    int64       `json:"last_page"`
	NextPageUrl string      `json:"next_page_url,omitempty"`
	Path        string      `json:"path"`
	PerPage     int64       `json:"per_page"`
	PrevPageUrl string      `json:"prev_page_url,omitempty"`
	To          int64       `json:"to"`
	Total       int64       `json:"total"`
	Data        interface{} `json:"data"`
	Meta        SimpleMeta  `json:"meta"`
}

type PaginateResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type SimpleMeta struct {
	StatusCode int      `json:"status_code"`
	Message    interface{} `json:"message"`
}

type PaginationMeta struct {
	Meta       SimpleMeta `json:"meta"`
	Pagination Paginator  `json:"pagination"`
}

type Paginator struct {
	Total       int               `json:"total"`
	Count       int               `json:"count"`
	PerPage     int64             `json:"per_page"`
	CurrentPage int64             `json:"current_page"`
	TotalPages  int               `json:"total_pages"`
	Links       map[string]string `json:"links"`
}

type Paginate struct {
	Data        interface{} `json:"data"`
	CurrentPage int64       `json:"current_page"`
	From        int64       `json:"from"`
	LastPage    int64       `json:"last_page"`
	NextPageUrl string      `json:"next_page_url,omitempty"`
	Path        string      `json:"path"`
	Page        int64       `json:"page"`
	PerPage     int64       `json:"per_page"`
	PrevPageUrl string      `json:"prev_page_url,omitempty"`
	To          int64       `json:"to"`
	Total       int64       `json:"total"`
}

func NewResponse(data interface{}) *SingleResponse {
	response := &SingleResponse{
		Data: data,
		Meta: SimpleMeta{
			StatusCode: 200,
			Message:    "Success",
		},
	}
	return response
}
