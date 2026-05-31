package dto

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int32   `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"` // 07-2025
	EndDate     *string `json:"end_date"`
}
