package services

type ISmsService interface {
	SendOTP(phonenumber string, code string)
}
