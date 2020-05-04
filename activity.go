package sendmail

import (
	"bytes"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"html/template"
	"net/smtp"
)

type Appointment struct {
	patient string
}

// ActivityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-flogo-sendmail")

// MyActivity is a stub for your Activity implementation
type sendmail struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &sendmail{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *sendmail) Metadata() *activity.Metadata {
	return a.metadata
}


// Eval implements activity.Activity.Eval
func (a *sendmail) Eval(ctx activity.Context) (done bool, err error) {

	server := ctx.GetInput("server").(string)
	port := ctx.GetInput("port").(string)
	sender := ctx.GetInput("sender").(string)
	apppass := ctx.GetInput("apppassword").(string)
	ercpnt := ctx.GetInput("rcpnt").(string)
	appointment := ctx.GetInput("appointment").(*string)
	speciality := ctx.GetInput("speciality").(*string)
	patient := ctx.GetInput("patient").(*string)
	practitioner := ctx.GetInput("practitioner").(*string)
	//date := ctx.GetInput("date").(*)
	clinic := "?"

	clientAppointment:= *appointment
	clientsSpeciality := *speciality
	clientPatient := *patient
	clientPractitioner := *practitioner
	//clientDate := *date

	auth := smtp.PlainAuth("", sender, apppass, server)
	templateData := struct {
		Name string
		Appointment  string
		Speciality string
		Practitioner string
		Date string
		Hour string
		Local string
	}{
		Name: clientPatient,
		Appointment:  clientAppointment,
		Speciality: clientsSpeciality,
		Practitioner: clientPractitioner,
		Date: "?",
		Hour: "?",
		Local: "?",
	}
	r := NewRequest([]string{ercpnt}, clientAppointment + " - " + clinic , "")
	error1 := r.ParseTemplate("template.html", templateData)
	if error1 := r.ParseTemplate("template.html", templateData); error1 == nil {
		ok, _ := r.SendEmail(auth, port)
		fmt.Println(ok)
	}
	fmt.Println(error1)


	
	//dt := time.Now()
	
	//auth := smtp.PlainAuth("", sender, apppass, server,)
	//auth = smtp.PlainAuth("", sender, apppass, server)
	//templateData := struct {
	//	Name string
	//	URL  string
	//}{
	//	Name: "Dhanush",
	//	URL:  "http://geektrust.in",
	//}
	//r := NewRequest([]string{"junk@junk.com"}, "Hello Junk!", "Hello, World!")
	//error := r.ParseTemplate("template.html", templateData)
	//if err := r.ParseTemplate("template.html", templateData); err == nil {
	//	ok, _ := r.SendEmail()
	//	fmt.Println(ok)
	//}

	

	ctx.SetOutput("output", "Mail_Sent_Successfully")
	////ctx.SetOutput("SentTime", dt)
	return true, nil
}

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail(auth smtp.Auth, port string) (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n";
	subject := "Subject: " + r.subject + "\n"
	msg := []byte(subject + mime + "\n" + r.body)

	addr := "smtp.gmail.com:"+port

	if err := smtp.SendMail(addr, auth, "carolinasoares.cps@gmail.com", r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}


func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}