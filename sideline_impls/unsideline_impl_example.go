package sideline_impls

import "github.com/flipkart-incubator/go-dmux/sideline"

type ScanImpl struct {
}

type UnsidelineImpl struct {
}

func (s *ScanImpl) ScanWithStartRowEndRow(request sideline.ScanWithStartRowEndRowRequest) ([]string, error) {
	response := []string{"abc", "cde"}
	return response, nil
}

func (s *ScanImpl) ScanWithStartTimeEndTime(request sideline.ScanWithStartTimeEndTimeRequest) ([]string, error) {
	response := []string{"abc", "cde"}
	return response, nil
}

func (u *UnsidelineImpl) UnsidelineByKey(request sideline.UnsidelineByKeyRequest) (string, error) {
	return "success", nil
}

func unsidelineInitExample() {
	scanImpl := &ScanImpl{}
	unsidelineImpl := &UnsidelineImpl{}
	path := "" // config path
	sideline.UnsidelineStart(scanImpl, unsidelineImpl, path)
}
