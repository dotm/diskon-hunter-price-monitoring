package constenum

type AlertMethod string

const (
	AlertMethodEmail            AlertMethod = "Email"
	AlertMethodPushNotification AlertMethod = "PushNotification"
	AlertMethodSMS              AlertMethod = "SMS"
	AlertMethodWhatsApp         AlertMethod = "WhatsApp"

	AlertMethodUnknown AlertMethod = "Unknown"
)

// kebab-case
func (x AlertMethod) ToKebabCase() string {
	switch x {
	case AlertMethodEmail:
		return "email"
	case AlertMethodPushNotification:
		return "push-notification"
	case AlertMethodSMS:
		return "sms"
	case AlertMethodWhatsApp:
		return "whatsapp"
	default:
		return "unknown"
	}
}

// camelCase
func (x AlertMethod) ToCamelCase() string {
	switch x {
	case AlertMethodEmail:
		return "email"
	case AlertMethodPushNotification:
		return "pushNotification"
	case AlertMethodSMS:
		return "sms"
	case AlertMethodWhatsApp:
		return "whatsapp"
	default:
		return "unknown"
	}
}

// PascalCase
func (x AlertMethod) ToPascalCase() string {
	switch x {
	case AlertMethodEmail:
		return "Email"
	case AlertMethodPushNotification:
		return "PushNotification"
	case AlertMethodSMS:
		return "SMS"
	case AlertMethodWhatsApp:
		return "WhatsApp"
	default:
		return "Unknown"
	}
}

// Used to parse value from frontend client
func NewAlertMethod(dtoValue string) AlertMethod {
	switch dtoValue {
	case "Email":
		return AlertMethodEmail
	case "PushNotification":
		return AlertMethodPushNotification
	case "SMS":
		return AlertMethodSMS
	case "WhatsApp":
		return AlertMethodWhatsApp
	default:
		return AlertMethodUnknown
	}
}
