package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strconv"
)

func getCatalogs(serviceID int, language string) ([]catalogDetailsStruct, error) {
	var returnParams []catalogDetailsStruct
	espXmlmc.SetParam("application", "com.hornbill.servicemanager")
	espXmlmc.SetParam("entity", "Catalogs")
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
	debugLog("data::entityBrowseRecords2::Catalogs")
	debugLog(XMLSTRING)

	xmlmcResponse, err := espXmlmc.Invoke("data", "entityBrowseRecords2")
	if err != nil {
		return returnParams, errors.New("API Call failed when returning Catalogs from Hornbill: " + err.Error())
	}

	var xmlRespon espGetCatalogsStruct
	err = xml.Unmarshal([]byte(xmlmcResponse), &xmlRespon)
	if err != nil {
		return returnParams, errors.New("Unmarshal failed when returning Catalogs from Hornbill: " + err.Error())
	}
	if xmlRespon.MethodResult != "ok" {
		return returnParams, errors.New("Error from Hornbill when returning Catalog data: " + xmlRespon.State.ErrorRet)
	}
	return xmlRespon.CatalogData, err
}

func makeCatalogExtraParams(catalogDetails catalogDetailsStruct) string {
	var extraParams catalogExtraParamsStruct
	extraParams.HID = catalogDetails.HID
	extraParams.HBpm = catalogDetails.HBpm
	extraParams.HIcon = catalogDetails.HIcon
	extraParams.HProCapture = catalogDetails.HProCapture
	extraParams.HRequestType = catalogDetails.HRequestType
	extraParams.HServiceID = catalogDetails.HServiceID
	extraParams.HVisibility = catalogDetails.HVisibility
	returnParams, _ := json.Marshal(extraParams)
	return string(returnParams)
}

func makeCatalogEntityObj() string {
	var entityObj entityObjStruct
	entityObj.Name = "Catalogs"
	entityObj.Title = "h_catalog_title"
	entityObj.Description = "h_catalog_description"
	entityObj.EntityColumn = "h_id"
	entityObj.LinkedColumn = "h_request_catalog_id"
	returnParams, _ := json.Marshal(entityObj)
	return string(returnParams)
}
