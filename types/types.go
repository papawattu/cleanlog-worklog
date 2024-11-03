package types

type WorkResponse struct {
	WorkID      int    `json:"workId"`
	Description string `json:"description"`
	TaskIds     []int  `json:"taskIds"`
	Date        string `json:"date"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
type CreateWorkRequest struct {
	Description string `json:"description"`
	Date        string `json:"date"`
}

type CreateWorkResponse struct {
	WorkResponse
}
type ListWorkResponse struct {
	WorkResponses []WorkResponse `json:"worklogs"`
}
