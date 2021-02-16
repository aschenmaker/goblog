package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"irisweb/config"
	"path/filepath"
	"strings"
	"time"
)

type Product struct {
	Model
	Id          uint           `json:"id" gorm:"column:id;type:int(10) unsigned not null AUTO_INCREMENT;primary_key"`
	Title       string         `json:"title" gorm:"column:title;type:varchar(250) not null;default:''"`
	UrlToken    string         `json:"url_token" gorm:"column:url_token;type:varchar(250) not null;default:'';index"`
	Keywords    string         `json:"keywords" gorm:"column:keywords;type:varchar(250) not null;default:''"`
	Description string         `json:"description" gorm:"column:description;type:varchar(250) not null;default:''"`
	CategoryId  uint           `json:"category_id" gorm:"column:category_id;type:int(10) unsigned not null;default:0;index:idx_category_id"`
	Price       float64        `json:"price" gorm:"column:price;type:decimal(10,2) unsigned not null;default:0;index:idx_price"`
	Stock       uint           `json:"stock" gorm:"column:stock;type:int(10) unsigned not null;default:0;index:idx_stock"`
	Views       uint           `json:"views" gorm:"column:views;type:int(10) unsigned not null;default:0;index:idx_views"`
	Images      pq.StringArray `json:"images" gorm:"column:images;type:text default null"`
	Status      uint           `json:"status" gorm:"column:status;type:tinyint(1) unsigned not null;default:0;index:idx_status"`
	CreatedTime int64          `json:"created_time" gorm:"column:created_time;type:int(11) not null;default:0;index:idx_created_time"`
	UpdatedTime int64          `json:"updated_time" gorm:"column:updated_time;type:int(11) not null;default:0;index:idx_updated_time"`
	DeletedTime int64          `json:"-" gorm:"column:deleted_time;type:int(11) not null;default:0"`
	Category    *Category      `json:"category" gorm:"-"`
	ProductData *ProductData   `json:"data" gorm:"-"`
	Logo        string         `json:"logo" gorm:"-"`
	Thumb       string         `json:"thumb" gorm:"-"`
}

type ProductData struct {
	Model
	Id      uint   `json:"id" gorm:"column:id;type:int(10) unsigned not null;primary_key"`
	Content string `json:"content" gorm:"column:content;type:longtext default null"`
}

func (product *Product) AddViews(db *gorm.DB) error {
	product.Views = product.Views + 1
	db.Model(Product{}).Where("`id` = ?", product.Id).Update("views", product.Views)
	return nil
}

func (product *Product) Save(db *gorm.DB) error {
	if product.Id == 0 {
		product.CreatedTime = time.Now().Unix()
	}
	product.UpdatedTime = time.Now().Unix()

	if err := db.Save(product).Error; err != nil {
		return err
	}
	if product.ProductData != nil {
		product.ProductData.Id = product.Id
		if err := db.Save(product.ProductData).Error; err != nil {
			return err
		}
	}

	return nil
}

func (product *Product) Delete(db *gorm.DB) error {
	if err := db.Model(product).Updates(Product{Status: 99, DeletedTime: time.Now().Unix()}).Error; err != nil {
		return err
	}

	return nil
}

func (product *Product) GetThumb() string {
	//取第一张
	if len(product.Images) > 0 {
		product.Logo = product.Images[0]
		//如果是一个远程地址，则缩略图和原图地址一致
		if strings.HasPrefix(product.Logo, "http") {
			product.Thumb = product.Logo
		} else {
			paths, fileName := filepath.Split(product.Logo)
			product.Thumb = config.JsonData.System.BaseUrl + paths + "thumb_" + fileName
		}
	} else if config.JsonData.Content.DefaultThumb != "" {
		product.Logo = config.JsonData.System.BaseUrl + config.JsonData.Content.DefaultThumb
		product.Thumb = product.Logo
	}

	return product.Thumb
}
