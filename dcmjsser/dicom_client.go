package main

import "errors"
import "dtools/gdcmgobr"
import "encoding/json"

type DClient struct {
	CallerAE_Title string
}

func (dc *DClient) checRequisites() error {
	if len(dc.CallerAE_Title) == 0 {
		return errors.New("error: CallerAE_Title is empty")
	}
	return nil
}

func (dc *DClient) CStore(cStorReq CStorReq) (CStorReq, error) {
	if err := dc.checRequisites(); err != nil {
		return CStorReq{}, err
	}
	cae := dc.CallerAE_Title
	sae := cStorReq.ServerSet.ServerAE_Title
	ip := cStorReq.ServerSet.Address
	port := cStorReq.ServerSet.Port
	fls := cStorReq.File
	isStore := gdcmgobr.CStore(ip, port, sae, cae, fls)
	if !isStore {
		return CStorReq{}, errors.New("error: can't store dicom file " + fls)
	}
	return cStorReq, nil
}

func (dc *DClient) CGet(cgt CGetReq) (CGetReq, error) {
	if err := dc.checRequisites(); err != nil {
		return CGetReq{}, err
	}
	cae := dc.CallerAE_Title
	sae := cgt.FindReq.ServerSet.ServerAE_Title
	ip := cgt.FindReq.ServerSet.Address
	port := cgt.FindReq.ServerSet.Port
	pn := cgt.FindReq.PatientName
	an := cgt.FindReq.AccessionNumber
	bd := cgt.FindReq.PatienDateOfBirth
	sd := cgt.FindReq.StudyDate
	fp := cgt.Folder
	cget := gdcmgobr.CGet(sae, cae, ip, port, pn, an, bd, sd, fp)
	if !cget {
		return CGetReq{}, errors.New("error: can't cget dicom file " + pn)
	}
	return cgt, nil
}

func (dc *DClient) CFind(freq FindReq) ([]FindRes, error) {
	if err := dc.checRequisites(); err != nil {
		return nil, err
	}
	cae := dc.CallerAE_Title
	sae := freq.ServerSet.ServerAE_Title
	ip := freq.ServerSet.Address
	port := freq.ServerSet.Port
	pn := freq.PatientName
	an := freq.AccessionNumber
	bd := freq.PatienDateOfBirth
	sd := freq.StudyDate
	cfindResult := gdcmgobr.CFind(sae, cae, ip, port, pn, an, bd, sd)
	var fdat []FindRes
	err := json.Unmarshal([]byte(cfindResult), &fdat)
	if err != nil {
		return nil, errors.New("error: can't parse dicom cFind result data " + err.Error() + cfindResult)
	}
	return fdat, nil
}
func (dc *DClient) CEcho(dicomCEchoRequest EchoReq) (EchoRes, error) {
	if err := dc.checRequisites(); err != nil {
		return EchoRes{}, err
	}
	//	log.Printf("info: dicomCEchoRequest.ServerAE_Title=%v dc.CallerAE_Title=%v ",dicomCEchoRequest.ServerAE_Title, dc.CallerAE_Title)
	isAlive := gdcmgobr.CEcho(dicomCEchoRequest.Address, dicomCEchoRequest.Port, dicomCEchoRequest.ServerAE_Title, dc.CallerAE_Title)
	dicomCEchoResult := EchoRes{IsAlive: isAlive}
	return dicomCEchoResult, nil
}
