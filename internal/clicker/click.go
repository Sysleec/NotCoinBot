package clicker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/robertkrimen/otto"
)

func checkHashShit(checkCode string) (bool, int) {
	shit0 := "querySelectorAll"
	shit1 := "window.Telegram.WebApp.initDataUnsafe.user.id"
	shit11 := "? 5 : 10"
	shit2 := "window.location.host == 'clicker.joincommunity.xyz' ? 129 : 578"

	if strings.Contains(checkCode, shit1) && strings.Contains(checkCode, shit11) {
		return true, 5 // 5||10
	}
	if strings.Contains(checkCode, shit0) {
		return true, 1
	}
	if strings.Contains(checkCode, shit2) {
		return true, 129
	}

	return false, 0
}

func getNeedplus(str string) (string, int) {
	splitstr := strings.Split(str, " ")
	last := splitstr[len(splitstr)-1]
	plus, err := strconv.Atoi(last)
	if err != nil {
		ErrorLogger.Println("Error on get_needplus, err =", err)
	}
	firstStr := strings.Join(splitstr[:len(splitstr)-1], " ")
	return firstStr, plus
}

func hashResolve(resp []string) int {
	var needPlus int
	encodedString := strings.Join(resp, "")
	if len(encodedString) < 3 {
		return -1
	}

	decodedBytes, _ := base64.StdEncoding.DecodeString(encodedString)
	codejsStr := string(decodedBytes)

	if len(resp) > 1 {
		codejsStr, needPlus = getNeedplus(codejsStr)
	}

	findshit, shit := checkHashShit(codejsStr)
	if findshit {
		return shit + needPlus
	}
	vm := otto.New()
	result, _ := vm.Run(codejsStr)
	resultint, _ := result.ToInteger()
	return int(resultint) + needPlus

}

func getDivide(coef, divide int) int {
	divided := coef / divide
	if divided < 1 {
		divided = 1
	}
	return divided
}

func (Notcoin *Notcoin) getCountClick() int { // in hand 40/sec, turbo = hand*3
	var coinscount int
	var minus int
	if Notcoin.Turbo {
		//conv to int
		// limCoinsINT, _ := strconv.Atoi(Notcoin.LimitCoins)
		minus = getRandomint(132, 311, 1)
		coinscount = Notcoin.LimitCoins/4 - minus
		if Notcoin.Timestart_turbo == 0 {
			Notcoin.Timestart_turbo = time.Now().Unix()
		} else if Notcoin.Timestart_turbo+11 <= time.Now().Unix() {
			Notcoin.Turbo = false
			Notcoin.Timestart_turbo = 0
		}
	} else {
		minus = getDivide(Notcoin.Coefficient, 3)
		minn := 100*Notcoin.Coefficient/getDivide(Notcoin.Coefficient, 3) + getRandomint(2, 47, 1)
		if minus <= 2 {
			coinscount = minn
		} else {
			coinscount = Notcoin.LastAvailableCoins / minus
		}

		if coinscount <= minn && Notcoin.LastAvailableCoins > minn {
			coinscount = minn
		} else if coinscount < minn {
			coinscount = Notcoin.LastAvailableCoins
		}
	}

	return coinscount
}

func getRandomint(min, max, coef int) int {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(max-min+1) + min
	return randomInt * coef
}

//func get_Turbo(resp string) bool {
//	re := regexp.MustCompile(`"Turbo":(.*?)}`)
//	match := re.FindStringSubmatch(resp)
//	if len(match) > 1 {
//		if strings.Contains(strings.ToLower(match[1]), "true") {
//			return true
//		}
//	}
//	return false
//}

func parseRespclick(content []byte) *Click_resp {
	var response Click_resp
	err := json.Unmarshal(content, &response)
	if err != nil {
		var respnoslice Click_resp_no_slice
		err := json.Unmarshal(content, &respnoslice)
		if err != nil {
			fmt.Println(string(content))
			ErrorLogger.Println("Error on unmarshal click response, err =", err)
			return &Click_resp{Ok: false}
		}
		return &Click_resp{Ok: respnoslice.Ok, Data: []Click_respdata{respnoslice.Data}}
	}
	return &response
}

func (Notcoin *Notcoin) click() {
	var data string
	urlstr := "https://clicker-api.joincommunity.xyz/clicker/core/click"
	count := Notcoin.getCountClick()
	webAppData := Notcoin.TGWebAppData
	if Notcoin.Turbo {
		if Notcoin.Hash != -1 {
			data = fmt.Sprintf(`{"count":%d, "hash":%d, "Turbo": true, "webAppData":"%v"}`, count, Notcoin.Hash, webAppData)
		} else {
			data = fmt.Sprintf(`{"count":%d, "Turbo": true, "webAppData":"%v"}`, count, webAppData)
		}
	} else {
		if Notcoin.Hash != -1 {
			data = fmt.Sprintf(`{"count":%d, "hash":%d, "webAppData":"%v"}`, count, Notcoin.Hash, webAppData)
		} else {
			data = fmt.Sprintf(`{"count":%d,"webAppData":"%v"}`, count, webAppData)
		}
	}

	resp := Notcoin.Ses.Postreq(urlstr, data)
	parsed_resp := parseRespclick(resp.body)

	if parsed_resp.Ok {
		Notcoin.Count_400 = 0
		Notcoin.LimitCoins = parsed_resp.Data[0].LimitCoins
		Notcoin.Hash = hashResolve(parsed_resp.Data[0].Hash)
		Notcoin.BalanceCoins = parsed_resp.Data[0].BalanceCoins
		//Notcoin.Coefficient = parsed_resp.Data[0].MultipleClicks
		Notcoin.Turbo_boost_count = parsed_resp.Data[0].TurboTimes
		Notcoin.LastAvailableCoins = parsed_resp.Data[0].LastAvailableCoins

	} else {
		if resp.status == 400 {
			Notcoin.Count_400++
		}
		Notcoin.Hash = -1
	}

	//conv to int
	// limCoinsINT, _ := strconv.Atoi(Notcoin.LimitCoins)
	if Notcoin.LastAvailableCoins < Notcoin.LimitCoins/2 &&
		Notcoin.Turbo_boost_count > 0 &&
		Notcoin.Count_400 == 0 &&
		Notcoin.Hash != -1 &&
		!Notcoin.Turbo {
		Notcoin.TurboActivate()
	}

	//fmt.Println(string(resp.body))

	if parsed_resp.Ok && resp.status < 400 {
		SuccessLogger.Printf("[%s] clicked and get %d coins, status = %d, balance = %s\n", Notcoin.UserId, count, resp.status, Notcoin.BalanceCoins)
	} else {
		WarningLogger.Printf("[%s] not success clicked %d times, status = %d\n", Notcoin.UserId, count, resp.status)
	}
}

func (not *Notcoin) TurboActivate() {
	var urlActivateTurbo = "https://clicker-api.joincommunity.xyz/clicker/core/active-turbo"
	var parsedResp Active_turbo_resp
	var ok bool

	resp := not.Ses.Postreq(urlActivateTurbo, "{}")

	err := json.Unmarshal(resp.body, &parsedResp)
	if err != nil {
		ErrorLogger.Println("Error on unmarshal turbo activate, err =", err)
		return
	}
	ok = parsedResp.Ok
	if !ok {
		ErrorLogger.Println("Error on turbo activate, response =", resp.String())
		return
	}

	not.Turbo = true
	not.Turbo_boost_count--
	not.Timestart_turbo = time.Now().Unix()
	SuccessLogger.Printf("[%d] Activated turbo", not.UserId)
}
