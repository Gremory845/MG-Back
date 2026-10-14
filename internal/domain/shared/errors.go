package shared

// DomainError representa un error propio del dominio (no un error técnico
// de infraestructura). presentation/http/response lo traduce a un status
// HTTP usando el Code, en un solo lugar, para que ningún handler tenga
// que repetir ese switch.
type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewDomainError(code, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}
