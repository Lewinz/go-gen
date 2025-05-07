package {{.PackageName}}

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	{{.TypePascal}} struct {
		Id          string    `bson:"_id,omitempty" json:"id,omitempty"`
		// TOTO: Add your fields here
		CreatedTime time.Time `bson:"createdTime"   json:"createdTime"`
		UpdatedTime time.Time `bson:"updatedTime"   json:"updatedTime"`
	}

	{{.TypePascal}}Model interface {
		Insert(ctx context.Context, {{.TypeCamel}} *{{.TypePascal}}) error
		Update(ctx context.Context, {{.TypeCamel}} *{{.TypePascal}}) error
		Delete(ctx context.Context, id string) error
		FindById(ctx context.Context, id string) (*{{.TypePascal}}, error)
		Search(ctx context.Context, cond *{{.TypePascal}}Cond) ([]*{{.TypePascal}}, error)
	}

	default{{.TypePascal}}Model struct {
		model *mongo.Collection
	}

	{{.TypePascal}}Cond struct {
		Id  string
		Ids []string
	}
)

func New{{.TypePascal}}Model(db *mongo.Database) {{.TypePascal}}Model {
	return &default{{.TypePascal}}Model{
		model: db.Collection("{{.TypeSnake}}"),
	}
}

func (m *default{{.TypePascal}}Model) Insert(ctx context.Context, {{.TypeCamel}} *{{.TypePascal}}) error {
	{{.TypeCamel}}.CreatedTime = time.Now()
	{{.TypeCamel}}.UpdatedTime = time.Now()
	
	_, err := m.model.InsertOne(ctx, {{.TypeCamel}})
	return err
}

func (m *default{{.TypePascal}}Model) Update(ctx context.Context, {{.TypeCamel}} *{{.TypePascal}}) error {
	{{.TypeCamel}}.UpdatedTime = time.Now()
	
	_, err := m.model.UpdateOne(
		ctx,
		bson.M{"_id": {{.TypeCamel}}.Id},
		bson.M{"$set": {{.TypeCamel}}},
	)
	return err
}

func (m *default{{.TypePascal}}Model) Delete(ctx context.Context, id string) error {
	_, err := m.model.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (m *default{{.TypePascal}}Model) FindById(ctx context.Context, id string) (*{{.TypePascal}}, error) {
	var {{.TypeCamel}} {{.TypePascal}}
	err := m.model.FindOne(ctx, bson.M{"_id": id}).Decode(&{{.TypeCamel}})
	if err != nil {
		return nil, err
	}
	return &{{.TypeCamel}}, nil
}

func (c *{{.TypePascal}}Cond) genCond() bson.M {
	filter := bson.M{}

	if c.Id != "" {
		filter["_id"] = c.Id
	} else if len(c.Ids) > 0 {
		filter["_id"] = bson.M{"$in": c.Ids}
	}

	return filter
}

func (m *default{{.TypePascal}}Model) Search(ctx context.Context, cond *{{.TypePascal}}Cond) ([]*{{.TypePascal}}, error) {
	var result []*{{.TypePascal}}
	filter := cond.genCond()

	cursor, err := m.model.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
} 