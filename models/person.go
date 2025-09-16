package models

type Person struct {
	UID      string `json:"uid"` // UUID
}

type Operator struct {
	ID   string `json:"id"` // id
}

type Transaction struct {
	ID          string `json:"id"` // Id
	Person      `json:"person"`
	IpAddress   string   `json:"ip_address"`
	Timestamp   string   `json:"timestamp"`
	Operator_id string   `json:"operator_id"` // Operator
	Machine_id  string   `json:"machine_id"`  // Machine
	Devices     []Device `json:"devices"`     // Devices
}

type EA struct {
	ID   string `json:"id"` // Id
	Name string `json:"name"`
	// Location string `json:"location"`
}

type Station struct {
	ID string `json:"id"` // Id
	EA string `json:"ea"` // EA id
}

type Machine struct {
	ID            string `json:"id"`           // Id
	Serial_Number string `json:"serial"`       // Serial Number
	Installation  string `json:"installation"` // Installation date
	Station       string `json:"station"`      // Station id
}

type Device struct {
	ID   string `json:"id"` // Id
	Type string `json:"type"`
}
