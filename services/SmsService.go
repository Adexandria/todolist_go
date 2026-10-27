package services

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type SmsService struct {
	Logger *slog.Logger
}

func SmsCon(jsonHandler *slog.JSONHandler) *SmsService {
	return &SmsService{
		Logger: slog.New(jsonHandler),
	}
}
func (s *SmsService) SendOTP(phonenumber string, code string) {
	from := os.Getenv("TWILIO_NUMBER")
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	apiKey := os.Getenv("TWILIO_API_KEY")
	apiSecret := os.Getenv("TWILIO_API_SECRET")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   apiKey,
		Password:   apiSecret,
		AccountSid: accountSid,
	})

	msg := fmt.Sprintf("Your verification code is: %s. Don't share this code with anyone; our employees will never ask for the code. This code is only valid for 3 minutes", code)

	params := &api.CreateMessageParams{}
	params.SetBody(msg)
	params.SetFrom(from)
	params.SetTo(phonenumber)

	_, err := client.Api.CreateMessage(params)

	if err != nil {
		s.Logger.Info(err.Error())
	} else {
		s.Logger.Info("Sent OTP Successfully to User")
	}
}
