package model

import (
	"context"
	"sip-monitor/src/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository implements Repository for MongoDB
type MongoRepository struct {
	db                       *mongo.Database
	recordCollection         *mongo.Collection
	recordCallCollection     *mongo.Collection
	recordRegisterCollection *mongo.Collection
	userCollection           *mongo.Collection
}

// NewMongoRepository creates a new MongoDB repository
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db:                       db,
		recordCollection:         db.Collection("call_records"),
		recordCallCollection:     db.Collection("call_records_call"),
		recordRegisterCollection: db.Collection("call_records_register"),
		userCollection:           db.Collection("users"),
	}
}

// buildMongoFilter constructs a MongoDB filter based on search parameters
func (r *MongoRepository) buildMongoFilter(params entity.SearchParams) bson.M {
	filter := bson.M{}

	if params.BeginTime != nil && params.EndTime != nil {
		filter["create_time"] = bson.M{
			"$gte": params.BeginTime,
			"$lte": params.EndTime,
		}
	}

	if params.SipCallID != "" {
		filter["sip_call_id"] = params.SipCallID
	}

	if params.FromUser != "" {
		filter["from_user"] = bson.M{"$regex": params.FromUser, "$options": "i"}
	}

	if params.SrcHost != "" {
		filter["src_host"] = params.SrcHost
	}

	if params.ToUser != "" {
		filter["to_user"] = bson.M{"$regex": params.ToUser, "$options": "i"}
	}

	if params.DstHost != "" {
		filter["dst_host"] = params.DstHost
	}

	return filter
}

// getSortOptions returns MongoDB sort options based on search parameters
func (r *MongoRepository) getSortOptions(params entity.SearchParams) *options.FindOptions {
	opts := options.Find()

	// Apply pagination
	if params.Page > 0 && params.PageSize > 0 {
		skip := int64((params.Page - 1) * params.PageSize)
		limit := int64(params.PageSize)
		opts.SetSkip(skip)
		opts.SetLimit(limit)
	}

	// Apply sorting
	if params.SortBy != "" {
		sortValue := 1 // ASC
		if params.SortDesc {
			sortValue = -1 // DESC
		}
		opts.SetSort(bson.D{{Key: params.SortBy, Value: sortValue}})
	} else {
		// Default to sorting by creation time descending
		opts.SetSort(bson.D{{Key: "create_time", Value: -1}})
	}

	return opts
}

// Record operations

// CreateRecord creates a new record in MongoDB
func (r *MongoRepository) CreateRecord(ctx context.Context, record *entity.Record) error {
	if record.CreateTime.IsZero() {
		record.CreateTime = time.Now()
	}
	_, err := r.recordCollection.InsertOne(ctx, record)
	return err
}

// GetRecordByID retrieves a record by ID from MongoDB
func (r *MongoRepository) GetRecordByID(ctx context.Context, id string) (*entity.Record, error) {
	var record entity.Record
	err := r.recordCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// GetRecordList retrieves records based on search parameters from MongoDB
func (r *MongoRepository) GetRecordList(ctx context.Context, params entity.SearchParams) ([]entity.Record, *entity.Meta, error) {
	var records []entity.Record
	filter := r.buildMongoFilter(params)
	opts := r.getSortOptions(params)

	// Count total documents
	totalCount, err := r.recordCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	// Execute query
	cursor, err := r.recordCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &records); err != nil {
		return nil, nil, err
	}

	// Create meta information
	meta := &entity.Meta{
		Total:    totalCount,
		PageSize: int64(params.PageSize),
	}

	return records, meta, nil
}

// DeleteRecord deletes a record from MongoDB
func (r *MongoRepository) DeleteRecord(ctx context.Context, id string) error {
	_, err := r.recordCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoRepository) GetRecordsBySIPCallID(ctx context.Context, sipCallID string) ([]entity.Record, error) {
	var records []entity.Record
	filter := bson.M{"sip_call_id": sipCallID}
	opts := r.getSortOptions(entity.SearchParams{})
	cursor, err := r.recordCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	return records, nil
}

// SIP Call record operations

// CreateSIPCallRecord creates a new SIP call record in MongoDB
func (r *MongoRepository) CreateSIPCallRecord(ctx context.Context, record *entity.SIPRecordCall) error {
	if record.CreateTime.IsZero() {
		now := time.Now()
		record.CreateTime = &now
	}
	_, err := r.recordCallCollection.InsertOne(ctx, record)
	return err
}

// GetSIPCallRecordByID retrieves a SIP call record by ID from MongoDB
func (r *MongoRepository) GetSIPCallRecordByID(ctx context.Context, id string) (*entity.SIPRecordCall, error) {
	var record entity.SIPRecordCall
	err := r.recordCallCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// GetSIPCallRecordList retrieves SIP call records based on search parameters from MongoDB
func (r *MongoRepository) GetSIPCallRecordList(ctx context.Context, params entity.SearchParams) ([]entity.SIPRecordCall, *entity.Meta, error) {
	var records []entity.SIPRecordCall
	filter := r.buildMongoFilter(params)
	opts := r.getSortOptions(params)

	// Count total documents
	totalCount, err := r.recordCallCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	// Execute query
	cursor, err := r.recordCallCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &records); err != nil {
		return nil, nil, err
	}

	// Create meta information
	meta := &entity.Meta{
		Total:    totalCount,
		PageSize: int64(params.PageSize),
	}

	return records, meta, nil
}

// DeleteSIPCallRecord deletes a SIP call record from MongoDB
func (r *MongoRepository) DeleteSIPCallRecord(ctx context.Context, id string) error {
	_, err := r.recordCallCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// SIP Register record operations

// CreateSIPRegisterRecord creates a new SIP register record in MongoDB
func (r *MongoRepository) CreateSIPRegisterRecord(ctx context.Context, record *entity.SIPRecordRegister) error {
	if record.CreateTime.IsZero() {
		now := time.Now()
		record.CreateTime = &now
	}
	_, err := r.recordRegisterCollection.InsertOne(ctx, record)
	return err
}

// GetSIPRegisterRecordByID retrieves a SIP register record by ID from MongoDB
func (r *MongoRepository) GetSIPRegisterRecordByID(ctx context.Context, id string) (*entity.SIPRecordRegister, error) {
	var record entity.SIPRecordRegister
	err := r.recordRegisterCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// GetSIPRegisterRecordList retrieves SIP register records based on search parameters from MongoDB
func (r *MongoRepository) GetSIPRegisterRecordList(ctx context.Context, params entity.SearchParams) ([]entity.SIPRecordRegister, *entity.Meta, error) {
	var records []entity.SIPRecordRegister
	filter := r.buildMongoFilter(params)
	opts := r.getSortOptions(params)

	// Count total documents
	totalCount, err := r.recordRegisterCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	// Execute query
	cursor, err := r.recordRegisterCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &records); err != nil {
		return nil, nil, err
	}

	// Create meta information
	meta := &entity.Meta{
		Total:    totalCount,
		PageSize: int64(params.PageSize),
	}

	return records, meta, nil
}

// DeleteSIPRegisterRecord deletes a SIP register record from MongoDB
func (r *MongoRepository) DeleteSIPRegisterRecord(ctx context.Context, id string) error {
	_, err := r.recordRegisterCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// User operations

// CreateUser creates a new user
func (r *MongoRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if user.CreateAt.IsZero() {
		user.CreateAt = time.Now()
	}
	user.UpdateAt = time.Now()
	_, err := r.userCollection.InsertOne(ctx, user)
	return err
}

// GetUserByID retrieves a user by ID
func (r *MongoRepository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	err := r.userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *MongoRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user
func (r *MongoRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	user.UpdateAt = time.Now()
	_, err := r.userCollection.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	return err
}

// DeleteUser deletes a user by ID
func (r *MongoRepository) DeleteUser(ctx context.Context, id int64) error {
	_, err := r.db.Collection("users").DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoRepository) CreateDefaultAdminUser(ctx context.Context) error {
	return nil
}

func (r *MongoRepository) GetUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	cursor, err := r.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
