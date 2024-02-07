package entity

type MongoDBStatsVO struct {
	Name       string `bson:"-" json:"name"`
	Capped     bool   `bson:"capped" json:"capped"`        //表示集合是否是固定大小的集合
	Count      int64  `bson:"count" json:"count"`          //集合中文档的总数
	IndexCount int64  `bson:"nindexes" json:"index_count"` //集合中索引的数量

	AvgObjSize  int64 `bson:"avgObjSize" json:"avg_obj_size"`  //平均对象大小（字节）,这个数字表示存储在集合中的文档平均大小。
	Size        int64 `bson:"size" json:"size"`                //集合中所有文档的总大小（字节）
	StorageSize int64 `bson:"storageSize" json:"storage_size"` //分配给集合的物理存储空间的总大小（字节）

	TotalIndexSize int64 `bson:"totalIndexSize" json:"total_index_size"` //所有索引占用的总空间大小（字节）

	FreeStorageSize int64 `bson:"freeStorageSize" json:"free_storage_size"` //空闲存储空间的大小（字节）
	TotalSize       int64 `bson:"totalSize" json:"total_size"`              //集合的总大小，包括所有文档和索引的大小
}
