package dto

import "github.com/google/uuid"

type ResponseDto struct {
	RequestId   uuid.UUID `json:"requestId,omitempty"`
	Description string    `json:"description"`
	Data        any       `json:"data"`
}
