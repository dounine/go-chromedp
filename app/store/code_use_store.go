package store

//go:generate mockgen -destination=./code_use_store_mock.go -package=store go-chromedp/app/store ICodeUseStore
import (
	"context"
	"encoding/json"
	"go-chromedp/app/models"
	"gorm.io/gorm"
	"time"
)

type (
	ICodeUseStore interface {
		Create(ctx *context.Context, entity *CodeUseEntity) error
		Update(ctx *context.Context, entity *CodeUseEntity, fields []string) error
		Updates(ctx *context.Context, entity *CodeUseEntity) error
		Delete(ctx *context.Context, id int) error
		All(ctx *context.Context) ([]CodeUseEntity, error)
		List(ctx *context.Context, offset int, limit int) ([]CodeUseEntity, error)
		QueryByUdidCodeID(ctx *context.Context, udid string, codeid string) ([]CodeUseEntity, error)
		Get(ctx *context.Context, id int) (*CodeUseEntity, error)
		Total(ctx *context.Context) (int, error)
		InitSchema()
	}
	CodeUseStore struct {
		ICodeUseStore
	}
	CodeUseEntity struct {
		ID        int
		Udid      string `gorm:"index:,composite:udid_codeid"` // 绑定的设备
		CodeId    string `gorm:"index:,composite:udid_codeid"` // 绑定的兑换码
		P8ID      int    // 绑定的p8
		DeviceId  string // 设备ID
		UpdatedAt time.Time
		CreatedAt time.Time
		DeletedAt gorm.DeletedAt
	}
)

func (c CodeUseEntity) String() string {
	str, _ := json.Marshal(c)
	return string(str)
}

func NewCodeUseStore() ICodeUseStore {
	return &CodeUseStore{}
}

func (*CodeUseStore) Total(ctx *context.Context) (int, error) {
	var count int64
	t := TraceDB(ctx).
		Model(&CodeUseEntity{}).
		Count(&count)
	return int(count), t.Error
}
func (*CodeUseStore) All(ctx *context.Context) (lists []CodeUseEntity, err error) {
	t := TraceDB(ctx).
		Find(&lists)
	err = t.Error
	return
}
func (*CodeUseStore) List(ctx *context.Context, offset int, limit int) (lists []CodeUseEntity, err error) {
	t := TraceDB(ctx).
		Offset(offset).
		Limit(limit).
		Find(&lists)
	err = t.Error
	return
}
func (*CodeUseStore) Get(ctx *context.Context, id int) (entity *CodeUseEntity, err error) {
	t := TraceDB(ctx).
		Where(&CodeUseEntity{
			ID: id,
		}).
		First(entity)
	err = t.Error
	return
}
func (*CodeUseStore) QueryByUdidCodeID(ctx *context.Context, udid string, codeid string) (lists []CodeUseEntity, err error) {
	t := TraceDB(ctx).
		Where(&CodeUseEntity{
			Udid:   udid,
			CodeId: codeid,
		}).
		Find(lists)
	err = t.Error
	return
}
func (*CodeUseStore) Create(ctx *context.Context, entity *CodeUseEntity) error {
	t := TraceDB(ctx).
		Create(entity)
	return t.Error
}
func (*CodeUseStore) Update(ctx *context.Context, entity *CodeUseEntity, fields []string) error {
	t := TraceDB(ctx).
		Model(entity).
		Select(fields).
		Updates(*entity)
	return t.Error
}
func (*CodeUseStore) Updates(ctx *context.Context, entity *CodeUseEntity) error {
	t := TraceDB(ctx).
		Model(entity).
		Select("*").
		Updates(*entity)
	return t.Error
}
func (*CodeUseStore) Delete(ctx *context.Context, id int) (err error) {
	t := TraceDB(ctx).
		Delete(&CodeUseEntity{
			ID: id,
		})
	err = t.Error
	return
}
func (*CodeUseStore) InitSchema() {
	if err := models.GetDbInstance().AutoMigrate(&CodeUseEntity{}); err != nil {
		panic(err)
	}
}
