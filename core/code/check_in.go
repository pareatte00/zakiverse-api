package code

const (
	// @Http StatusBadRequest
	// @Locale EN "Cannot claim yet"
	CheckInTooEarly Code = "CHK-001"

	// @Http StatusBadRequest
	// @Locale EN "Maximum claims reached"
	CheckInMaxReached Code = "CHK-002"

	// @Http StatusBadRequest
	// @Locale EN "Check-in plan is not active"
	CheckInPlanInactive Code = "CHK-003"
)
