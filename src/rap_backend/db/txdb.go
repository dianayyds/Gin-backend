package db

import (
	"context"
	"errors"
	"github.com/cihub/seelog"
	"gorm.io/gorm"
)

// TxDBHandler 事务处理方法
type TxDBHandler func(tx *gorm.DB) (interface{}, error)

// TxDBGetter 获取 db 连接方法
type TxDBGetter func(ctx context.Context) (*gorm.DB, error)

// Option 事务的选项方法
type Option func(txDB *TxDB)

// TxDB 新的事务管理器
type TxDB struct {
	tx       *gorm.DB
	commit   bool
	dbname   string
	handler  TxDBHandler
	dbGetter TxDBGetter
}

// WithTxDBName 设置 dbname，默认为 constant.DBName
func WithTxDBName(dbname string) Option {
	return func(txDB *TxDB) {
		txDB.dbname = dbname
	}
}

// WithTxCommit 设置是否主动提交，默认为 true
func WithTxCommit(commit bool) Option {
	return func(txDb *TxDB) {
		txDb.commit = commit
	}
}

// WithTx 设置使用传入的 db 连接
func WithTx(tx *gorm.DB) Option {
	return func(txDb *TxDB) {
		txDb.tx = tx
	}
}

// WithDbGetter 设置一个获取 db 连接的方法
func WithDbGetter(dbGetter TxDBGetter) Option {
	return func(txDb *TxDB) {
		txDb.dbGetter = dbGetter
	}
}

// 统一设置所有选项值
func (txDB *TxDB) applyOption(options ...Option) {
	for _, option := range options {
		option(txDB)
	}
}

// NewTxDB 新建一个事务管理
func NewTxDB(handler TxDBHandler, options ...Option) *TxDB {
	txDB := &TxDB{
		commit:  true,
		handler: handler,
	}

	txDB.applyOption(options...)
	return txDB
}

// Execute 执行这个事务
func (txDB *TxDB) Execute(ctx context.Context, options ...Option) (interface{}, error) {
	txDB.applyOption(options...)

	// 事务配置不正确，无法获取到一个 db 连接
	if txDB.tx == nil && txDB.dbGetter == nil {
		seelog.Errorf("初始化事务错误：缺少默认的 db 配置")
		return nil, errors.New("db inner error")
	}

	if txDB.tx == nil {
		// 获取 db 连接
		client, err := txDB.dbGetter(ctx)
		if err != nil {
			return nil, err
		}

		// 开启一个事务
		txDB.tx = client.Begin()
		if txDB.commit {
			defer txDB.tx.Rollback()
		}
	}

	// 调用实际的业务逻辑
	response, err := txDB.handler(txDB.tx)
	if err != nil {
		return response, err
	}

	if txDB.commit {
		txDB.tx.Commit()
	}

	return response, err
}
