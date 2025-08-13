package dto

type CreateSpaceDTO struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Status   string `json:"status" binding:"required"`
	Notice   string `json:"notice"`
	Capacity int    `json:"capacity" binding:"required"`
}

type UpdateSpaceDTO struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Notice   string `json:"notice"`
	Capacity int    `json:"capacity"`
}

type SpaceOutputDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Notice   string `json:"notice"`
	Capacity int    `json:"capacity"`
}

type UpdateStatusDTO struct {
	Status string `json:"status" binding:"required"`
}

type UpdateNoticeDTO struct {
	Notice string `json:"notice" binding:"required"`
}
