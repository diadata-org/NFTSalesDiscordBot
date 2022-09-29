package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/models"
	"database/sql"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func StopAllHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {
	optionsMap := ParseCommandOptions(interaction)

	if optionsMap["channel"] != nil {
		channel := optionsMap["channel"].ChannelValue(nil)
		sub := models.Subscriptions{ChannelID: sql.NullString{
			String: channel.ID,
			Valid:  true,
		}}

		channelSubs := sub.LoadChannelSubscriptions()

		for _, channelSub := range channelSubs {
			switch channelSub.Command {
			case "sales":
				for index, channel := range config.ActiveSales[channelSub.Address.String] {
					if channel == channelSub.ChannelID.String {
						//Update Golang to 1.18
						config.ActiveSales[channelSub.Address.String] = slices.Delete(config.ActiveSales[channelSub.Address.String], index, index+1)
						break
					}
				}
			case "all_sales":
				for index, channel := range config.ActiveAllSales[channelSub.Threshold] {
					if channel == channelSub.ChannelID.String {
						config.ActiveAllSales[channelSub.Threshold] = slices.Delete(config.ActiveAllSales[channelSub.Threshold], index, index+1)
						break
					}
				}
			}
		}
	} else {
		go models.Subscriptions{}.DeactivateAllSubscriptions()
		// Delete Global variables
		go maps.Clear(config.ActiveAllSales)
		config.ActiveAllSalesKeys = append(config.ActiveAllSalesKeys[:1], config.ActiveAllSalesKeys[2:]...)
		go maps.Clear(config.ActiveSales)
		go config.NftEventWSCancelFunc()
	}
}