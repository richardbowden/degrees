package notification

type TemplateType string

const (
	TPL_WELCOME_EMAIL               TemplateType = "system-welcome-email"
	TPL_SYSTEM_VERIFY_EMAIL_ADDRESS TemplateType = "system-verify-email-address"
	TPL_SYSTEM_PASSWORD_RESET       TemplateType = "system-password-reset"
	TPL_BOOKING_CONFIRMATION        TemplateType = "booking-confirmation"
)

func (s TemplateType) String() string {
	return string(s)

}

type Template struct {
	ID      int          `json:"id"`
	Type    TemplateType `json:"type"`
	Name    string       `json:"name"`
	Version string       `json:"version"`
	Content string       `json:"content"`
	Scope   string       `json:"scope"`
}
