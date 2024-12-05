package model

type Activity struct {
	Model
	UUID   string `validate:"token" json:"uuid"`
	Name   string `validate:"required,min=3,max=50" json:"name"`
	Url    string `validate:"omitempty,url" json:"url"`
	Type   uint8  `json:"type"`
	Status uint8  `json:"status"`
	Cost   uint   `json:"cost"`
	Info   string `validate:"omitempty,max=255" json:"info"`
}

type Prize struct {
	Model
	ActivityUUID string `json:"activity_uuid" validate:"token"`
	UUID         string `validate:"token" json:"uuid"`
	Name         string `validate:"required,min=3,max=10" json:"name"`
	Type         uint8  `json:"type"`
	Value        string `validate:"required,max=20" json:"value"`
	InitialStock uint   `json:"initial_stock" validate:"min=1"`
	Stock        uint   `validate:"ltefield=InitialStock" json:"stock"`
	Score        uint   `validate:"gt=0"`
}

type ActivityRecord struct {
	Model
	Uid          int64  `json:"uid"`
	ActivityUUID string `json:"activity_uuid"`
	PrizeUUID    string `json:"prize_uuid"`
	Remark       string `json:"remark"`
}
