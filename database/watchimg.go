package database

// Img List
type ImgList struct {
	ID        int64  `json:"id" gorm:"primary_key, column:id"`
	Account   string `json:"account" gorm:"varchar(128);index:idx_account_id;column:account"`
	Cover     string `json:"cover" gorm:"column:cover"`
	Today     int64  `json:"today" gorm:"column:today"`
	YesterDay int64  `json:"yesterday" gorm:"column:yesterday"`
	DateTime  int64  `json:"datetime" gorm:"column:datetime"`
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
	if err = Eloquent.Model(&imglist).Count(&count).Error; err != nil {
		return
	}
	return
}

// GetImgList 列表
func (imglist *ImgList) GetImgList(page int64) (imglists []ImgList, err error) {
	p := makePage(page)
	if err = Eloquent.
		Select("id, account, cover, today, yesterday, datetime").
		Order("today desc").
		Limit(100).Offset(p).
		Find(&imglists).Error; err != nil {
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
