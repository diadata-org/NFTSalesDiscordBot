package models

import (
	"DIA-NFT-Sales-Bot/config"
	"database/sql"
	"gorm.io/gorm"
)

type Subscriptions struct {
	gorm.Model
	Command   string         `gorm:"column:command;not null"`
	ChannelID sql.NullString `gorm:"column:channel_id"`
	Address   sql.NullString `gorm:"column:address"`
	Threshold float64        `gorm:"column:threshold"`
	All       bool           `gorm:"column:all"`
	Active    bool           `gorm:"column:is_active"`
}

func (subscription Subscriptions) SaveSubscription() {
	result := config.DBClient.Model(&subscription).Create(&subscription)

	if result.Error != nil {
		err := "Error Saving Subscription: \n" + result.Error.Error()
		panic(err)
	}
}

func (subscription Subscriptions) LoadChannelSubscriptions() []Subscriptions {
	var subscriptions []Subscriptions

	result := config.DBClient.Model(&subscription).Where("is_active = ? AND channel_id = ?", true, subscription.ChannelID.String).Find(subscriptions)
	if result.Error != nil {
		err := "Error Fetching Subscriptions by Channel: \n" + result.Error.Error()
		panic(err)
	}

	return subscriptions
}

func LoadAllSubscriptions() {}