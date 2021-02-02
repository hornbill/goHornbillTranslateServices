package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strconv"
)

func getBulletins(serviceID int, language string) ([]bulletinDetailsStruct, error) {
	var returnParams []bulletinDetailsStruct
	espXmlmc.SetParam("application", "com.hornbill.servicemanager")
	espXmlmc.SetParam("entity", "ServiceBulletin")
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
	debugLog("data::entityBrowseRecords2::ServiceBulletin")
	debugLog(XMLSTRING)

	xmlmcResponse, err := espXmlmc.Invoke("data", "entityBrowseRecords2")
	if err != nil {
		return returnParams, errors.New("API Call failed when returning Catalogs from Hornbill: " + err.Error())
	}

	var xmlRespon espGetBulletinsStruct
	err = xml.Unmarshal([]byte(xmlmcResponse), &xmlRespon)
	if err != nil {
		return returnParams, errors.New("Unmarshal failed when returning Catalogs from Hornbill: " + err.Error())
	}
	if xmlRespon.MethodResult != "ok" {
		return returnParams, errors.New("Error from Hornbill when returning Catalog data: " + xmlRespon.State.ErrorRet)
	}
	return xmlRespon.BulletinDetails, err
}

func makeBulletinExtraParams(bulletinDetails bulletinDetailsStruct) string {
	var extraParams bulletinExtraParamsStruct
	extraParams.HServiceID = bulletinDetails.HServiceID
	extraParams.HOrder = bulletinDetails.HOrder
	extraParams.HServiceBulletinImage = bulletinDetails.HServiceBulletinImage
	extraParams.HDisplayBulletinText = bulletinDetails.HDisplayBulletinText
	extraParams.HDisplayTextShadow = bulletinDetails.HDisplayTextShadow
	extraParams.HLink = bulletinDetails.HLink
	extraParams.HStartTimer = bulletinDetails.HStartTimer
	extraParams.HEndTimer = bulletinDetails.HEndTimer
	returnParams, _ := json.Marshal(extraParams)
	return string(returnParams)
}

func makeBulletinEntityObj() string {
	var entityObj entityObjStruct
	entityObj.Name = "ServiceBulletin"
	entityObj.Title = "h_bulletin_title"
	entityObj.Description = "h_bulletin_description"
	entityObj.EntityColumn = "h_id"
	entityObj.LinkedColumn = "h_service_bulletin_id"
	returnParams, _ := json.Marshal(entityObj)
	return string(returnParams)
}
