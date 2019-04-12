package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/zhiburt/openbox/monitor/services/monitor"
	"go.mongodb.org/mongo-driver/bson"
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
	id, _ := uuid.NewV4()
	f.ID = id.String()

	cleanBody(&f)

	var root monitor.File
	level.Debug(*rep.logger).Log("get file", fmt.Sprint(f))
	if err := rep.files.FindOne(ctx, bson.D{{Key: "ownerid", Value: f.OwnerID}}).Decode(&root); err != nil {
		level.Error(*rep.logger).Log("error didn't find root", err)

		_, err := rep.files.InsertOne(ctx, monitor.File{
			OwnerID:   f.OwnerID,
			Name:      "/",
			Status:    "root",
			CreatedOn: time.Now().UnixNano(),
			IsFolder:  true,
			Files:     []monitor.File{f},
		})
		if err != nil {
			level.Error(*rep.logger).Log("err in first <IF>", err)
			return "", err
		}

		return f.ID, err
	}
	d, _ := json.MarshalIndent(root, " ", "")
	level.Debug(*rep.logger).Log("found root", fmt.Sprint(string(d)))

	updatedfile, err := insertInto(root, f)
	if err != nil {
		level.Error(*rep.logger).Log("err cannot find place for file, structure your file cantains erorr ", err)
		return "", err
	}

	_, err = rep.files.ReplaceOne(ctx, bson.D{{Key: "ownerid", Value: f.OwnerID}}, updatedfile)
	if err != nil {
		level.Error(*rep.logger).Log("err", err)
		return "", err
	}

	level.Info(*rep.logger).Log("msg", "created new file in repo")
	return f.ID, nil
}

func (*MongoRepo) GetFileByID(ctx context.Context, id string) (monitor.File, error) {
	return monitor.File{}, nil
}

func (rep *MongoRepo) GetFilesByOwner(ctx context.Context, id string) ([]monitor.File, error) {
	var root monitor.File
	if err := rep.files.FindOne(ctx, bson.D{{Key: "ownerid", Value: id}}).Decode(&root); err != nil {
		level.Error(*rep.logger).Log("error didn't find users files", err, "user", id)
		return nil, err
	}

	return root.Files, nil
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

func insertInto(root monitor.File, insertedFile monitor.File) (monitor.File, error) {
	r := &root
	err := insertInPlace(&root, &insertedFile)
	return *r, err
}

func insertInPlace(root, insertedFile *monitor.File) error {
	if insertedFile.IsFolder == false || insertedFile.Files == nil {
		root.Files = append(root.Files, *insertedFile)
		return nil
	}

	for i := 0; i < len(root.Files); i++ {
		if root.Files[i].Name == insertedFile.Name && root.Files[i].IsFolder {
			return insertInPlace(&root.Files[i], &insertedFile.Files[0])
		}
	}

	return fmt.Errorf("canno't find such directory")
}

func cleanBody(f *monitor.File) {
	if f == nil {
		return
	}

	f.Body = nil
	for i := 0; i < len(f.Files); i++ {
		cleanBody(&f.Files[i])
	}
}
