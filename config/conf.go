package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ditto/booking/logger"
	"ditto/booking/utils"

	"github.com/kardianos/osext"
)

//
var (
	configfile = "booking.json"
	setting    *Config
)

//AccountConfig -
type AccountConfig struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

//RsaConfig -
type RsaConfig struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

//DbConfig for sqlite3
type DbConfig struct {
	Type            string `json:"type"`
	DNSfmt          string `json:"dns"`
	DNS             string `json:"-"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	User            string `json:"user"`
	Password        string `json:"password"`
	Database        string `json:"database"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	MaxOpenConns    int    `json:"max_open_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
	Debug           bool   `json:"debug"`
}

//Log -
type Log struct {
	Level  string `json:"level"`  //ログレベル
	Dir    string `json:"dir"`    //ログファイルの出力先
	Output string `json:"output"` //出力先[screen, file]
	Format string `json:"format"` //ログファイル名フォーマット
}

//CacheConfig -
type CacheConfig struct {
	Server   string `json:"server"`   //Server
	Address  string `json:"address"`  //address
	Password string `json:"password"` //password
	Enable   bool   `json:"enable"`   //Enable
}

//HolidaysConfig -
type HolidaysConfig struct {
	URL    string `json:"url"`
	Encode string `json:"encode"`
	Header bool   `json:"header"`
}

//Config -
type Config struct {
	Error          error            `json:"-"`               //Error
	Version        string           `json:"-"`               //Version
	ExeName        string           `json:"-"`               //モジュールパス
	FileName       string           `json:"-"`               //s3service.jsonパス
	Dir            string           `json:"-"`               //work directory
	LogFile        string           `json:"-"`               //logfile fullpath
	Host           string           `json:"host"`            //サーバーIP
	Port           int              `json:"port"`            //ポート
	BaseURL        string           `json:"baseurl"`         //BaseURL
	Timeout        int              `json:"timeout"`         //タイムアウト（秒）
	Log            Log              `json:"log"`             //ログ
	Expires        int64            `json:"expires"`         //token expires in hour (access token)（時間）
	ExpiresConfirm int64            `json:"expires_confirm"` //アカウント作成時確認コードの有効期限（時間）
	Rsa            RsaConfig        `json:"rsa"`             //Rsa key
	Db             DbConfig         `json:"db"`              //DB設定
	Account        []*AccountConfig `json:"account"`         //account
	Cache          CacheConfig      `json:"cache"`           //Cache server (redis)
	Holidays       HolidaysConfig   `json:"holidays"`        //Holidays
}

//Load -
func Load() *Config {
	if setting != nil {
		return setting
	}

	//initialize
	setting = &Config{
		Host: "localhost",
		Port: 9898,
	}

	exename, _ := osext.Executable()
	setting.ExeName = exename

	dir := filepath.Dir(exename)
	name := filepath.Join(dir, configfile)
	setting.FileName = name
	setting.Dir = dir
	fmt.Println(name)

	fr, err := os.Open(name)
	if err != nil {
		setting.Error = err
		return setting
	}
	defer fr.Close()

	//s3service -json
	err = utils.JSON.NewDecoder(fr).Decode(&setting)
	if err != nil {
		setting.Error = err
		return setting
	}
	if setting.Dir == "" {
		setting.Dir = dir
	}

	//initializer
	setting.Init()
	setting.Error = nil

	return setting
}

//Init -
func (c *Config) Init() {
	if c.Log.Output == "" {
		c.Log.Output = "file"
	}
	var logfile string
	if c.Log.Output == "file" || c.Log.Output == "both" {
		err := utils.MakeDirectory(c.Log.Dir)
		if err != nil {
			c.Log.Dir = ""
		}

		logfile = getLogFilename(c.Log.Format)
		logfile = filepath.Join(c.Log.Dir, logfile)

		c.LogFile = logfile
	}
	logger.SetOutput(c.Log.Output, logfile)

	//Replace
	port := fmt.Sprintf("%+v", c.Db.Port)
	dns := c.Db.DNSfmt
	dns = strings.Replace(dns, "${user}", c.Db.User, -1)
	dns = strings.Replace(dns, "${password}", c.Db.Password, -1)
	dns = strings.Replace(dns, "${host}", c.Db.Host, -1)
	dns = strings.Replace(dns, "${port}", port, -1)
	dns = strings.Replace(dns, "${database}", c.Db.Database, -1)

	c.Db.DNS = dns

	//Default
	//Access tokenの有効期間：24時間
	if c.Expires == 0 {
		c.Expires = 24
	}
	//アカウント作成時確認コードの有効期限：48時間
	if c.ExpiresConfirm == 0 {
		c.ExpiresConfirm = 48
	} else if c.ExpiresConfirm < 0 {
		c.ExpiresConfirm = 0
	}
}

//getLogFilename -
func getLogFilename(format string) string {
	logfile := "booking.log" //20060102150405

	pos := strings.Index(format, "%")
	if pos > -1 {
		lpos := strings.LastIndex(format, "%")
		if lpos > -1 {
			if lpos > pos {
				fstr := format[pos+1 : lpos]
				if len(fstr) > 0 {
					ffmt := format[:pos] + "%v" + format[lpos+1:]
					//20060102150405 "%yyyymmddHHMMSS%"
					fstr = strings.ReplaceAll(fstr, "yyyy", "2006")
					fstr = strings.ReplaceAll(fstr, "mm", "01")
					fstr = strings.ReplaceAll(fstr, "dd", "02")
					fstr = strings.ReplaceAll(fstr, "HH", "15")
					fstr = strings.ReplaceAll(fstr, "MM", "04")
					fstr = strings.ReplaceAll(fstr, "SS", "05")

					logfile = fmt.Sprintf(ffmt, utils.NowJST().Format(fstr)) //20060102150405
				}
			}
		}
	} else {
		if format != "" {
			logfile = format
		}
	}

	return logfile
}

func dirWindows() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// Prefer standard environment variable USERPROFILE
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home, nil
	}

	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, or USERPROFILE are blank")
	}

	return home, nil
}
