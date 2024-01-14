package clicker

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := os.Truncate("logs.txt", 0); err != nil {
		ErrorLogger.Fatalln(err)
	}
	if err := godotenv.Load(); err != nil {
		ErrorLogger.Fatalln(err)
	}
}

func (notcoin *Notcoin) Set_default_values() {
	notcoin.LimitCoins = getRandomint(148, 170, 1)
	notcoin.LastAvailableCoins = getRandomint(1, 147, 1)
	notcoin.Coefficient = 1
	notcoin.Hash = -1

}

func (notcoin *Notcoin) work(wg *sync.WaitGroup) {
	defer wg.Done()
	var sleep_time int
	var is_slept bool = true
	var userid int = notcoin.UserId
	InfoLogger.Printf("[%v] Starting thread id = %v, Name = %v, Proxy = %v", userid, userid, notcoin.Clear_name, notcoin.Proxy)
	notcoin.Ses = CreateSession()

	err := notcoin.SetSession()
	if err != nil {
		ErrorLogger.Printf("[%v] Err in set session, close thread...", userid)
		return
	}

	notcoin.Set_default_values()

	if os.Getenv("enable_daily_boosters") == "1" {
		notcoin.UpdateBoosters()
	}

	for count_iteration := 0; count_iteration != -1; count_iteration++ {
		if count_iteration%30 == 0 && count_iteration > 1 || count_iteration == 1 {
			notcoin.UpdateShop()

			if os.Getenv("enable_daily_boosters") == "1" {
				notcoin.UpdateBoosters()
			}
		}

		if notcoin.LastAvailableCoins >= 100 || is_slept || notcoin.Turbo {
			notcoin.click()
			is_slept = false
		} else {
			if notcoin.Fullenergy_boost >= 1 && notcoin.Hash != -1 {
				InfoLogger.Printf("[%v] ACTIVATED FULLENERGY\n", userid)
				ok := notcoin.ActiveTask(2)
				InfoLogger.Println("ok active=", ok)
				if ok {
					time.Sleep(time.Second * time.Duration(2))
					continue
				}
			}
			InfoLogger.Printf("[%v] AvailableCoins is done, wait 50-100 seconds....\n", userid)
			notcoin.LastAvailableCoins += getRandomint(1, 30, 1)
			sleep_time = getRandomint(51, 98, 1)
			time.Sleep(time.Second * time.Duration(sleep_time))
			is_slept = true
			continue
		}
		if notcoin.Count_400 >= 7 {
			notcoin.Ses = CreateSession()
			err := notcoin.SetSession()
			if err != nil {
				ErrorLogger.Printf("[%v] Err in set session, close thread...", userid)
				return
			}
			notcoin.Count_400 = 0
			continue
		}
		if notcoin.Turbo {
			if notcoin.Count_400 > 0 {
				time.Sleep(time.Second * time.Duration(5))
			} else {
				time.Sleep(time.Second * time.Duration(3))
			}
			continue
		}
		if notcoin.LastAvailableCoins < 150 {
			time.Sleep(time.Second * time.Duration(2))
			continue
		}
		if notcoin.Hash == -1 {
			InfoLogger.Printf("[%v] bad hash, wait %v seconds...\n", userid, 10)
			time.Sleep(time.Second * time.Duration(10))
			continue
		}

		sleep_time = getRandomint(2, 13, 1)
		InfoLogger.Printf("[%v] wait %v seconds...\n", userid, sleep_time)
		time.Sleep(time.Second * time.Duration(sleep_time))

	}

}

func (notcoin *Notcoin) setreqwebappdata(raw_data string) {
	split1 := strings.Split(raw_data, "#tgWebAppData=")[1]
	split2 := strings.Split(split1, "&tgWebAppVersion")[0]
	decodedUrl, _ := url.QueryUnescape(split2)
	notcoin.TGWebAppData = decodedUrl
}

func (notcoin *Notcoin) SetSession() error {
	var rwebses Webappses_resp

	if len(notcoin.Proxy) > 1 {
		err := notcoin.Ses.Set_proxy(notcoin.Proxy)
		if err != nil {
			ErrorLogger.Printf("[%v]Err on set proxy '%v' %v\n", notcoin.UserId, notcoin.Proxy, err)
		}
	}

	raw_tgwebappdata, err := notcoin.getAppdata()
	if err != nil {
		ErrorLogger.Printf("[%v] Anything err in get appdata, err = %v\n", notcoin.UserId, err)
		return err
	}
	notcoin.setreqwebappdata(raw_tgwebappdata)
	url := "https://clicker-api.joincommunity.xyz/auth/webapp-session"
	webAppData := notcoin.TGWebAppData
	data := fmt.Sprintf(`{"webAppData":"%v"}`, webAppData)

	resp := notcoin.Ses.Postreq(url, data)
	_ = json.Unmarshal(resp.body, &rwebses)
	notcoin.Ses.headers.Add("Authorization", fmt.Sprintf("Bearer %v", rwebses.Data.AccessToken))
	return nil
}

func get_files(path string) []fs.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		ErrorLogger.Fatalf("err on get files in %v err = %v\n", path, err)
	}
	return files
}

func get_accounts_data() Accounts_data {
	var data Accounts_data
	file, err := os.ReadFile("accounts.json")
	if err != nil {
		ErrorLogger.Fatalf("err on read accounts.json err = %v\n", err)
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		ErrorLogger.Fatalf("err on decode json accounts.json err = %v\n", err)
	}

	return data
}

func ClickerStart() {
	var wg sync.WaitGroup

	path_sessions := "./sessions//"
	files := get_files(path_sessions)
	accounts := get_accounts_data()

	for i, file := range files {
		filename := file.Name()

		path := fmt.Sprintf("%v%v", path_sessions, filename)
		clear_name := strings.Split(filename, ".")[0]
		account, ok := accounts[clear_name]
		if !ok {
			ErrorLogger.Printf("Account %v Not found in accounts.json\n", clear_name)
			continue
		}

		work := Notcoin{
			Clear_name: clear_name,
			Path_file:  path,
			UserId:     i,
			TG_appHash: account.APIHash,
			TG_appID:   account.APIID,
			Proxy:      account.Proxy,
		}
		wg.Add(1)
		go work.work(&wg)

	}
	wg.Wait()
}
