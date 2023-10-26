package store

//go:generate mockgen -destination=./code_store_mock.go -package=store go-chromedp/app/store ICodeStore
import (
	"context"
	"encoding/json"
	"go-chromedp/app/models"
	"gorm.io/gorm"
	"time"
)

const (
	CodePlatformAll CodePlatform = iota + 1
	CodePlatformIOS
	CodePlatformMAC
)
const (
	CodeStatusNormal CodeStatus = iota + 1
	CodeStatusRecycle
)
const (
	CodeTypeCategory CodeType = iota + 1
	CodeTypeCertificate
)
const (
	CodeStrategyAverage CodeStrategy = iota + 1
	CodeStrategyFull
	CodeStrategyAverageLimit
	CodeStrategyFullLimit
)

type (
	CodePlatform int16
	CodeStatus   int16
	CodeType     int16
	CodeStrategy int16

	ICodeStore interface {
		Create(ctx *context.Context, entity *CodeEntity) error
		Update(ctx *context.Context, entity *CodeEntity, fields []string) error
		Updates(ctx *context.Context, entity *CodeEntity) error
		Delete(ctx *context.Context, id string) error
		All(ctx *context.Context) ([]CodeEntity, error)
		List(ctx *context.Context, offset int, limit int) ([]CodeEntity, error)
		Get(ctx *context.Context, id string) (*CodeEntity, error)
		Total(ctx *context.Context) (int, error)
		FindByUDID(ctx *context.Context, udid string) ([]CodeEntity, error)
		InitSchema()
	}
	CodeStore struct {
		ICodeStore
	}
	CodeEntity struct {
		ID            string       `gorm:"primaryKey"` // 主键，由业务方生成
		Udid          string       `gorm:"index"`      // 绑定的设备
		Status        CodeStatus   `gorm:"index"`      // 状态，1:正常，2:回收
		Type          CodeType     `gorm:"index"`      // 类型，1:绑定分类，2:绑定证书
		TypeValue     int          // 绑定类型的值，如分类ID，证书ID
		Platform      CodePlatform `gorm:"index"` // 分配平台，ALL/IOS/MAC_OS
		Strategy      CodeStrategy `gorm:"index"` // 分配策略, 1:平均，2:满载，3:平均+限制，4:满载+限制
		StrategyValue int          // 自定义策略限制值，如：10
		Expire        time.Time    `gorm:"index"` // 过期时间
		UpdatedAt     time.Time
		CreatedAt     time.Time
		DeletedAt     gorm.DeletedAt
	}
)

func (c CodeEntity) String() string {
	str, _ := json.Marshal(c)
	return string(str)
}

func NewCodeStore() ICodeStore {
	return &CodeStore{}
}

func (*CodeStore) Total(ctx *context.Context) (int, error) {
	var count int64
	t := TraceDB(ctx).
		Model(&CodeEntity{}).
		Count(&count)
	return int(count), t.Error
}
func (*CodeStore) All(ctx *context.Context) (lists []CodeEntity, err error) {
	t := TraceDB(ctx).
		Find(&lists)
	err = t.Error
	return
}
func (*CodeStore) List(ctx *context.Context, offset int, limit int) (lists []CodeEntity, err error) {
	t := TraceDB(ctx).
		Offset(offset).
		Limit(limit).
		Find(&lists)
	err = t.Error
	return
}
func (*CodeStore) Get(ctx *context.Context, id string) (entity *CodeEntity, err error) {
	t := TraceDB(ctx).
		Where("id = ?", id).
		First(entity)
	err = t.Error
	return
}
func (*CodeStore) Create(ctx *context.Context, entity *CodeEntity) error {
	t := TraceDB(ctx).
		Create(entity)
	return t.Error
}
func (*CodeStore) Update(ctx *context.Context, entity *CodeEntity, fields []string) error {
	t := TraceDB(ctx).
		Model(entity).
		Select(fields).
		Updates(*entity)
	return t.Error
}
func (*CodeStore) FindByUDID(ctx *context.Context, udid string) (lists []CodeEntity, err error) {
	t := TraceDB(ctx).
		Where(&CodeEntity{
			Udid: udid,
		}).
		Find(&lists)
	err = t.Error
	return
}
func (*CodeStore) Updates(ctx *context.Context, entity *CodeEntity) error {
	t := TraceDB(ctx).
		Model(entity).
		Select("*").
		Updates(*entity)
	return t.Error
}
func (*CodeStore) Delete(ctx *context.Context, id string) (err error) {
	t := TraceDB(ctx).
		Delete(&CodeEntity{
			ID: id,
		})
	err = t.Error
	return
}
func (*CodeStore) InitSchema() {
	if err := models.GetDbInstance().AutoMigrate(&CodeEntity{}); err != nil {
		panic(err)
	}
}
