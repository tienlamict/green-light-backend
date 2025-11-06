package dto

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

func SuccessResponse(data interface{}, message string) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(err string) Response {
	return Response{
		Success: false,
		Error:   err,
	}
}

func PaginatedSuccessResponse(data interface{}, total int64, page, limit int) PaginatedResponse {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}
}
