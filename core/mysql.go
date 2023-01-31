package core

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"time"

	"strconv"
	"strings"

	// mysql SQL driver
	_ "github.com/go-sql-driver/mysql"
)

// Mysql structure
type Mysql struct {
	Enabled    bool   `json:"enabled"`
	ServerAddr string `json:"server_addr"`
	ServerPort int    `json:"server_port"`
	Database   string `json:"database"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Cafile     string `json:"cafile"`
}

// User: User table record structure
type User struct {
	ID          uint
	Username    string
	Password    string
	EncryptPass string
	Quota       int64
	Download    uint64
	Upload      uint64
	UseDays     uint
	ExpiryDate  string
}

// PageQuery Paging structure
type PageQuery struct {
	PageNum  int
	CurPage  int
	Total    int
	PageSize int
	DataList []*User
}

// CreateTablesql Create database table SQL
var CreateTableSql = `
CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL,
    password CHAR(56) NOT NULL,
    passwordShow VARCHAR(255) NOT NULL,
    quota BIGINT NOT NULL DEFAULT 0,
    download BIGINT UNSIGNED NOT NULL DEFAULT 0,
    upload BIGINT UNSIGNED NOT NULL DEFAULT 0,
    useDays int(10) DEFAULT 0,
    expiryDate char(10) DEFAULT '',
    PRIMARY KEY (id),
    INDEX (password)
) DEFAULT CHARSET=utf8mb4;
`

// GetDB Get mysql database connection
func (mysql *Mysql) GetDB() *sql.DB {
	// Mask the log output of the MySQL drive package
	mysqlDriver.SetLogger(log.New(ioutil.Discard, "", 0))
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysql.Username, mysql.Password, mysql.ServerAddr, mysql.ServerPort, mysql.Database)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return db
}

// CreateTable Create trojan user table if it is not exists
func (mysql *Mysql) CreateTable() {
	db := mysql.GetDB()
	defer db.Close()
	if _, err := db.Exec(CreateTableSql); err != nil {
		fmt.Println(err)
	}
}

func queryUserList(db *sql.DB, sql string) ([]*User, error) {
	var (
		username    string
		encryptPass string
		passShow    string
		download    uint64
		upload      uint64
		quota       int64
		id          uint
		useDays     uint
		expiryDate  string
	)
	var userList []*User
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &username, &encryptPass, &passShow, &quota, &download, &upload, &useDays, &expiryDate); err != nil {
			return nil, err
		}
		userList = append(userList, &User{
			ID:          id,
			Username:    username,
			Password:    passShow,
			EncryptPass: encryptPass,
			Download:    download,
			Upload:      upload,
			Quota:       quota,
			UseDays:     useDays,
			ExpiryDate:  expiryDate,
		})
	}
	return userList, nil
}

func queryUser(db *sql.DB, sql string) (*User, error) {
	var (
		username    string
		encryptPass string
		passShow    string
		download    uint64
		upload      uint64
		quota       int64
		id          uint
		useDays     uint
		expiryDate  string
	)
	row := db.QueryRow(sql)
	if err := row.Scan(&id, &username, &encryptPass, &passShow, &quota, &download, &upload, &useDays, &expiryDate); err != nil {
		return nil, err
	}
	return &User{ID: id, Username: username, Password: passShow, EncryptPass: encryptPass, Download: download, Upload: upload, Quota: quota, UseDays: useDays, ExpiryDate: expiryDate}, nil
}

// CreateUser Create a trojan user
func (mysql *Mysql) CreateUser(username string, base64Pass string, originPass string) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("Can't connect to mysql!")
	}
	defer db.Close()
	encryPass := sha256.Sum224([]byte(originPass))
	if _, err := db.Exec(fmt.Sprintf("INSERT INTO users(username, password, passwordShow, quota) VALUES ('%s', '%x', '%s', -1);", username, encryPass, base64Pass)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// UpdateUser Update trojan user name and password
func (mysql *Mysql) UpdateUser(id uint, username string, base64Pass string, originPass string) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("Can't connect to mysql!")
	}
	defer db.Close()
	encryPass := sha256.Sum224([]byte(originPass))
	if _, err := db.Exec(fmt.Sprintf("UPDATE users SET username='%s', password='%x', passwordShow='%s' WHERE id=%d;", username, encryPass, base64Pass, id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// DeleteUser delete users
func (mysql *Mysql) DeleteUser(id uint) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("Can't connect to mysql!")
	}
	defer db.Close()
	if userList, err := mysql.GetData(strconv.Itoa(int(id))); err != nil {
		return err
	} else if userList != nil && len(userList) == 0 {
		return fmt.Errorf("There are no users with id %d", id)
	}
	if _, err := db.Exec(fmt.Sprintf("DELETE FROM users WHERE id=%d;", id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// MonthlyResetData The user who has set an expired time shall clear the traffic regularly every month.
func (mysql *Mysql) MonthlyResetData() error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("Can't connect to mysql!")
	}
	defer db.Close()
	userList, err := queryUserList(db, "SELECT * FROM users WHERE useDays != 0 AND quota != 0")
	if err != nil {
		return err
	}
	for _, user := range userList {
		if _, err := db.Exec(fmt.Sprintf("UPDATE users SET download=0, upload=0 WHERE id=%d;", user.ID)); err != nil {
			return err
		}
	}
	return nil
}

// DailyCheckExpire Check whether there is an expired, the upper limit of the setting flow is 0
func (mysql *Mysql) DailyCheckExpire() (bool, error) {
	needRestart := false
	now := time.Now()
	utc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return false, err
	}
	addDay, _ := time.ParseDuration("-24h")
	yesterdayStr := now.Add(addDay).In(utc).Format("2006-01-02")
	yesterday, _ := time.Parse("2006-01-02", yesterdayStr)
	db := mysql.GetDB()
	if db == nil {
		return false, errors.New("can't connect to mysql")
	}
	defer db.Close()
	userList, err := queryUserList(db, "SELECT * FROM users WHERE quota != 0")
	if err != nil {
		return false, err
	}
	for _, user := range userList {
		if expireDate, err := time.Parse("2006-01-02", user.ExpiryDate); err == nil {
			if yesterday.Sub(expireDate).Seconds() >= 0 {
				if _, err := db.Exec(fmt.Sprintf("UPDATE users SET quota=0 WHERE id=%d;", user.ID)); err != nil {
					return false, err
				}
				if !needRestart {
					needRestart = true
				}
			}
		}
	}
	return needRestart, nil
}

// CancelExpire Cancel the expiration time
func (mysql *Mysql) CancelExpire(id uint) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("Can't connect to mysql!")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE users SET useDays=0, expiryDate='' WHERE id=%d;", id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// SetExpire Set the expiration time
func (mysql *Mysql) SetExpire(id uint, useDays uint) error {
	now := time.Now()
	utc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return err
	}
	addDay, _ := time.ParseDuration(strconv.Itoa(int(24*useDays)) + "h")
	expiryDate := now.Add(addDay).In(utc).Format("2006-01-02")

	db := mysql.GetDB()
	if db == nil {
		return errors.New("Can't connect to mysql!")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE users SET useDays=%d, expiryDate='%s' WHERE id=%d;", useDays, expiryDate, id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// SetQuota Restricting traffic
func (mysql *Mysql) SetQuota(id uint, quota int) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("Can't connect to mysql!")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE users SET quota=%d WHERE id=%d;", quota, id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// CleanData Clear traffic statistics
func (mysql *Mysql) CleanData(id uint) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("Can't connect to mysql!")
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("UPDATE users SET download=0, upload=0 WHERE id=%d;", id)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// CleanDataByName Clear statistics of designated user presence volume
func (mysql *Mysql) CleanDataByName(usernames []string) error {
	db := mysql.GetDB()
	if db == nil {
		return errors.New("Can't connect to mysql!")
	}
	defer db.Close()
	runSql := "UPDATE users SET download=0, upload=0 WHERE BINARY username in ("
	for i, name := range usernames {
		runSql = runSql + "'" + name + "'"
		if i == len(usernames)-1 {
			runSql = runSql + ")"
		} else {
			runSql = runSql + ","
		}
	}
	if _, err := db.Exec(runSql); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// GetUserByName Find user with username
func (mysql *Mysql) GetUserByName(name string) *User {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	defer db.Close()
	user, err := queryUser(db, fmt.Sprintf("SELECT * FROM users WHERE BINARY username='%s'", name))
	if err != nil {
		return nil
	}
	return user
}

// GetUserByPass Find users by password
func (mysql *Mysql) GetUserByPass(pass string) *User {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	defer db.Close()
	user, err := queryUser(db, fmt.Sprintf("SELECT * FROM users WHERE BINARY passwordShow='%s'", pass))
	if err != nil {
		return nil
	}
	return user
}

// PageList Obtain user records by pages
func (mysql *Mysql) PageList(curPage int, pageSize int) (*PageQuery, error) {
	var (
		total int
	)

	db := mysql.GetDB()
	if db == nil {
		return nil, errors.New("Can't connect to mysql!")
	}
	defer db.Close()
	offset := (curPage - 1) * pageSize
	querySQL := fmt.Sprintf("SELECT * FROM users LIMIT %d, %d", offset, pageSize)
	userList, err := queryUserList(db, querySQL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	db.QueryRow("SELECT COUNT(id) FROM users").Scan(&total)
	return &PageQuery{
		CurPage:  curPage,
		PageSize: pageSize,
		Total:    total,
		DataList: userList,
		PageNum:  (total + pageSize - 1) / pageSize,
	}, nil
}

// GetData Get user record
func (mysql *Mysql) GetData(ids ...string) ([]*User, error) {
	querySQL := "SELECT * FROM users"
	db := mysql.GetDB()
	if db == nil {
		return nil, errors.New("Can't connect to mysql!")
	}
	defer db.Close()
	if len(ids) > 0 {
		querySQL = querySQL + " WHERE id in (" + strings.Join(ids, ",") + ")"
	}
	userList, err := queryUserList(db, querySQL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return userList, nil
}
