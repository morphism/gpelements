package gpelements

import "encoding/xml"

func NewElements() *Elements {
	return &Elements{
		OMMId:      "CCSDS_OMM_VERS",
		OMMVersion: "2.0",
	}
}

type ElementsList struct {
	XMLName xml.Name   `xml:"ndm"`
	Es      []Elements `xml:"omm"`
}
