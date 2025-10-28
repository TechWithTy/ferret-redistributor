package instagram

type creationResponse struct {
    ID string `json:"id"`
}

type publishResponse struct {
    ID string `json:"id"`
}

type containerStatusResponse struct {
    StatusCode string `json:"status_code"`
}

type graphErrorPayload struct {
    Message       string `json:"message"`
    Type          string `json:"type"`
    Code          int    `json:"code"`
    ErrorSubcode  int    `json:"error_subcode"`
    FBTraceID     string `json:"fbtrace_id"`
}

type graphErrorResponse struct {
    Error graphErrorPayload `json:"error"`
}

