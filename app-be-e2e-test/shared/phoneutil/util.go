package phoneutil

import "fmt"

func StandardizePhoneNumberPrefix(rawPhoneNumber string) string {
	phoneNumber := rawPhoneNumber
	if rawPhoneNumber[0] == '0' {
		phoneNumber = fmt.Sprintf("+62%s", rawPhoneNumber[1:]) //default calling code is Indonesia
		return phoneNumber
	}
	if rawPhoneNumber[0] != '+' {
		phoneNumber = fmt.Sprintf("+%s", rawPhoneNumber)
		return phoneNumber
	}
	return phoneNumber
}
