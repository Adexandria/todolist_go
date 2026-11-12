package services

type SmsService interface {
	SendOTP(phonenumber string, code string)
}
