package response

type AvailabilityResponse struct {
	Available bool   `json:"available"`
	Reason    string `json:"reason,omitempty"`
}

func NewAvailabilityResponse(available bool, reason string) *AvailabilityResponse {
	resp := &AvailabilityResponse{
		Available: available,
	}
	if reason != "" {
		resp.Reason = reason
	}
	return resp
}