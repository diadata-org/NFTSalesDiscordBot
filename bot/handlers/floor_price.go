package handlers

import (
	"DIA-NFT-Sales-Bot/config"
	"DIA-NFT-Sales-Bot/services"
	"DIA-NFT-Sales-Bot/utils"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func FloorPriceHandler(discordSession *discordgo.Session, interaction *discordgo.InteractionCreate) {

	channel := interaction.ChannelID

	payload, found := config.ActiveSalesInfoBot[channel]

	if !found {
		_, err := discordSession.ChannelMessageSend(channel, "This Channel is not registered. Kindly use /set_up_info_bot and select this channel.")
		if err != nil {
			panic(err)
		}
		return
	}

	//Respond Channel is being Setup
	message := fmt.Sprintf("Get FloorPrice for Collection: %s ", utils.CreateHyperLink(payload["address"], utils.GetScanLink("address", payload["address"], payload["blockchain"])))
	err := discordSession.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("https://api.diadata.org/v1/NFTVolume/%s/%s", payload["address"], payload["blockchain"])
	response := services.VolumeAPI(payload["address"], payload["blockchain"], url)
	embedMsg := createFloorPriceMessage(response, payload["address"], payload["blockchain"])

	_, err = config.DiscordBot.ChannelMessageSendEmbed(interaction.ChannelID, embedMsg)
	if err != nil {
		panic(err)
	}
}

func createFloorPriceMessage(payload services.VolumeAPIResponse, address, blockchain string) *discordgo.MessageEmbed {
	scanLink := utils.GetScanLink("address", address, blockchain)
	price := fmt.Sprintf("%f %s", math.Round(payload.Floor*100)/100, currencies[strings.ToLower(blockchain)])

	embed := &discordgo.MessageEmbed{
		Color:       0x5f3267,
		Title:       "Floor Price",
		Description: "NFT Discord Bot Floor Price Response",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Collection Name",
				Value:  payload.Collection,
				Inline: false,
			}, {
				Name:   "Collection Address",
				Value:  utils.CreateHyperLink(address, scanLink),
				Inline: true,
			}, {
				Name:   "Collection FloorPrice",
				Value:  price,
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Footer:    &config.MessageFooter,
	}

	return embed
}
