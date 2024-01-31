package clicker

import (
	"io"
	"time"
)

type Webappses_resp struct {
	Ok   bool `json:"ok"`
	Data struct {
		UserID       int    `json:"userId"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	} `json:"data"`
}

type Accounts_struct struct {
	APIID   int    `json:"api_id"`
	APIHash string `json:"api_hash"`
	Proxy   string `json:"proxy"`
}

type Accounts_data map[string]Accounts_struct

type Click_respdata struct {
	Id                 int         `json:"id"`
	UserId             int         `json:"userId"`
	TeamId             int         `json:"teamId"`
	LeagueId           int         `json:"leagueId"`
	LimitCoins         int         `json:"limitCoins"`
	TotalCoins         string      `json:"totalCoins"`
	BalanceCoins       string      `json:"balanceCoins"`
	SpentCoins         string      `json:"spentCoins"`
	MiningPerTime      int         `json:"miningPerTime"`
	MultipleClicks     int         `json:"multipleClicks"`
	AutoClicks         int         `json:"autoClicks"`
	WithRobot          bool        `json:"withRobot"`
	LastMiningAt       string      `json:"lastMiningAt"`
	LastAvailableCoins int         `json:"lastAvailableCoins"`
	TurboTimes         int         `json:"turboTimes"`
	Avatar             interface{} `json:"avatar"`
	CreatedAt          string      `json:"createdAt"`
	Hash               []string    `json:"hash"`
	AvailableCoins     int         `json:"availableCoins"`
}

type Click_resp struct {
	Ok   bool             `json:"ok"`
	Data []Click_respdata `json:"data"`
}

type Click_resp_no_slice struct {
	Ok   bool           `json:"ok"`
	Data Click_respdata `json:"data"`
}

type Active_turbo_data struct {
	Multiple int   `json:"multiple"`
	Expire   int64 `json:"expire"`
}

type Active_turbo_resp struct {
	Ok   bool                `json:"ok"`
	Data []Active_turbo_data `json:"data"`
}

type Notcoin struct {
	Ses   Session
	Proxy string

	Clear_name   string
	Path_file    string
	TG_appID     int
	TG_appHash   string
	TGWebAppData string

	Need_sleep_10 bool
	//UserId        int
	UserId      string
	Coefficient int // multipleClicks

	Fullenergy_boost   int
	Count_400          int
	LimitCoins         int
	LastAvailableCoins int
	BalanceCoins       string
	Turbo              bool
	Timestart_turbo    int64
	Turbo_boost_count  int
	Hash               int
}

type Task_completed_data struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Coins     int       `json:"coins"`
	Max       int       `json:"max"`
	Link      string    `json:"link"`
	Image     string    `json:"image"`
	Status    string    `json:"status"`
	IsDaily   bool      `json:"isDaily"`
	EntityId  string    `json:"entityId"`
	CreatedAt time.Time `json:"createdAt"`
}

type Tasks_completed struct {
	Id        int                 `json:"id"`
	TaskId    int                 `json:"taskId"`
	UserId    int                 `json:"userId"`
	CreatedAt time.Time           `json:"createdAt"`
	Task      Task_completed_data `json:"task"`
}

type Task_completed_resp struct {
	Ok   bool              `json:"ok"`
	Data []Tasks_completed `json:"data"`
}

type Itemstore_data struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Type              string    `json:"type"`
	Image             string    `json:"image"`
	Coins             int       `json:"coins"`
	Price             int       `json:"price"`
	Max               int       `json:"max"`
	Coefficient       int       `json:"coefficient"`
	IsPartner         bool      `json:"isPartner"`
	IsTask            bool      `json:"isTask"`
	IsFeatured        bool      `json:"isFeatured"`
	Status            string    `json:"status"`
	MinLeagueID       int       `json:"minLeagueId"`
	ChallengeID       int       `json:"challengeId"`
	LiveTimeInSeconds int       `json:"liveTimeInSeconds"`
	CreatedAt         time.Time `json:"createdAt"`
	IsCompleted       bool      `json:"isCompleted"`
	Count             int       `json:"count"`
}

type Store_resp struct {
	OK   bool             `json:"ok"`
	Data []Itemstore_data `json:"data"`
}

type Item_buying struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"userId"`
	TeamID             int       `json:"teamId"`
	LeagueID           int       `json:"leagueId"`
	LimitCoins         int       `json:"limitCoins"`
	TotalCoins         string    `json:"totalCoins"`
	BalanceCoins       string    `json:"balanceCoins"`
	SpentCoins         string    `json:"spentCoins"`
	MiningPerTime      int       `json:"miningPerTime"`
	MultipleClicks     int       `json:"multipleClicks"`
	AutoClicks         int       `json:"autoClicks"`
	WithRobot          bool      `json:"withRobot"`
	LastMiningAt       time.Time `json:"lastMiningAt"`
	LastAvailableCoins int       `json:"lastAvailableCoins"`
	TurboTimes         int       `json:"turboTimes"`
	Avatar             string    `json:"avatar"`
	CreatedAt          time.Time `json:"createdAt"`
}

type Buy_item_resp struct {
	OK   bool        `json:"ok"`
	Data Item_buying `json:"data"`
}

type ColorableWriter struct {
	Console io.Writer
	File    io.Writer
	Prefix  string
	Color   func(a ...interface{}) string
}
