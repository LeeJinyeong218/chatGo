package model

type TestTable struct {
	ID   int64 `gorm:"primaryKey"`
	Text string
}

func (TestTable) TableName() string {
	return "test_table"
}
