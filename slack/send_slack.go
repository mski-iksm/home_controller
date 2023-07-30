package slack

import "github.com/slack-go/slack"

type SlackObject struct {
	SlackToken   string
	SlackChannel string
}

// 要検証
func ConstructSlackObject(slackToken string, slackChannel string) SlackObject {
	slackObject := SlackObject{
		SlackToken:   slackToken,
		SlackChannel: slackChannel,
	}
	return slackObject
}

// objectに対してメソッドを生やして使う
// senf_slackメソッドはメッセージだけをパラメータにもつ

func send_slack(slackObject SlackObject) {

	c := slack.New(slackObject.SlackToken)

	// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
	_, _, err := c.PostMessage(slackObject.SlackChannel, slack.MsgOptionText("Hello World", true))
	if err != nil {
		panic(err)
	}

}
