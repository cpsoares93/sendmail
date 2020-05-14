package sendmail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/arran4/golang-ical"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
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

	server := ctx.GetInput("a_server").(string)
	port := ctx.GetInput("b_port").(string)
	sender := ctx.GetInput("c_sender").(string)
	apppass := ctx.GetInput("d_password").(string)


	ercpnt := ctx.GetInput("l_patient_contact").(string)
	appointment := ctx.GetInput("e_appointment").(string)
	speciality := ctx.GetInput("f_speciality").(string)
	patient := ctx.GetInput("j_patient").(string)
	practitioner := ctx.GetInput("m_practitioner").(string)
	date := ctx.GetInput("i_date").(string)
	local := ctx.GetInput("h_local").(string);
	template:= ctx.GetInput("p_template").(string)
	clinic := ctx.GetInput("g_hospital").(string)
	meet := ctx.GetInput("n_meet").(string)
	subject := ctx.GetInput("o_subject").(string)
	image_footer := ctx.GetInput("q_image_footer").(string)
	link_footer := ctx.GetInput("r_link_footer").(string)
	image_footer_alt := ctx.GetInput("s_image_footer_alt").(string)
	fdate := strings.Split(date, " ")

	hour := strings.Split(fdate[1], ":");




	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	cal.SetProductId(" Integrations")
	cal.SetVersion("2.0")
	event := cal.AddEvent("teste@google.com")
	event.SetDtStampTime(time.Now())
	event.SetOrganizer("sender@domain", ics.WithCN("Saúde"))
	event.SetStartAt(time.Now())
	event.SetEndAt(time.Now())
	event.SetStatus(ics.ObjectStatusConfirmed)
	event.SetDescription("teste")
	event.SetSummary("teste1")


	tmpFile, err := ioutil.TempFile(os.TempDir(), "*.ics")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	// Remember to clean up the file afterwards
	//defer os.Remove(tmpFile.Name())

	fmt.Println("Created File: " + tmpFile.Name())

	// Example writing to the file
	text := []byte(cal.Serialize())
	if _, err = tmpFile.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	auth := smtp.PlainAuth("", sender, apppass, server)
	templateData := struct {
		Name string
		Appointment  string
		Speciality string
		Practitioner string
		Date string
		Hour string
		Local string
		Meet string
		Hospital string
		Footer string
		Image string
		Alt string
	}{
		Name: patient,
		Appointment:  appointment,
		Speciality: speciality,
		Practitioner: practitioner,
		Date: fdate[0],
		Hour: hour[0] + ":" + hour[1],
		Local: local,
		Meet: meet,
		Hospital: clinic,
		Footer: link_footer,
		Image: image_footer,
		Alt: image_footer_alt,
	}
	r := NewRequest([]string{ercpnt}, subject , "")
	error1 := r.ParseTemplate(template + ".html", templateData)
	if error1 := r.ParseTemplate(template + ".html", templateData); error1 == nil {
		ok, _ := r.SendEmail(auth, port, sender, tmpFile.Name())
		fmt.Println(ok)
	}
	fmt.Println(error1)
	ctx.SetOutput("output", "Mail_Sent_Successfully")
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

func (r *Request) SendEmail(auth smtp.Auth, port string, sender string, filename string) (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: multipart/mixed; charset=\"UTF-8\";Content-Transfer-Encoding: 7bit\n\n";
	subject := "Subject: " + r.subject + "\n"

	attachment := "Content-Type: text/calendar; charset=\"utf-8\"\r\n"
	attachment += "Content-Transfer-Encoding: base64\r\n"
	attachment += "Content-Disposition: attachment;filename=\"invite.ics\"\r\n"
	//read file
	rawFile, fileErr := ioutil.ReadFile(filename)
	if fileErr != nil {
		log.Panic(fileErr)
	}
	r.body += "\r\n" + base64.StdEncoding.EncodeToString(rawFile)
	msg := []byte(subject + mime + "\n" + r.body + attachment)


	addr := "smtp.gmail.com:"+port

	if err := smtp.SendMail(addr, auth, sender, r.to, msg); err != nil {
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