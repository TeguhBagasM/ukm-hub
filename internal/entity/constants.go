package entity

const (
	OrganizationStatusActive   = "active"
	OrganizationStatusInactive = "inactive"

	DivisionStatusActive   = "active"
	DivisionStatusInactive = "inactive"

	EventStatusDraft     = "DRAFT"
	EventStatusPublished = "PUBLISHED"
	EventStatusClosed    = "CLOSED"
	EventStatusArchived  = "ARCHIVED"

	RegistrationStatusPending  = "PENDING"
	RegistrationStatusAccepted = "ACCEPTED"
	RegistrationStatusRejected = "REJECTED"

	MemberStatusActive   = "active"
	MemberStatusInactive = "inactive"

	FormFieldTypeText     = "TEXT"
	FormFieldTypeTextarea = "TEXTAREA"
	FormFieldTypeEmail    = "EMAIL"
	FormFieldTypeNumber   = "NUMBER"
	FormFieldTypePhone    = "PHONE"
	FormFieldTypeDate     = "DATE"
	FormFieldTypeSelect   = "SELECT"
	FormFieldTypeRadio    = "RADIO"
	FormFieldTypeCheckbox = "CHECKBOX"
	FormFieldTypeURL      = "URL"
)

var SupportedFormFieldTypes = []string{
	FormFieldTypeText,
	FormFieldTypeTextarea,
	FormFieldTypeEmail,
	FormFieldTypeNumber,
	FormFieldTypePhone,
	FormFieldTypeDate,
	FormFieldTypeSelect,
	FormFieldTypeRadio,
	FormFieldTypeCheckbox,
	FormFieldTypeURL,
}

func IsSupportedFieldType(t string) bool {
	for _, s := range SupportedFormFieldTypes {
		if s == t {
			return true
		}
	}
	return false
}
