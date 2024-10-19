package types

type CreateWorkRequest struct {
	Description string `json:"description"`
}

type CreateWorkResponse struct {
	WorkID      int    `json:"workId"`
	Description string `json:"description"`
	TaskIds     []int  `json:"taskIds"`
}
