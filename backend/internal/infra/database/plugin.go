package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// UUIDPlugin は、INSERT時に自動的にUUID IDを付与し、ユーザーID、バージョン情報を管理するGormプラグイン
type UUIDPlugin struct{}

// Name はプラグインの名前を返す
func (p *UUIDPlugin) Name() string {
	return "UUIDPlugin"
}

// Initialize はプラグインを初期化する
func (p *UUIDPlugin) Initialize(db *gorm.DB) error {
	// Createコールバックを登録
	err := db.Callback().Create().Before("gorm:create").Register("uuid:before_create", p.beforeCreate)
	if err != nil {
		return err
	}

	// Updateコールバック を登録（楽観ロックプラグインの後に実行）
	return db.Callback().Update().After("optimistic_lock:before_update").Register("uuid:before_update", p.beforeUpdate)
}

// beforeCreate は、レコード作成前にUUID ID、CreateUserID、Versionを付与する
func (p *UUIDPlugin) beforeCreate(db *gorm.DB) {
	if db.Statement.Schema != nil {
		// ReflectValueがスライスかどうかをチェック
		reflectValue := db.Statement.ReflectValue

		// スライスの場合は各要素に対して値を設定
		if reflectValue.Kind() == reflect.Slice {
			for i := 0; i < reflectValue.Len(); i++ {
				elem := reflectValue.Index(i)
				p.setFieldsForCreate(db, elem)
			}
		} else {
			// 単一のレコードの場合
			p.setFieldsForCreate(db, reflectValue)
		}
	}
}

// beforeUpdate は、レコード更新前にUpdateUserID、Versionを付与する
func (p *UUIDPlugin) beforeUpdate(db *gorm.DB) {
	if db.Statement.Schema != nil {
		// ReflectValueがスライスかどうかをチェック
		reflectValue := db.Statement.ReflectValue

		// スライスの場合は各要素に対して値を設定
		if reflectValue.Kind() == reflect.Slice {
			for i := 0; i < reflectValue.Len(); i++ {
				elem := reflectValue.Index(i)
				p.setFieldsForUpdate(db, elem)
			}
		} else {
			// 単一のレコードの場合
			p.setFieldsForUpdate(db, reflectValue)
		}
	}
}

// setFieldsForCreate は、作成時に必要なフィールドを設定する
func (p *UUIDPlugin) setFieldsForCreate(db *gorm.DB, reflectValue reflect.Value) {
	// IDフィールドの設定
	if field := db.Statement.Schema.LookUpField("ID"); field != nil {
		p.setUUIDIfEmpty(db, field, reflectValue)
	}

	// CreateUserIDフィールドの設定
	if field := db.Statement.Schema.LookUpField("CreateUserID"); field != nil {
		p.setUserIDFromContext(db, field, reflectValue)
	}

	// Versionフィールドの設定
	if field := db.Statement.Schema.LookUpField("Version"); field != nil {
		p.setVersionForCreate(db, field, reflectValue)
	}
}

// setFieldsForUpdate は、更新時に必要なフィールドを設定する
func (p *UUIDPlugin) setFieldsForUpdate(db *gorm.DB, reflectValue reflect.Value) {
	// UpdateUserIDフィールドの設定
	if field := db.Statement.Schema.LookUpField("UpdateUserID"); field != nil {
		p.setUserIDFromContext(db, field, reflectValue)
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

// setUserIDFromContext は、contextからユーザーIDを取得してフィールドに設定する
func (p *UUIDPlugin) setUserIDFromContext(db *gorm.DB, field *schema.Field, reflectValue reflect.Value) {
	// contextからユーザーIDを取得
	userID := context.GetCtxFromUser(db.Statement.Context)
	if userID != "" {
		// ユーザーIDフィールドに値を設定
		_ = field.Set(db.Statement.Context, reflectValue, userID)
	}
}

// setVersionForCreate は、作成時にVersionフィールドを1に設定する
func (p *UUIDPlugin) setVersionForCreate(db *gorm.DB, field *schema.Field, reflectValue reflect.Value) {
	// Versionフィールドの値を取得
	_, isZero := field.ValueOf(db.Statement.Context, reflectValue)

	// Versionがゼロ値の場合のみ0を設定（明示的に値が設定されている場合は変更しない）
	if isZero {
		_ = field.Set(db.Statement.Context, reflectValue, 1)
	}
}

// NewUUIDPlugin は、新しいUUIDPluginインスタンスを作成する
func NewUUIDPlugin() *UUIDPlugin {
	return &UUIDPlugin{}
}

// OptimisticLockPlugin は楽観的排他制御を行うGormプラグイン
type OptimisticLockPlugin struct{}

// Name はプラグインの名前を返す
func (p *OptimisticLockPlugin) Name() string {
	return "OptimisticLockPlugin"
}

// Initialize はプラグインを初期化する
func (p *OptimisticLockPlugin) Initialize(db *gorm.DB) error {
	// Updateコールバックを登録（楽観ロックチェック）
	// UUIDPluginのbeforeUpdateより前に実行されるように別の位置に登録
	err := db.Callback().Update().Before("gorm:before_update").Register("optimistic_lock:before_update", p.beforeUpdate)
	if err != nil {
		return err
	}

	// 更新後のコールバックを登録（影響行数チェック）
	return db.Callback().Update().After("gorm:update").Register("optimistic_lock:after_update", p.afterUpdate)
}

// beforeUpdate は、更新前に楽観ロックチェック用のWHERE条件を追加する
func (p *OptimisticLockPlugin) beforeUpdate(db *gorm.DB) {
	if db.Statement.Schema != nil {
		// Versionフィールドが存在するかチェック
		if field := db.Statement.Schema.LookUpField("Version"); field != nil {
			// ReflectValueがスライスかどうかをチェック
			reflectValue := db.Statement.ReflectValue

			if reflectValue.Kind() == reflect.Slice {
				// スライスの場合は各要素に対して処理
				for i := 0; i < reflectValue.Len(); i++ {
					elem := reflectValue.Index(i)
					p.addOptimisticLockCondition(db, field, elem)
				}
			} else {
				// 単一のレコードの場合
				p.addOptimisticLockCondition(db, field, reflectValue)
			}
		}
	}
}

// afterUpdate は、更新後に影響行数をチェックして楽観ロック例外を発生させる
func (p *OptimisticLockPlugin) afterUpdate(db *gorm.DB) {
	// エラーが既に発生している場合は何もしない
	if db.Error != nil {
		return
	}

	// Versionフィールドが存在するスキーマかチェック
	if db.Statement.Schema != nil {
		if field := db.Statement.Schema.LookUpField("Version"); field != nil {
			// 影響行数が0の場合は楽観ロックエラーとして扱う
			if db.RowsAffected == 0 {
				db.AddError(errors.ErrOptimisticLock)
			}
		}
	}
}

// addOptimisticLockCondition は、楽観ロックチェック用のWHERE条件を追加し、バージョンをインクリメントする
func (p *OptimisticLockPlugin) addOptimisticLockCondition(db *gorm.DB, field *schema.Field, reflectValue reflect.Value) {
	// 現在のVersionフィールドの値を取得
	fieldValue, _ := field.ValueOf(db.Statement.Context, reflectValue)

	// Versionフィールドの値をWHERE条件に追加
	if currentVersion, ok := fieldValue.(int); ok {
		// 元のバージョンをWHERE条件に追加
		db.Where(fmt.Sprintf("%s = ?", field.DBName), currentVersion)

		// バージョンを明示的にインクリメントしてSETに追加
		newVersion := currentVersion + 1
		_ = field.Set(db.Statement.Context, reflectValue, newVersion)

		// Statement内のSETフィールドを直接設定してバージョンをインクリメント
		if db.Statement.Dest != nil {
			db.Statement.SetColumn(field.DBName, newVersion)
		}
	} else {
		// Versionが取得できない場合はエラー
		db.AddError(errors.ErrVersionNotFound)
	}
}

// NewOptimisticLockPlugin は、新しいOptimisticLockPluginインスタンスを作成する
func NewOptimisticLockPlugin() *OptimisticLockPlugin {
	return &OptimisticLockPlugin{}
}

// ZeroValueOmitPlugin はゼロ値やsql.Null*のValidがfalseのフィールドを自動的にomitするGormプラグイン
type ZeroValueOmitPlugin struct{}

// Name はプラグインの名前を返す
func (p *ZeroValueOmitPlugin) Name() string {
	return "ZeroValueOmitPlugin"
}

// Initialize はプラグインを初期化する
func (p *ZeroValueOmitPlugin) Initialize(db *gorm.DB) error {
	// Updateコールバックを登録（ゼロ値チェック）
	return db.Callback().Update().Before("gorm:before_update").Register("zero_value_omit:before_update", p.beforeUpdate)
}

// beforeUpdate は、更新前にゼロ値やsql.Null*のValidがfalseのフィールドをomitする
func (p *ZeroValueOmitPlugin) beforeUpdate(db *gorm.DB) {
	if db.Statement.Schema != nil {
		reflectValue := db.Statement.ReflectValue

		// Check if reflectValue kind is valid (struct or slice of structs)
		// If map is passed to Updates, reflectValue might be map, which we should skip or handle differently
		if reflectValue.Kind() == reflect.Map {
			return
		}

		// スライスの場合は各要素に対して処理
		if reflectValue.Kind() == reflect.Slice {
			for i := 0; i < reflectValue.Len(); i++ {
				elem := reflectValue.Index(i)
				p.omitZeroValueFields(db, elem)
			}
		} else {
			// 単一のレコードの場合
			p.omitZeroValueFields(db, reflectValue)
		}
	}
}

// omitZeroValueFields は、ゼロ値やsql.Null*のValidがfalseのフィールドをomitする
func (p *ZeroValueOmitPlugin) omitZeroValueFields(db *gorm.DB, reflectValue reflect.Value) {
	if reflectValue.Kind() == reflect.Ptr {
		if reflectValue.IsNil() {
			return
		}
		reflectValue = reflectValue.Elem()
	}

	var omitFields []string

	// 各フィールドをチェック
	for _, field := range db.Statement.Schema.Fields {
		// CreatedAtやDeletedAtなどの自動管理フィールドはスキップ
		if field.AutoCreateTime != 0 || field.AutoUpdateTime != 0 || field.Name == "DeletedAt" {
			continue
		}

		fieldValue, isZero := field.ValueOf(db.Statement.Context, reflectValue)

		// ゼロ値の場合はomit（time.Timeのゼロ値も含む）
		if isZero {
			omitFields = append(omitFields, field.DBName)
			continue
		}

		// sql.Null*系の型をチェック
		if p.isInvalidNullType(fieldValue) {
			omitFields = append(omitFields, field.DBName)
		}
	}

	// omitするフィールドがある場合は適用
	if len(omitFields) > 0 {
		db.Statement.Omits = append(db.Statement.Omits, omitFields...)
	}
}

// isInvalidNullType は、sql.Null*系の型でValidがfalseかどうかをチェック
func (p *ZeroValueOmitPlugin) isInvalidNullType(value interface{}) bool {
	switch v := value.(type) {
	case sql.NullString:
		return !v.Valid
	case sql.NullInt64:
		return !v.Valid
	case sql.NullInt32:
		return !v.Valid
	case sql.NullInt16:
		return !v.Valid
	case sql.NullByte:
		return !v.Valid
	case sql.NullFloat64:
		return !v.Valid
	case sql.NullBool:
		return !v.Valid
	case sql.NullTime:
		return !v.Valid
	default:
		// driver.Valuerインターフェースを実装しているカスタムNull型もチェック
		if valuer, ok := value.(driver.Valuer); ok {
			val, err := valuer.Value()
			if err != nil || val == nil {
				return true
			}
		}
	}
	return false
}

// NewZeroValueOmitPlugin は、新しいZeroValueOmitPluginインスタンスを作成する
func NewZeroValueOmitPlugin() *ZeroValueOmitPlugin {
	return &ZeroValueOmitPlugin{}
}
