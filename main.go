package lora_events

import (
	"github.com/guelfey/go.dbus"
        "time"
        "errors"
)


type ConnDetails struct {
        Status int16
        Obj *dbus.Object
        Err error
}

func NewLoraConnection() ConnDetails {
        var lc ConnDetails
        conn, err := dbus.SystemBus()
        if err != nil {
                lc.Err=err
                return lc
        }

        obj := conn.Object("org.cacophony.Lora", "/org/cacophony/Lora")
        lc.Obj = obj
        return lc
}

func (lc *ConnDetails) GetStatus() (int16, error) {
        call := lc.obj.Call("org.cacophony.Lora.GetStatus")

        if call.Err != nil {
                return -1, call.Err
        }

        lc.status = 1
        return call.Body[0].(int16), nil
}

func (lc *ConnDetails) Start() (int16, error) {
        call := lc.Obj.Call("org.cacophony.Lora.Connect", 0)

        if call.Err != nil {
                return -1, call.Err
        }

        lc.Status = 1
        return call.Body[0].(int16), nil
}


func (lc *ConnDetails) Stop() {

}


func (lc *ConnDetails) WaitUntilUp(connectRequestId int16, connTimeout int16) (error) {
        err := lc.WaitUntilComplete(connectRequestId, 60)
        lc.Status = 2
        return err
}

func (lc *ConnDetails) ReportEvent(description string, times ... interface{}) (int16, error) {
  call := lc.Obj.Call("org.cacophony.Lora.Message", 0, description)
        if call.Err != nil {
                return -1, call.Err
        }
        return call.Body[0].(int16), nil
}

func (lc *ConnDetails) WaitUntilComplete(requestId int16, timeout int16) (error) {
        complete := false
        var attempts int16
        attempts = 0
        var final_result error
        final_result = errors.New("Timeout")

        for (complete == false && attempts < timeout) {

                result := lc.Obj.Call("org.cacophony.Lora.GetResponse", 0,  requestId)
                if result.Err != nil {
                        panic(result.Err)
                }

                Status := result.Body[0].(int16)
                switch Status {
                case 5:
                        complete = true
                        final_result = nil
                case 6:
                        complete = true
                        final_result = errors.New("Warning")
                case 7:
                        complete = true
                        final_result = errors.New("ERROR")

                default:
                        time.Sleep(1 * time.Second)
                }

                attempts += 1
        }

        return final_result
}


