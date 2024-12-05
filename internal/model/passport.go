package model

type Passport struct {
	Model
	Uid      int64  `json:"-"`
	Token    string `json:"-"`
	Ip       string `json:"ip"`
	DeviceId string `json:"device_id"`
	Ua       string `json:"ua"`
}

type PassportDto struct {
	Uid      int64  `json:"uid"`
	DeviceId string `json:"device_id"`
}
