package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strconv"
)

func getStatuses(serviceID int, language string) ([]statusDetailsStruct, error) {
	var returnParams []statusDetailsStruct
	espXmlmc.SetParam("application", "com.hornbill.servicemanager")
	espXmlmc.SetParam("entity", "Statuses")
	espXmlmc.OpenElement("searchFilter")
	espXmlmc.SetParam("column", "h_language")
	espXmlmc.SetParam("value", language)
	espXmlmc.SetParam("matchType", "exact")
	espXmlmc.CloseElement("searchFilter")
	espXmlmc.OpenElement("searchFilter")
	espXmlmc.SetParam("column", "h_service_id")
	espXmlmc.SetParam("value", strconv.Itoa(serviceID))
	espXmlmc.SetParam("matchType", "exact")
	espXmlmc.CloseElement("searchFilter")

	var XMLSTRING = espXmlmc.GetParam()
	debugLog("data::entityBrowseRecords2::Statuses")
	debugLog(XMLSTRING)

	xmlmcResponse, err := espXmlmc.Invoke("data", "entityBrowseRecords2")
	if err != nil {
		return returnParams, errors.New("API Call failed when returning Statuses from Hornbill: " + err.Error())
	}

	var xmlRespon espGetStatusStruct
	err = xml.Unmarshal([]byte(xmlmcResponse), &xmlRespon)
	if err != nil {
		return returnParams, errors.New("Unmarshal failed when returning Statuses from Hornbill: " + err.Error())
	}
	if xmlRespon.MethodResult != "ok" {
		return returnParams, errors.New("Error from Hornbill when returning Statuses data: " + xmlRespon.State.ErrorRet)
	}
	return xmlRespon.StatusDetails, err
}

func makeStatusExtraParams(statusDetails statusDetailsStruct) string {
	var extraParams statusExtraParamsStruct
	extraParams.HServiceID = statusDetails.HServiceID
	extraParams.HRequestType = statusDetails.HRequestType
	extraParams.HParentStatus = statusDetails.HParentStatus
	extraParams.HPauseIndef = statusDetails.HPauseIndef
	extraParams.HReasonRequired = statusDetails.HReasonRequired
	extraParams.HTimelineVisibility = statusDetails.HTimelineVisibility
	extraParams.HSupplierEnabled = statusDetails.HSupplierEnabled
	extraParams.HStatus = statusDetails.HStatus
	extraParams.HDatePublished = statusDetails.HDatePublished
	returnParams, _ := json.Marshal(extraParams)
	return string(returnParams)
}

func makeStatusEntityObj() string {
	var entityObj entityObjStruct
	entityObj.Name = "Statuses"
	entityObj.Title = "h_name"
	entityObj.Description = "h_customer_label"
	entityObj.EntityColumn = "h_id"
	entityObj.LinkedColumn = "h_status_id"
	returnParams, _ := json.Marshal(entityObj)
	return string(returnParams)
}
