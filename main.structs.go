package main

import (
	apiLib "github.com/hornbill/goApiLib"
	"github.com/sirupsen/logrus"
)

const (
	version    = "1.1.0"
	minSMBuild = 2144
)

var (
	logr    = logrus.New()
	logFile = logrus.New()

	configAPIKey      string
	configDebug       bool
	configDestination string
	configInstance    string
	configLangs       bool
	configSource      string
	configVersion     bool

	counters counterStruct

	espXmlmc *apiLib.XmlmcInstStruct

	srcServices []smServiceStruct
	dstServices []smServiceStruct
	languages   = make(map[string]langsStruct)

	smStatusSupport = false
)

// General Structs

type counterStruct struct {
	ServicesTranslated int
	ServicesSkipped    int
	ServicesError      int
	CatalogsCreated    int
	CatalogsSkipped    int
	CatalogsError      int
	BulletinsCreated   int
	BulletinsSkipped   int
	BulletinsError     int
	FAQsCreated        int
	FAQsSkipped        int
	FAQsError          int
	FeedbackCreated    int
	FeedbackSkipped    int
	FeedbackError      int
	StatusCreated      int
	StatusSkipped      int
	StatusError        int
}

type entityObjStruct struct {
	Description  string `json:"description"`
	EntityColumn string `json:"entityColumn"`
	LinkedColumn string `json:"linkedColumn"`
	Name         string `json:"name"`
	Title        string `json:"title"`
}

type espStateStruct struct {
	Code     string `xml:"code"`
	ErrorRet string `xml:"error"`
}

//-- Service Records Structs

type espGetServicesStruct struct {
	MethodResult string            `xml:"status,attr"`
	State        espStateStruct    `xml:"state"`
	Rows         []smServiceStruct `xml:"params>rowData>row"`
}

type smServiceStruct struct {
	HPkServiceid                    int    `xml:"h_pk_serviceid"`
	HLinkedServiceID                int    `xml:"h_linked_service_id"`
	HServicename                    string `xml:"h_servicename"`
	HServicedescription             string `xml:"h_servicedescription"`
	HAvailable                      string `xml:"h_available"`
	HChangeBpmName                  string `xml:"h_change_bpm_name"`
	HAccess                         string `xml:"h_access"`
	HHassupportteam                 string `xml:"h_hassupportteam"`
	HIncidentActionVisibility       string `xml:"h_incident_action_visibility"`
	HIncidentBpmName                string `xml:"h_incident_bpm_name"`
	HPortfolioStatus                string `xml:"h_portfolio_status"`
	HFkServicecategory              string `xml:"h_fk_servicecategory"`
	HServiceBpmName                 string `xml:"h_service_bpm_name"`
	HAllowIncidents                 string `xml:"h_allow_incidents"`
	HAllowServicerequests           string `xml:"h_allow_servicerequests"`
	HAllowChangerequests            string `xml:"h_allow_changerequests"`
	HAllowProblems                  string `xml:"h_allow_problems"`
	HAllowKnownerrors               string `xml:"h_allow_knownerrors"`
	HAllowReleases                  string `xml:"h_allow_releases"`
	HServiceIncidentPortalHud       string `xml:"h_service_incident_portal_hud"`
	HServiceServicerequestPortalHud string `xml:"h_service_servicerequest_portal_hud"`
	HStatus                         string `xml:"h_status"`
	HIcon                           string `xml:"h_icon"`
	HLanguage                       string `xml:"h_language"`
	HDefaultLanguage                string `xml:"h_default_language"`
	HLastUpdated                    string `xml:"h_last_updated"`
}

type smServiceDetailsStruct struct {
	CatalogCategory              string `json:"catalog_category"`
	CatalogDomain                string `json:"catalog_domain"`
	HAccess                      string `json:"h_access"`
	HAllowChangerequests         string `json:"h_allow_changerequests"`
	HAllowIncidents              string `json:"h_allow_incidents"`
	HAllowKnownerrors            string `json:"h_allow_knownerrors"`
	HAllowProblems               string `json:"h_allow_problems"`
	HAllowReleases               string `json:"h_allow_releases"`
	HAllowServicerequests        string `json:"h_allow_servicerequests"`
	HAvailable                   string `json:"h_available"`
	HDateCreated                 string `json:"h_date_created"`
	HFkServicecategory           string `json:"h_fk_servicecategory"`
	HHassupportteam              string `json:"h_hassupportteam"`
	HIcon                        string `json:"h_icon"`
	HIconRef                     string `json:"h_icon_ref"`
	HLastUpdated                 string `json:"h_last_updated"`
	HLinkedServiceID             string `json:"h_linked_service_id"`
	HPkServiceid                 string `json:"h_pk_serviceid"`
	HPortfolioStatus             string `json:"h_portfolio_status"`
	HRelType                     string `json:"h_rel_type"`
	HServiceCatalogCategory      string `json:"h_service_catalog_category"`
	HServiceCatalogDomain        string `json:"h_service_catalog_domain"`
	HServicestatus               string `json:"h_servicestatus"`
	HStatus                      string `json:"h_status"`
	HUserID                      string `json:"h_user_id"`
	HUserName                    string `json:"h_user_name"`
	Iscurrentuserownerordelegate string `json:"iscurrentuserownerordelegate"`
	ServiceDescription           string `json:"service_description"`
	ServiceName                  string `json:"service_name"`
	TranslatedID                 string `json:"translated_id"`
	HMbid                        string `json:"h_mbid"`
}

type serviceExtraParamsStruct struct {
	HAccess                 string `json:"h_access"`
	HAvailable              string `json:"h_available"`
	HFkServicecategory      string `json:"h_fk_servicecategory"`
	HIcon                   string `json:"h_icon"`
	HMbid                   string `json:"h_mbid"`
	HPortfolioStatus        string `json:"h_portfolio_status"`
	HServiceCatalogCategory int    `json:"h_service_catalog_category"`
	HServiceCatalogDomain   string `json:"h_service_catalog_domain"`
	HStatus                 string `json:"h_status"`
}

//-- Translation Structs
type espGetTranslationStruct struct {
	MethodResult string                  `xml:"status,attr"`
	State        espStateStruct          `xml:"state"`
	Params       translationParamsStruct `xml:"params"`
}

type translationParamsStruct struct {
	TranslatedDescription string        `xml:"translatedDescription"`
	TranslatedTitle       string        `xml:"translatedTitle"`
	LanguageInfo          []langsStruct `xml:"languageInfo"`
	ServiceDetails        string        `xml:"queryExecJSON"`
}

type langsStruct struct {
	Language  string `xml:"language"`
	Name      string `xml:"languageName"`
	Supported bool   `xml:"supported"`
}

//-- Catalog Structs

type espGetCatalogsStruct struct {
	MethodResult string                 `xml:"status,attr"`
	State        espStateStruct         `xml:"state"`
	CatalogData  []catalogDetailsStruct `xml:"params>rowData>row"`
}

type catalogDetailsStruct struct {
	HID                 string `xml:"h_id"`
	HBpm                string `xml:"h_bpm"`
	HIcon               string `xml:"h_icon"`
	HProCapture         string `xml:"h_pro_capture"`
	HRequestType        string `xml:"h_request_type"`
	HServiceID          string `xml:"h_service_id"`
	HVisibility         string `xml:"h_visibility"`
	HLanguage           string `xml:"h_language"`
	HRequestCatalogID   int    `xml:"h_request_catalog_id"`
	HCatalogTitle       string `xml:"h_catalog_title"`
	HCatalogDescription string `xml:"h_catalog_description"`
}

type catalogExtraParamsStruct struct {
	HID          string `json:"h_id"`
	HBpm         string `json:"h_bpm"`
	HIcon        string `json:"h_icon"`
	HProCapture  string `json:"h_pro_capture"`
	HRequestType string `json:"h_request_type"`
	HServiceID   string `json:"h_service_id"`
	HVisibility  string `json:"h_visibility"`
}

//-- Bulletin Structs

type espGetBulletinsStruct struct {
	MethodResult    string                  `xml:"status,attr"`
	State           espStateStruct          `xml:"state"`
	BulletinDetails []bulletinDetailsStruct `xml:"params>rowData>row"`
}

type bulletinDetailsStruct struct {
	HID                   string `xml:"h_id"`
	HServiceBulletinID    int    `xml:"h_service_bulletin_id"`
	HBulletinTitle        string `xml:"h_bulletin_title"`
	HBulletinDescription  string `xml:"h_bulletin_description"`
	HStatus               string `xml:"h_status"`
	HLanguage             string `xml:"h_language"`
	HServiceID            string `xml:"h_service_id"`
	HOrder                string `xml:"h_order"`
	HServiceBulletinImage string `xml:"h_service_bulletin_image"`
	HDisplayBulletinText  string `xml:"h_display_bulletin_text"`
	HDisplayTextShadow    string `xml:"h_display_text_shadow"`
	HLink                 string `xml:"h_link"`
	HStartTimer           string `xml:"h_start_timer"`
	HEndTimer             string `xml:"h_end_timer"`
}

type bulletinExtraParamsStruct struct {
	HServiceID            string `json:"h_service_id"`
	HOrder                string `json:"h_order"`
	HServiceBulletinImage string `json:"h_service_bulletin_image"`
	HDisplayBulletinText  string `json:"h_display_bulletin_text"`
	HDisplayTextShadow    string `json:"h_display_text_shadow"`
	HLink                 string `json:"h_link"`
	HStartTimer           string `json:"h_start_timer"`
	HEndTimer             string `json:"h_end_timer"`
}

//-- FAQ Structs

type espGetFAQsStruct struct {
	MethodResult string             `xml:"status,attr"`
	State        espStateStruct     `xml:"state"`
	FAQDetails   []faqDetailsStruct `xml:"params>rowData>row"`
}

type faqDetailsStruct struct {
	HID              string `xml:"h_id"`
	HFAQID           int    `xml:"h_faq_id"`
	HFAQQuestion     string `xml:"h_question"`
	HFAQAnswer       string `xml:"h_answer"`
	HStatus          string `xml:"h_status"`
	HLanguage        string `xml:"h_language"`
	HServiceID       string `xml:"h_service_id"`
	HServiceName     string `xml:"h_service_name"`
	HViewCount       string `xml:"h_view_count"`
	HMediaLink       string `xml:"h_media_link"`
	HVisibility      string `xml:"h_visibility"`
	HCreatedByUserID string `xml:"h_createdby_userid"`
}

type faqExtraParamsStruct struct {
	HServiceID       string `json:"h_service_id"`
	HServiceName     string `json:"h_service_name"`
	HViewCount       string `json:"h_view_count"`
	HMediaLink       string `json:"h_media_link"`
	HVisibility      string `json:"h_visibility"`
	HCreatedByUserID string `json:"h_createdby_userid"`
}

//-- Feedback Structs

type espGetFeedbackStruct struct {
	MethodResult    string                  `xml:"status,attr"`
	State           espStateStruct          `xml:"state"`
	FeedbackDetails []feedbackDetailsStruct `xml:"params>rowData>row"`
}

type feedbackDetailsStruct struct {
	HID            string `xml:"h_id"`
	HQuestionID    int    `xml:"h_question_id"`
	HQuestion      string `xml:"h_question"`
	HRequestType   string `xml:"h_requesttype"`
	HServiceID     string `xml:"h_service_id"`
	HFieldRequired string `xml:"h_field_required"`
	HFieldType     string `xml:"h_field_type"`
	HLanguage      string `xml:"h_language"`
}

type feedbackExtraParamsStruct struct {
	HRequestType   string `json:"h_requesttype"`
	HServiceID     string `json:"h_service_id"`
	HFieldRequired string `json:"h_field_required"`
	HFieldType     string `json:"h_field_type"`
}

//-- Sub Status Structs
type espGetStatusStruct struct {
	MethodResult  string                `xml:"status,attr"`
	State         espStateStruct        `xml:"state"`
	StatusDetails []statusDetailsStruct `xml:"params>rowData>row"`
}

type statusDetailsStruct struct {
	HID                 string `xml:"h_id"`
	HStatusID           int    `xml:"h_status_id"`
	HStatus             string `xml:"h_status"`
	HName               string `xml:"h_name"`
	HCustomerLabel      string `xml:"h_customer_label"`
	HRequestType        string `xml:"h_request_type"`
	HServiceID          string `xml:"h_service_id"`
	HParentStatus       string `json:"h_parent_status"`
	HPauseIndef         string `json:"h_pause_indef"`
	HReasonRequired     string `json:"h_reason_required"`
	HTimelineVisibility string `json:"h_timeline_visibility"`
	HSupplierEnabled    string `json:"h_supplier_enabled"`
	HLanguage           string `xml:"h_language"`
	HDatePublished      string `xml:"h_date_published"`
}

type statusExtraParamsStruct struct {
	HRequestType        string `json:"h_request_type"`
	HServiceID          string `json:"h_service_id"`
	HParentStatus       string `json:"h_parent_status"`
	HPauseIndef         string `json:"h_pause_indef"`
	HReasonRequired     string `json:"h_reason_required"`
	HTimelineVisibility string `json:"h_timeline_visibility"`
	HSupplierEnabled    string `json:"h_supplier_enabled"`
	HStatus             string `json:"h_status"`
	HDatePublished      string `json:"h_date_published"`
}

//-- App Version Structs
type espGetAppVerStruct struct {
	MethodResult string         `xml:"status,attr"`
	State        espStateStruct `xml:"state"`
	Params       struct {
		Application []struct {
			Name  string `xml:"name"`
			Build int    `xml:"build"`
		} `xml:"application"`
	} `xml:"params"`
}
