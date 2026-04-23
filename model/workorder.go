package model

import (
	"encoding/json"
	"time"
)

// Generic request wrapper untuk semua function
type GenericRequest struct {
	Function  string          `json:"function"`
	Payload   json.RawMessage `json:"payload"`
	RequestID string          `json:"-"` // Generated server-side
}

// Specific payload untuk ff_updateWorkorder
type WorkorderRequest struct {
	Req  RequestData  `json:"req"`
	Res  ResponseData `json:"res"`
	Date time.Time    `json:"date"`
}

type RequestData struct {
	Memo               string `json:"memo"`
	Task               string `json:"task"`
	Wolo1              string `json:"wolo1"`
	Wolo3              string `json:"wolo3"`
	Wonum              string `json:"wonum"`
	Siteid             string `json:"siteid"`
	Status             string `json:"status"`
	Latitude           string `json:"latitude"`
	CpeModel           string `json:"cpe_model"`
	Errorcode          string `json:"errorcode"`
	Longitude          string `json:"longitude"`
	TaskName           string `json:"task_name"`
	CpeVendor          string `json:"cpe_vendor"`
	LaborScmt          string `json:"labor_scmt"`
	Statusiface        string `json:"statusiface"`
	Urlevidence        string `json:"urlevidence"`
	Engineermemo       string `json:"engineermemo"`
	Suberrorcode       string `json:"suberrorcode"`
	NpStatusmemo       string `json:"np_statusmemo"`
	CpeSerialNumber    string `json:"cpe_serial_number"`
	TkCustomHeader03   string `json:"tk_custom_header_03"`
	TkCustomHeader04   string `json:"tk_custom_header_04"`
	TkCustomHeader09   string `json:"tk_custom_header_09"`
	TkCustomHeader10   string `json:"tk_custom_header_10"`
}

type ResponseData struct {
	Data    string      `json:"data"`
	Errors  interface{} `json:"errors"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
}

// Standard API Response dengan RequestID
type APIResponse struct {
	RequestID string      `json:"requestId"`
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
}
