package repositories

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/openbox/monitor/services/monitor"
	"go.mongodb.org/mongo-driver/mongo"
)

// "go.mongodb.org/mongo-driver/mongo"

type MongoRepo struct {
	logger *log.Logger
	cmongo *mongo.Client
	files  *mongo.Collection
}

func NewRepository(client *mongo.Client, logger *log.Logger) monitor.Repository {
	files := client.Database("monitor_files").Collection("files")

	return &MongoRepo{
		logger: logger,
		cmongo: client,
		files:  files,
	}
}

func (rep *MongoRepo) CreateFile(ctx context.Context, f monitor.File) (string, error) {
	// if rep.files.FindOne(ctx, bson.D{{"user_id", f.OwnerID}}).Err() != nil {
	// 	_, err := rep.files.InsertOne(ctx, bson.D{{"user_id", f.OwnerID}, {"file", monitor.File{
	// 		OwnerID:   f.OwnerID,
	// 		Name:      "/",
	// 		Status:    "root",
	// 		CreatedOn: time.Now().UnixNano(),
	// 		IsFolder:  true,
	// 		Files:     []monitor.File{f},
	// 	}}})
	// 	if err != nil {
	// 		level.Error(*rep.logger).Log("err in first <IF>", err)
	// 		return "", err
	// 	}
	// }
	// if f.IsFolder && f.Files != nil {
	// 	rep.files.UpdateOne(ctx, bson.D{{"name", f.Name}})
	// }

	// res, err := rep.files.InsertOne(ctx, f)
	// if err != nil {
	// 	level.Error(*rep.logger).Log("err", err)
	// 	return "", err
	// }
	// level.Info(*rep.logger).Log("msg", "created new file in repo")
	// return res.InsertedID.(string), nil
	return "", nil
}

func (*MongoRepo) GetFileByID(ctx context.Context, id string) (monitor.File, error) {
	return monitor.File{}, nil
}

func (*MongoRepo) GetFilesByOwner(ctx context.Context, id string) ([]monitor.File, error) {
	return nil, nil
}

func (*MongoRepo) ChangeFileName(ctx context.Context, id, newname string) error {
	return nil
}

func (*MongoRepo) ChangeFileBody(ctx context.Context, id string, b []byte) error {
	return nil
}

func (*MongoRepo) RemoveFileByID(ctx context.Context, id string) error {
	return nil
}
