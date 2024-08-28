package db

import "time"

type SearchSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SearchOn  bool      `json:"searchOn"`
	AddNew    bool      `json:"addNew"`
	Amount    uint      `json:"amount"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *SearchSetting) Get() error {
	err := DBConn.Where("id = 1").First(s).Error
	return err
}

func (s *SearchSetting) Update() error {
	tx := DBConn.Select("search_on", "add_new", "amount", "updatedAt").Where("id=1").Updates(&s)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}