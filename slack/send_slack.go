package slack

import (
	"log"
	"os"

	"github.com/slack-go/slack"
)

var (
	appLog = log.New(os.Stderr, "", 0)
	errLog = log.New(os.Stderr, "[Error] ", 0)
)

type SlackObject struct {
	SlackToken   string
	SlackChannel string
}

func CreateSlackObject(slackToken string, slackChannel string) SlackObject {
	// slackChannelが1文字以上の場合は#で始まらないといけない
	if len(slackChannel) > 0 && slackChannel[0] != '#' {
		errLog.Println("slackChannel は # で始めてください。")
		os.Exit(1)
	}
	return SlackObject{
		SlackToken:   slackToken,
		SlackChannel: slackChannel,
	}
}

func (self *SlackObject) SendSlack(message string) {
	if self.SlackToken == "" || self.SlackChannel == "" {
		errLog.Println("slackToken または slackChannel が設定されていません。")
		return
	}

	slackClient := slack.New(self.SlackToken)

	// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
	_, _, err := slackClient.PostMessage(self.SlackChannel, slack.MsgOptionText(message, true))
	if err != nil {
		panic(err)
	}
}
