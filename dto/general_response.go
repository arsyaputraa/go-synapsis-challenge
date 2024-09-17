package dto

// GeneralResponse represents a standard API response
type GeneralResponse struct {
    Status  string      `json:"status"`             // "success" or "error"
    Message string      `json:"message,omitempty"`  // A message describing the response
    Data    interface{} `json:"data,omitempty"`     // The actual data, if any
    Error   interface{} `json:"error,omitempty"`    // Error details, if any
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}, message string) GeneralResponse {
    return GeneralResponse{
        Status:  "success",
        Message: message,
        Data:    data,
    }
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, err interface{}) GeneralResponse {
    return GeneralResponse{
        Status:  "error",
        Message: message,
        Error:   err,
    }
}