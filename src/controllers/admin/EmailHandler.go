package admin

import (
	"github.com/Demistry/Hotel-Management-System/src/utils"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

func sendMail(emailAddress string, username string, userId string){
	from := mail.NewEmail("HotSys", "Hotsys@mail.com")
	subject := "Email Verification for HotSys"
	to := mail.NewEmail(username, emailAddress)
	content := mail.NewContent("text/plain", "Click on the link below to verify your email address for " + username + "\n " + utils.ConfirmAdminMailEndpoint+ userId + "\nThis link expires in 7 days.")
	m := mail.NewV3MailInit(from, subject, to, content)
	apiKey,ok := os.LookupEnv("SENDGRID_API_KEY")
	if ok == false{
		apiKey = os.Getenv("SENDGRID_API_KEY")
	}
	request := sendgrid.GetRequest(apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	}
}

func sendResetPasswordMail(emailAddress string, username string, userId string){
	from := mail.NewEmail("HotSys", "Hotsys@mail.com")
	subject := "Password Reset for HotSys"
	to := mail.NewEmail(username, emailAddress)
	content := mail.NewContent("text/plain", "Click on the link below to reset the password for your HotSys account\n " + utils.ResetPasswordAdminMailEndpoint+ userId + "\nThis link expires in 15 minutes. Ignore this mail if you had nothing to do with this.")
	m := mail.NewV3MailInit(from, subject, to, content)
	apiKey,ok := os.LookupEnv("SENDGRID_API_KEY")
	if ok == false{
		apiKey = os.Getenv("SENDGRID_API_KEY")
	}
	request := sendgrid.GetRequest(apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	}
}