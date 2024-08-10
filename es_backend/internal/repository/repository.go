package repository

import (
    "context"
    "fmt"
    "es_backend/pkg/log"
    "es_backend/pkg/zapgorm2"
    "time"

    "github.com/elastic/go-elasticsearch/v8"
    "github.com/glebarez/sqlite"
    "github.com/redis/go-redis/v9"
    "github.com/spf13/viper"
    "go.uber.org/zap"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

const ctxTxKey = "TxKey"

type Repository struct {
    db     *gorm.DB
    rdb    *redis.Client
    logger *log.Logger
    es     *elasticsearch.Client
}

func NewRepository(logger *log.Logger, db *gorm.DB, rdb *redis.Client, es *elasticsearch.Client) *Repository {
    return &Repository{
        db:     db,
        rdb:    rdb,
        logger: logger,
        es:     es,
    }
}

type Transaction interface {
    Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTransaction(r *Repository) Transaction {
    return r
}

// DB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (r *Repository) DB(ctx context.Context) *gorm.DB {
    v := ctx.Value(ctxTxKey)
    if v != nil {
        if tx, ok := v.(*gorm.DB); ok {
            return tx
        }
    }
    return r.db.WithContext(ctx)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        ctx = context.WithValue(ctx, ctxTxKey, tx)
        return fn(ctx)
    })
}

func NewDB(conf *viper.Viper, l *log.Logger) *gorm.DB {
    var (
        db  *gorm.DB
        err error
    )

    logger := zapgorm2.New(l.Logger)
    driver := conf.GetString("data.db.user.driver")
    dsn := conf.GetString("data.db.user.dsn")

    // GORM doc: https://gorm.io/docs/connecting_to_the_database.html
    switch driver {
    case "mysql":
        db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
            Logger: logger,
        })
    case "postgres":
        db, err = gorm.Open(postgres.New(postgres.Config{
            DSN:                  dsn,
            PreferSimpleProtocol: true, // disables implicit prepared statement usage
        }), &gorm.Config{})
    case "sqlite":
        db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
    default:
        panic("unknown db driver")
    }
    if err != nil {
        panic(err)
    }
    if err != nil {
        zap.L().Error("failed to migrate db", zap.Error(err))
    }
    db = db.Debug()

    // Connection Pool config
    sqlDB, err := db.DB()
    if err != nil {
        panic(err)
    }
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    return db
}
func NewRedis(conf *viper.Viper) *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:     conf.GetString("data.redis.addr"),
        Password: conf.GetString("data.redis.password"),
        DB:       conf.GetInt("data.redis.db"),
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        panic(fmt.Sprintf("redis error: %s", err.Error()))
    }

    return rdb
}
func NewElasticSearch(conf *viper.Viper) *elasticsearch.Client {
    cfg := elasticsearch.Config{
        Addresses: []string{conf.GetString("data.elasticsearch.address")},
        Username:  conf.GetString("data.elasticsearch.username"),
        Password:  conf.GetString("data.elasticsearch.password"),
    }
    elasticClient, err := elasticsearch.NewClient(cfg)
    if err != nil {
        panic(fmt.Sprintf("elasticsearch error: %s", err.Error()))
    }
    return elasticClient
}
