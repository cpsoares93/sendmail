package main

import (

	"fmt"
	"log"
	"strings"
	"time"
	"net/smtp"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

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
	
	
	server := ctx.GetInput("A_server").(string)
	port := ctx.GetInput("B_port").(string)
	sender := ctx.GetInput("C_sender").(string)
	apppass := ctx.GetInput("D_apppassword").(string)
	ercpnt := ctx.GetInput("E_rcpnt").(string)
	fsub := ctx.GetInput("F_sub").(string)
	gbody := ctx.GetInput("G_body").(string)
	
	dt := time.Now()
	
	// Set up authentication information.
	//auth := smtp.PlainAuth(
	//	"",
	//	"sendalertsforq@gmail.com",
	//	"ptcxejoylzgtrfmh",
	//	"smtp.gmail.com",
	//)
	
	auth := smtp.PlainAuth("", sender, apppass, server,)
	
	t := []string{"To:", ercpnt}
	strings.Join(t, " ")
	
	s := []string{"Subject:", fsub}
	strings.Join(s, " ")
	
	serv := []string{server, port}
	strings.Join(serv, ":")
	
	
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	
	to := []string{ercpnt}
	msg := []byte(strings.Join(t, " ") + "\r\n" + strings.Join(s, " ") + "\r\n" + gbody + "\r\n")
	
	err = smtp.SendMail(strings.Join(serv, ":"), auth, sender, to, msg)
	if err != nil {
		activityLog.Debugf("Error ", err)
		fmt.Println("error: ", err)
		return
	}

	fmt.Println("Mail Sent")
	log.Println("Mail Sent")


	// Set the output as part of the context
	activityLog.Debugf("Activity has sent the mail Successfully")
	fmt.Println("Activity has sent the mail Successfully")

	ctx.SetOutput("output", "Mail_Sent_Successfully")
	ctx.SetOutput("SentTime", dt)
	return true, nil
}
