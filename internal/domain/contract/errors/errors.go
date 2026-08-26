package errors

import "errors"

var (
	ErrContractNotFound = errors.New("contract not found")
	ErrActiveContractNotFound = errors.New("no active contract found")
	ErrCompanyAlreadyHasActiveContract = errors.New("company already has an active contract")
    ErrInvalidStatusTransition = errors.New("invalid status transition")
    ErrCannotUpdateTerminated = errors.New("cannot update terminated contract")
    ErrContractExpired = errors.New("contract has expired")
    ErrServiceNotEnabled = errors.New("service not enabled")
    ErrServicesListEmpty = errors.New("services list cannot be empty")
    ErrDuplicateService = errors.New("duplicate service")
    ErrInvalidServiceType = errors.New("invalid service type")
    ErrInvalidStatus = errors.New("invalid status")
    ErrUnknownStatus = errors.New("unknown status")
    ErrInvalidContractID = errors.New("invalid contract id")
    ErrInvalidCompanyID = errors.New("invalid company id")
    ErrInvalidDateFormat = errors.New("invalid date format")
    ErrInvalidDateRange = errors.New("invalid date range")
    ErrMissingRequiredField = errors.New("missing required field")
    ErrServiceFailed = errors.New("service operation failed")
)

func GetErrorMessage(err error) string {
	switch {
	case errors.Is(err, ErrContractNotFound):
		return "Contract not found"
	case errors.Is(err, ErrActiveContractNotFound):
		return "Active contract not found"
	case errors.Is(err, ErrCompanyAlreadyHasActiveContract):
		return "Company already has an active contract"
	case errors.Is(err, ErrInvalidStatusTransition):
		return "Invalid status transition"
	case errors.Is(err, ErrCannotUpdateTerminated):
		return "Cannot update terminated contract"
	case errors.Is(err, ErrContractExpired):
		return "Contract has expired"
	case errors.Is(err, ErrServiceNotEnabled):
		return "Service is not enabled"
	case errors.Is(err, ErrServicesListEmpty):
		return "Services list cannot be empty"
	case errors.Is(err, ErrDuplicateService):
		return "Duplicate service found"
	case errors.Is(err, ErrInvalidServiceType):
		return "Invalid service type"
	case errors.Is(err, ErrInvalidStatus):
		return "Invalid status"
	case errors.Is(err, ErrUnknownStatus):
		return "Unknown status"
	case errors.Is(err, ErrInvalidContractID):
		return "Invalid contract ID"
	case errors.Is(err, ErrInvalidCompanyID):
		return "Invalid company ID"
	case errors.Is(err, ErrInvalidDateFormat):
		return "Invalid date format"
	case errors.Is(err, ErrInvalidDateRange):
		return "Invalid date range"
	case errors.Is(err, ErrMissingRequiredField):
		return "Missing required field"
	default:
		return "Internal server error"
	}
}

func GetHTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrContractNotFound),
		errors.Is(err, ErrActiveContractNotFound):
		return 404
	case errors.Is(err, ErrCompanyAlreadyHasActiveContract),
		errors.Is(err, ErrInvalidStatusTransition),
		errors.Is(err, ErrCannotUpdateTerminated),
		errors.Is(err, ErrContractExpired),
		errors.Is(err, ErrServiceNotEnabled):
		return 409
	case errors.Is(err, ErrServicesListEmpty),
		errors.Is(err, ErrDuplicateService),
		errors.Is(err, ErrInvalidServiceType),
		errors.Is(err, ErrInvalidStatus),
		errors.Is(err, ErrUnknownStatus),
		errors.Is(err, ErrInvalidContractID),
		errors.Is(err, ErrInvalidCompanyID),
		errors.Is(err, ErrInvalidDateFormat),
		errors.Is(err, ErrInvalidDateRange),
		errors.Is(err, ErrMissingRequiredField):
		return 400
	default:
		return 500
	}
}