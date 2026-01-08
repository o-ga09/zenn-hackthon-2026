package database

import (
	"reflect"

	"github.com/o-ga09/zenn-hackthon-2026/pkg/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// UUIDPlugin は、INSERT時に自動的にUUID IDを付与するGormプラグイン
type UUIDPlugin struct{}

// Name はプラグインの名前を返す
func (p *UUIDPlugin) Name() string {
	return "UUIDPlugin"
}

// Initialize はプラグインを初期化する
func (p *UUIDPlugin) Initialize(db *gorm.DB) error {
	// Createコールバックを登録
	return db.Callback().Create().Before("gorm:create").Register("uuid:before_create", p.beforeCreate)
}

// beforeCreate は、レコード作成前にUUID IDを付与する
func (p *UUIDPlugin) beforeCreate(db *gorm.DB) {
	// IDフィールドが存在するかチェック
	if db.Statement.Schema != nil {
		if field := db.Statement.Schema.LookUpField("ID"); field != nil {
			// ReflectValueがスライスかどうかをチェック
			reflectValue := db.Statement.ReflectValue
			
			// スライスの場合は各要素に対してUUIDを生成
			if reflectValue.Kind() == reflect.Slice {
				for i := 0; i < reflectValue.Len(); i++ {
					elem := reflectValue.Index(i)
					p.setUUIDIfEmpty(db, field, elem)
				}
			} else {
				// 単一のレコードの場合
				p.setUUIDIfEmpty(db, field, reflectValue)
			}
		}
	}
}

// setUUIDIfEmpty は、IDフィールドが空の場合にUUIDを設定する
func (p *UUIDPlugin) setUUIDIfEmpty(db *gorm.DB, field *schema.Field, reflectValue reflect.Value) {
	// IDフィールドの値を取得
	fieldValue, isZero := field.ValueOf(db.Statement.Context, reflectValue)
	
	// IDがゼロ値（空文字列）の場合のみUUIDを生成
	if isZero {
		// UUID V7を生成
		id, err := uuid.GenerateIDV7()
		if err != nil {
			// V7の生成に失敗した場合はV4にフォールバック
			id = uuid.GenerateID()
		}
		
		// IDフィールドに値を設定
		_ = field.Set(db.Statement.Context, reflectValue, id)
	} else if str, ok := fieldValue.(string); ok && str == "" {
		// 文字列型で空文字列の場合もUUIDを生成
		id, err := uuid.GenerateIDV7()
		if err != nil {
			id = uuid.GenerateID()
		}
		_ = field.Set(db.Statement.Context, reflectValue, id)
	}
}

// NewUUIDPlugin は、新しいUUIDPluginインスタンスを作成する
func NewUUIDPlugin() *UUIDPlugin {
	return &UUIDPlugin{}
}
