package main

import "encoding/json"
import "errors"

// common dicom request
type DicomCEchoRequest struct {
	Address        string `json:"Address"`
	Port           uint16 `json:"Port"`
	ServerAE_Title string `json:"ServerAE_Title"`
}

func (dicomCEchoRequest *DicomCEchoRequest) InitFromJsonData(data []byte) error {
	err := json.Unmarshal(data, &dicomCEchoRequest)
	if err != nil {
		return errors.New("error: Can't parse dicom cEcho request data")
	}
	return nil

}

type DicomCFindRequest struct {
	Address          string `json:"Address"`
	Port             uint16 `json:"Port"`
	ServerAE_Title   string `json:"ServerAE_Title"`
	PatientName      string `json:"PatientName"`
	PatientMRN       string `json:"PatientMRN"`
	StudyID          string `json:"StudyID"`
	PatienDateOfBird string `json:"PatienDateOfBird"`
	StudyDate        string `json:"StudyDate"`
}

func (dicomCFindRequest *DicomCFindRequest) InitFromJsonData(data []byte) error {
	err := json.Unmarshal(data, &dicomCFindRequest)
	if err != nil {
		return errors.New("error: Can't parse dicom cFind request data")
	}
	return nil

}