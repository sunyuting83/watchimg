package database

import (
	"github.com/jinzhu/gorm"
)

// Img List
type ImgList struct {
	ID         int64  `json:"id" gorm:"primary_key, column:id"`
	Account    string `json:"account" gorm:"varchar(128);index:idx_account_id;column:account"`
	Cover      string `json:"cover" gorm:"column:cover"`
	Password   string `json:"password" gorm:"column:password"`
	Today      int64  `json:"today" gorm:"column:today"`
	NewStatus  int64  `json:"new_status" gorm:"column:new_status"`
	UserID     int64  `json:"user_id" gorm:"column:user_id"`
	Multiple   int64  `json:"multiple" gorm:"column:multiple"`
	YesterDay  int64  `json:"yesterday" gorm:"column:yesterday"`
	UpDateTime int64  `json:"updatetime" gorm:"column:updatetime"`
	DateTime   int64  `json:"datetime" gorm:"column:datetime"`
	ExpTime    string `json:"expdate" gorm:"column:expdate"`
}
type ImgLists struct {
	ID         int64  `json:"id" gorm:"primary_key, column:id"`
	Account    string `json:"account" gorm:"varchar(128);index:idx_account_id;column:account"`
	Cover      string `json:"cover" gorm:"column:cover"`
	Password   string `json:"password" gorm:"column:password"`
	Today      int64  `json:"today" gorm:"column:today"`
	NewStatus  int64  `json:"new_status" gorm:"column:new_status"`
	UserName   string `json:"username" gorm:"column:username"`
	Multiple   int64  `json:"multiple" gorm:"column:multiple"`
	YesterDay  int64  `json:"yesterday" gorm:"column:yesterday"`
	UpDateTime int64  `json:"updatetime" gorm:"column:updatetime"`
	DateTime   int64  `json:"datetime" gorm:"column:datetime"`
	ExpTime    string `json:"expdate" gorm:"column:expdate"`
}

type DateTimeData struct {
	DateList  int64 `json:"date" gorm:"type:datetime"`
	NewStatus int64 `json:"new_status" gorm:"column:new_status"`
}

// TableName change table name
func (ImgList) TableName() string {
	return "imageData"
}

// GetImgOne 列表
func (imglist *ImgList) GetImgOne(account string) (imglists ImgList, err error) {
	if err = Eloquent.First(&imglists, "account = ?", account).Error; err != nil {
		return
	}
	return
}

// GetCount 列表
func (imglist *ImgList) GetCount() (count int64, err error) {
	if err = Eloquent.Model(&imglist).Where("new_status = ?", 0).Count(&count).Error; err != nil {
		return
	}
	return
}

// GetImgList 列表
func (imglist *ImgList) GetImgList(page int64) (imglists []ImgList, err error) {
	p := makePage(page)
	if err = Eloquent.
		Select("id, account, cover, today, yesterday, multiple, datetime").
		Where("new_status = ?", 0).
		Order("today desc").
		Limit(100).Offset(p).
		Find(&imglists).Error; err != nil {
		return
	}
	return
}

// GetImgList 列表
func (imglist *ImgList) GetOldCount() (count int64, err error) {
	// sql := `SELECT s.*,t.username FROM imageData AS s INNER JOIN user AS t WHERE s.new_status = 1 AND s.user_id = t.id`
	if err = Eloquent.Model(&imglist).Where("new_status = ?", 1).Count(&count).Error; err != nil {
		return
	}
	return
}

// GetImgList 列表
func (imglist *ImgList) GetOldImgList() (imglists []ImgLists, err error) {
	sql := `SELECT s.*,t.username FROM imageData AS s INNER JOIN user AS t WHERE s.new_status = 1 AND s.user_id = t.id`
	if err = Eloquent.
		Raw(sql).
		Order("updatetime desc").
		Scan(&imglists).Error; err != nil {
		return
	}
	return
}

// Insert Data
func (datalist *ImgList) Insert() error {
	Eloquent.Create(&datalist)
	return nil
}

// UpdateStatus data
func (datalist *ImgList) UpdateStatus(account string) (update ImgList, err error) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	if err = Eloquent.First(&update, "account = ?", account).Error; err != nil {
		return
	}
	if err = Eloquent.Model(&update).Updates(&datalist).Error; err != nil {
		return
	}
	return
}

// UpdateStatus data
func (datalist *ImgList) UpdateImg(account string) (update ImgList, err error) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	if err = Eloquent.First(&update, "account = ?", account).Error; err != nil {
		return
	}
	if err = Eloquent.Model(&update).Updates(&datalist).Error; err != nil {
		return
	}
	return
}

// DeleteOne data
func (datalist *ImgList) DeleteOne(account string) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	Eloquent.Where("account = ?", account).Delete(&datalist)
}

// Get DateTime Data
func GetDateTimeData(status string) (re []string, err error) {
	var d string = "updatetime"
	if status == "0" {
		d = "datetime"
	}
	sql := "SELECT DISTINCT DATE(" + d + ",'unixepoch','localtime'),new_status FROM imageData WHERE new_status = " + status + " ORDER BY " + d + " DESC"
	re, err = RawQuerySearchAndParseToMap(Eloquent, sql)
	return
}

// Get DateTime Data
func GetDateTimeDataNList(date string) (list []*ImgList, err error) {
	sql := "SELECT * from imageData where DATE(datetime,'unixepoch','localtime') = '" + date + "' AND new_status = 0"
	// fmt.Println(sql)
	if err = Eloquent.
		Raw(sql).
		Order("today desc").
		Scan(&list).Error; err != nil {
		return
	}
	return
}
func GetDateTimeDataYList(date, st string) (list []*ImgLists, err error) {
	sql := `SELECT s.*,t.username FROM imageData AS s INNER JOIN user AS t WHERE s.new_status = ` + st + ` AND s.user_id = t.id AND DATE(s.updatetime,'unixepoch','localtime') = '` + date + "'"
	// sql := "SELECT * from imageData where DATE(updatetime,'unixepoch') = '" + date + "' AND new_status = 1"
	// fmt.Println(sql)
	if err = Eloquent.
		Raw(sql).
		Order("today desc").
		Scan(&list).Error; err != nil {
		return
	}
	return
}

// SearchKey 列表
func (search *ImgList) SearchKey(key string) (searchs []*ImgList, err error) {
	if err = Eloquent.
		Where("account LIKE ? AND new_status = 0", "%"+key+"%").
		Order("id desc").
		Find(&searchs).Error; err != nil {
		return
	}
	return
}

// SearchKey 列表
func (imglist *ImgList) SearchKeyY(key, st string) (searchs []*ImgLists, err error) {
	sql := `SELECT s.*,t.username FROM imageData AS s INNER JOIN user AS t WHERE s.new_status = ` + st + ` AND s.user_id = t.id AND s.account LIKE '%` + key + "%'"
	if err = Eloquent.
		Raw(sql).
		Order("s.id").
		Scan(&searchs).Error; err != nil {
		return
	}
	return
}

// makePage make page
func makePage(p int64) int64 {
	p = p - 1
	if p <= 0 {
		p = 0
	}
	page := p * 100
	return page
}

// RawQuerySearchAndParseToMap ...
func RawQuerySearchAndParseToMap(db *gorm.DB, query string) ([]string, error) {
	//Use raw query
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	//取得搜尋回來的資料所擁有的column
	columns, er := rows.Columns()
	if er != nil {
		return nil, er
	}
	columnLength := len(columns)

	//make一個臨時儲存的地方，並賦予指標
	cache := make([]interface{}, columnLength)
	for index := range cache {
		var a interface{}
		cache[index] = &a
	}

	var list []map[string]interface{}
	for rows.Next() {
		rows.Scan(cache...)

		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //column可能有許多種data type，因此在這取出時不指定型別，否則會轉換錯誤，且在這取出時為uint8(btye array)格式
		}

		list = append(list, item)
	}
	var l []string
	//將byte array轉換為字串
	for index := range list {
		for _, column := range columns {
			if column == "DATE(datetime,'unixepoch','localtime')" || column == "DATE(updatetime,'unixepoch','localtime')" {
				list[index][column] = list[index][column].(string)
				l = append(l, list[index][column].(string))
			}
		}
	}

	return l, nil

}
