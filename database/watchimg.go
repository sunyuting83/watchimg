package database

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

// makePage make page
func makePage(p int64) int64 {
	p = p - 1
	if p <= 0 {
		p = 0
	}
	page := p * 100
	return page
}
