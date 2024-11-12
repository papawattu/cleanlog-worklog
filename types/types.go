package types

type WorkResponse struct {
	WorkID      int    `json:"workId"`
	Description string `json:"description"`
	TaskIds     []int  `json:"taskIds"`
	Date        string `json:"date"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	UserID      int    `json:"userId"`
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

type UpdateWorkRequest struct {
	Description string `json:"description"`
	Date        string `json:"date"`
	TaskIds     []int  `json:"taskIds"`
}

type AddTaskRequest struct {
	TaskId int `json:"taskId"`
}
