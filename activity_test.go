package sendmail

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"io/ioutil"
	"testing"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	//act := NewActivity(getActivityMetadata())
	//tc := test.NewTestActivityContext(getActivityMetadata())
	//
	////setup attrs
	//tc.SetInput("A_server", "smtp.gmail.com")
	//tc.SetInput("B_port", "587")
	//tc.SetInput("C_sender", "carolinasoares.cps@gmail.com")
	//tc.SetInput("D_apppassword", "hakunamatata93")
	//tc.SetInput("E_rcpnt", "carolina.soares@litthub.com")
	//tc.SetInput("F_sub", "Q_Subscriber_Down!")
	//tc.SetInput("G_body", "Subscriber_for_Queue_is_down.")
	//
	//
	//done, err := act.Eval(tc)
	//if !done {
	//	fmt.Println(err)
	//}
	//act.Eval(tc)
	////check output attr
	//
	//output := tc.GetOutput("output")
	//assert.Equal(t, output, output)
	//output := tc.GetOutput("SentTime")
	//assert.Equal(t, SentTime, SentTime)


}
