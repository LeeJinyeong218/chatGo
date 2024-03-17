package dto

type TestTableDto struct {
	ID   int64  `json:"id" form:"id"`
	Text string `json:"text" form:"text"`
}
