package util

import (
	"github.com/twilio/twilio-go"
	lookups "github.com/twilio/twilio-go/rest/lookups/v2"
)

type PhoneDetails struct {
	PhoneNumber      string
	Valid            bool
	ValidationErrors []string
}

func VerifyPhone(phoneNumber, accountSid, authToken string) (PhoneDetails, error) {
	// client := twilio.NewRestClient()
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	params := &lookups.FetchPhoneNumberParams{}

	resp, err := client.LookupsV2.FetchPhoneNumber(phoneNumber, params)
	if err != nil {
		return PhoneDetails{}, err
	}

	phoneDetails := PhoneDetails{
		PhoneNumber:      *resp.PhoneNumber,
		Valid:            *resp.Valid,
		ValidationErrors: *resp.ValidationErrors,
	}

	return phoneDetails, nil
}
