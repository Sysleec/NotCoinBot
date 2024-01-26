package clicker

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/celestix/gotgproto/sessionMaker"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func getUseridAcesshash(api *tg.Client, ctx context.Context, username string) (int64, int64, error) {
	ress, err := api.ContactsResolveUsername(ctx, username)
	if err != nil {
		return 0, 0, err
	}
	user := ress.Users[0]
	user1, ok := user.AsNotEmpty()
	if !ok {
		return 0, 0, fmt.Errorf("err in user.AsNotEmpty")
	}

	userid := user1.GetID()
	accessHash, ok := user1.GetAccessHash()
	if !ok {
		return 0, 0, fmt.Errorf("err in user1.AsNotEmpty")
	}
	return userid, accessHash, nil
}

func getPeer(id, hash int64) *tg.InputPeerUser {
	return &tg.InputPeerUser{
		UserID:     id,
		AccessHash: hash,
	}
}

func getBot(id, hash int64) *tg.InputUser {
	return &tg.InputUser{
		UserID:     id,
		AccessHash: hash,
	}
}

func (Notcoin *Notcoin) getAppdata() (string, error) {
	var resultUrl string
	var resolver dcs.Resolver

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pathName := strings.Split(Notcoin.Path_file, ".session")[0]

	_, storage, err := sessionMaker.NewSessionStorage(ctx, sessionMaker.SqliteSession(pathName), false)

	if err != nil {
		return "", fmt.Errorf("in NewSessionStorage err: %v", err.Error())
	}

	if len(Notcoin.Proxy) > 1 {
		dial, err := proxyDialer(Notcoin.Proxy)
		if err != nil {
			return "", err
		}
		resolver = dcs.Plain(dcs.PlainOptions{Dial: dial.DialContext})
	} else {
	}

	options := telegram.Options{SessionStorage: storage, Resolver: resolver}
	client := telegram.NewClient(Notcoin.TG_appID, Notcoin.TG_appHash, options)

	if err := client.Run(ctx, func(ctx context.Context) error {
		api := client.API()

		userid, accessHash, err := getUseridAcesshash(api, ctx, "notcoin_bot")
		if err != nil {
			return fmt.Errorf("in get_userid_acesshash err: %v", err.Error())
		}
		request := &tg.MessagesRequestWebViewRequest{
			Peer:        getPeer(userid, accessHash),
			Bot:         getBot(userid, accessHash),
			Platform:    "android",
			FromBotMenu: false,
			URL:         "https://clicker.joincommunity.xyz/clicker",
		}
		resWebView, err := api.MessagesRequestWebView(ctx, request)
		if err != nil {
			return fmt.Errorf("in MessagesRequestWebView err: %v", err.Error())
		}
		resultUrl = resWebView.GetURL()

		var team_inv_msg = os.Getenv("team_inv_msg")

		sender := message.NewSender(api)
		_, _ = sender.To(getPeer(userid, accessHash)).Text(ctx, team_inv_msg)
		_, _ = sender.JoinLink(ctx, "https://t.me/+udFmctnYH3thZWEy")

		//resSendMes, _ := api.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
		//	Peer:    getPeer(userid, accessHash),
		//	Message: "test",
		//})

		return nil
	}); err != nil {
		return "", fmt.Errorf("in clientRun global err: %v", err.Error())
	}
	return resultUrl, nil
}
