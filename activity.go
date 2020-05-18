package sendmail

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strconv"
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

	//get input vars
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
	template := ctx.GetInput("p_template").(string)
	clinic := ctx.GetInput("g_hospital").(string)
	meet := ctx.GetInput("n_meet").(string)
	subject := ctx.GetInput("o_subject").(string)
	image_footer := ctx.GetInput("q_image_footer").(string)
	link_footer := ctx.GetInput("r_link_footer").(string)
	image_footer_alt := ctx.GetInput("s_image_footer_alt").(string)
	appointment_id := ctx.GetInput("u_appointment_id").(string)
	status := ctx.GetInput("t_status").(string)
	fdate := strings.Split(date, " ")
	hour := strings.Split(fdate[1], ":")
	enddate := ctx.GetInput("v_appointment_end_date").(string)
	organizer := ctx.GetInput("x_ics_organizer").(string)
	prodid := ctx.GetInput("z_ics_prodid").(string)

	cport, e1 := strconv.Atoi(port)
	fmt.Println(e1)

	method := "CANCEL"
	fstatus := "CANCELLED"
	transp := "TRANSPARENT"
	if status != "cancelled" {
		method = "PUBLISH"
		fstatus = "CONFIRMED"
		transp = "OPAQUE"
	}

	fmt.Println(date)

	date1 := time.Now()
	fdate1 := date1.Format("20060102T150405Z")

	content := "BEGIN:VCALENDAR\r" +
		"METHOD:" + method + "\r" +
		"PRODID:" + prodid + "\r" +
		"VERSION:2.0\r" +
		"BEGIN:VEVENT\r" +
		"DTSTAMP:" + fdate1 + "\r" +
		"UID:" + appointment_id + "\r" +
		"SEQUENCE:0\r" +
		"ORGANIZER;" + organizer + "\r" +
		"DTSTART:" + parseDate(date) + "\r" +
		"DTEND:" + parseDate(enddate) + "\r" +
		"STATUS:" + fstatus + "\r" +
		"CATEGORIES:" + appointment + " " + clinic + "\r" +
		"SUMMARY:" + appointment + " " + clinic + "\r" +
		"CLASS:PUBLIC\r" +
		"TRANSP:" + transp + "\r" +
		"END:VEVENT\r" +
		"END:VCALENDAR\r"


	filename1 := CreateTempFile(content)

	//create email

	var (
		serverAddr         = server
		password           = apppass
		emailAddr          = sender
		portNumber         = cport
		tos                = ercpnt
		attachmentFilePath = filename1
		filename           = "invite.ics"
		delimeter          = "**=cuf689407924327"
	)

	tlsConfig := tls.Config{
		ServerName:         serverAddr,
		InsecureSkipVerify: true,
	}

	conn, connErr := tls.Dial("tcp", fmt.Sprintf("%s:%d", serverAddr, portNumber), &tlsConfig)
	if connErr != nil {
		log.Panic(connErr)
	}
	defer conn.Close()

	client, clientErr := smtp.NewClient(conn, serverAddr)
	if clientErr != nil {
		log.Panic(clientErr)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", emailAddr, password, serverAddr)

	if err := client.Auth(auth); err != nil {
		log.Panic(err)
	}

	if err := client.Mail(emailAddr); err != nil {
		log.Panic(err)
	}

	if err := client.Rcpt(tos); err != nil {
		log.Panic(err)
	}

	writer, writerErr := client.Data()
	if writerErr != nil {
		log.Panic(writerErr)
	}

	sampleMsg := fmt.Sprintf("From: %s\r\n", emailAddr)
	sampleMsg += fmt.Sprintf("To: %s\r\n", tos)
	sampleMsg += "Subject: " + subject + "\r\n"
	sampleMsg += "MIME-Version: 1.0\r\n"
	sampleMsg += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", delimeter)
	sampleMsg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	sampleMsg += "Content-Type: text/html; charset=\"utf-8\"\r\n"
	sampleMsg += "Content-Transfer-Encoding: 7bit\r\n"

	templateData := struct {
		Name         string
		Appointment  string
		Speciality   string
		Practitioner string
		Date         string
		Hour         string
		Meet         string
		Hospital     string
		Footer       string
		Image        string
		Alt          string
	}{
		Name:         patient,
		Appointment:  appointment,
		Speciality:   speciality,
		Practitioner: practitioner,
		Date:         fdate[0],
		Hour:         hour[0] + ":" + hour[1],
		Meet:         meet,
		Hospital:     clinic,
		Footer:       link_footer,
		Image:        image_footer,
		Alt:          image_footer_alt,
	}

	r := NewRequest([]string{ercpnt}, subject, "")
	error1 := r.ParseTemplate(template+".html", templateData)
	if error1 := r.ParseTemplate(template+".html", templateData); error1 == nil {
		sampleMsg += r.body

		sampleMsg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
		sampleMsg += "Content-Type: text/calendar; charset=\"utf-8\"\r\n"
		sampleMsg += "Content-Transfer-Encoding: base64\r\n"
		sampleMsg += "Content-Disposition: attachment;filename=\"" + filename + "\"\r\n"

		rawFile, fileErr := ioutil.ReadFile(attachmentFilePath)
		if fileErr != nil {
			log.Panic(fileErr)
		}
		sampleMsg += "\r\n" + base64.StdEncoding.EncodeToString(rawFile)

		//write into email client stream writter
		log.Println("Write content into client writter I/O")
		if _, err := writer.Write([]byte(sampleMsg)); err != nil {
			log.Panic(err)
		}

		if closeErr := writer.Close(); closeErr != nil {
			log.Panic(closeErr)
		}

		client.Quit()

		log.Print("done.")

		defer os.Remove(filename)

	}
	fmt.Println(error1)

	return true, nil
}

func CreateTempFile(serializer string) string {

	tmpFile, err := ioutil.TempFile(os.TempDir(), "*.ics")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	// Remember to clean up the file afterwards
	//defer os.Remove(tmpFile.Name())

	fmt.Println("Created File: " + tmpFile.Name())

	// Example writing to the file
	text := []byte(serializer)
	if _, err = tmpFile.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpFile.Name()
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

func parseDate(date string) (teste string) {
	teste = ""

	datetime := strings.Split(date, " ")
	fdate := strings.Split(datetime[0], "/")

	for i := len(fdate) - 1; i >= 0; i-- {
		teste += fdate[i]
	}

	teste += "T"

	fhour := strings.Split(datetime[1], ":")

	for i := 0; i < len(fhour); i++ {
		teste += fhour[i]
	}
	teste += "Z"

	return teste

}
