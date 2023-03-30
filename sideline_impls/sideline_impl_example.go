package sideline_impls

import (
	"encoding/json"
	"github.com/flipkart-incubator/go-dmux/sideline"
)

type CheckMessageSidelineImpl struct {
}

func sidelineInitExample() {
	custom := DmuxCustom{}
	sidelineImpl := &CheckMessageSidelineImpl{}
	path := "" // config path
	custom.DmuxStart(path, sidelineImpl)
}

func (c *CheckMessageSidelineImpl) CheckMessageSideline(key []byte) ([]byte, error) {
	checkMessageSidelineResponse := sideline.CheckMessageSidelineResponse{}
	return json.Marshal(checkMessageSidelineResponse)
}

func (c *CheckMessageSidelineImpl) SidelineMessage(msg []byte) sideline.SidelineMessageResponse {
	sidelineMessageResponse := sideline.SidelineMessageResponse{}
	return sidelineMessageResponse
}

func (c *CheckMessageSidelineImpl) InitialisePlugin(confBytes []byte) error {
	return nil
}
