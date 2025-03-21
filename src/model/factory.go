package model

import (
	"context"
	"fmt"
	"log"
	"sip-monitor/src/config"
	"sip-monitor/src/entity"
	"sip-monitor/src/model/sql"
	"time"

	mongorepo "sip-monitor/src/model/mongo"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	DBTypeMySQL    = "mysql"
	DBTypeSQLite   = "sqlite"
	DBTypeMongoDB  = "mongodb"
	DBTypePostgres = "postgres"
)

// RepositoryFactory creates the appropriate repository implementation based on database type
type RepositoryFactory struct{}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory() *RepositoryFactory {
	return &RepositoryFactory{}
}

// Helper function to initialize a repository
func InitRepository(cfg *config.Config) (Repository, error) {
	factory := NewRepositoryFactory()

	if cfg.DBType == "" {
		cfg.DBType = DBTypeSQLite
	}

	return factory.CreateRepository(cfg)
}

// CreateRepository creates the appropriate repository based on configuration
func (f *RepositoryFactory) CreateRepository(cfg *config.Config) (Repository, error) {
	switch cfg.DBType {
	case DBTypeMySQL:
		return f.createMySQLRepository(cfg)
	case DBTypeSQLite:
		return f.createSQLiteRepository(cfg)
	case DBTypeMongoDB:
		return f.createMongoRepository(cfg)
	case DBTypePostgres:
		return f.createPostgresRepository(cfg)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}
}

// createMySQLRepository creates a MySQL repository
func (f *RepositoryFactory) createMySQLRepository(cfg *config.Config) (Repository, error) {
	dsn := cfg.DSNURL
	if dsn == "" {
		// Construct DSN from individual components
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser, cfg.DBPassword, cfg.DBAddr, cfg.DBName)
	}

	db, err := f.openGormDB(mysql.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// Auto-migrate schema
	if err := f.migrateSchema(db); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	return sql.NewGormRepository(db), nil
}

// createSQLiteRepository creates a SQLite repository
func (f *RepositoryFactory) createSQLiteRepository(cfg *config.Config) (Repository, error) {
	filePath := cfg.DBPath
	if filePath == "" {
		filePath = "sip_monitor.db" // Default SQLite database file
	}

	db, err := f.openGormDB(sqlite.Open(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	// Auto-migrate schema
	if err := f.migrateSchema(db); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	return sql.NewGormRepository(db), nil
}

// createMongoRepository creates a MongoDB repository
func (f *RepositoryFactory) createMongoRepository(cfg *config.Config) (Repository, error) {
	// Use existing MongoDB initialization
	//todo::实现MongoDB的初始化

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if cfg.DSNURL == "" {
		logrus.Info("MongoDBInit BY DBUser、DBPassword、DBAddr")
		cfg.DSNURL = fmt.Sprintf("mongodb://%s:%s@%s", cfg.DBUser, cfg.DBPassword, cfg.DBAddr)
	} else {
		logrus.Info("MongoDBInit BY DSN_URL")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DSNURL))
	if err != nil {
		logrus.WithError(err).Error("MongoDBInit Error")
		return nil, err
	}

	md := client.Database(cfg.DBName)

	index := mongo.IndexModel{
		Keys:    bson.M{"sip_call_id": 1},
		Options: options.Index().SetUnique(true),
	}
	_, _ = md.Collection("call_records_call").Indexes().CreateOne(context.Background(), index)
	_, _ = md.Collection("call_records_register").Indexes().CreateOne(context.Background(), index)

	// Create MongoDB repository
	return mongorepo.NewMongoRepository(md), nil
}

// createPostgresRepository creates a PostgreSQL repository
func (f *RepositoryFactory) createPostgresRepository(cfg *config.Config) (Repository, error) {
	dsn := cfg.DSNURL
	if dsn == "" {
		// Construct DSN from individual components
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			cfg.DBAddr, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	}

	db, err := f.openGormDB(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Auto-migrate schema
	if err := f.migrateSchema(db); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	return sql.NewGormRepository(db), nil
}

func (f *RepositoryFactory) openGormDB(dialector gorm.Dialector) (*gorm.DB, error) {
	// Configure GORM logger
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Configure GORM
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Use singular table names
		},
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// migrateSchema migrates the database schema
func (f *RepositoryFactory) migrateSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Record{},
		&entity.RecordRaw{},
		&entity.Call{},
		&entity.User{},
		&entity.Gateway{},
	)
}
