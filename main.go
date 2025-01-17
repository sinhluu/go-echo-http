package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type TodoReq struct {
	Title string `json:"title" bson:"title"`
}

type Todo struct {
	ID        string    `json:"id" bson:"_id"`
	Title     string    `json:"title" bson:"title"`
	Complete  bool      `json:"complete" bson:"complete"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type DB struct {
	Mongo *mongo.Client
	Db    *mongo.Database
}

func (d *DB) Add(ctx context.Context, t Todo) error {
	_, err := d.Db.Collection("todo").InsertOne(ctx, t)
	return err
}

func (d *DB) List(ctx context.Context) ([]Todo, error) {
	cur, err := d.Db.Collection("todo").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var todos []Todo
	if err = cur.All(ctx, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

func (d *DB) Complete(ctx context.Context, id string) error {
	_, err := d.Db.Collection("todo").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"complete": true}})
	return err
}

func main() {
	uri := os.Getenv("MONGODB_URI")
	conn := options.Client().ApplyURI(uri)
	c, err := mongo.Connect(context.Background(), conn)
	if err != nil {
		log.Println("Cannot connect to DB: ", uri)
		panic(err)
	}
	d := DB{Mongo: c, Db: c.Database("todo")}

	e := echo.New()

	e.Use(middleware.Logger())

	g := e.Group("/todos")
	g.GET("", func(c echo.Context) error {
		list, err := d.List(c.Request().Context())
		if err != nil {
			log.Println("Get list error: ", err)
			return err
		}
		return c.JSON(200, list)
	})

	g.POST("", func(c echo.Context) error {
		var t TodoReq
		if err := c.Bind(&t); err != nil {
			log.Println("Bind error: ", err)
			return err
		}
		err = d.Add(c.Request().Context(), Todo{
			ID:        primitive.NewObjectID().Hex(),
			Title:     t.Title,
			Complete:  false,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}
		return c.JSON(200, t)
	})

	g.PUT("/:id", func(c echo.Context) error {
		id := c.Param("id")
		err := d.Complete(c.Request().Context(), id)
		if err != nil {
			log.Println("Complete error: ", err)
			return err
		}
		return c.JSON(200, id)
	})

	e.GET("/ifconfig", func(c echo.Context) error {
		return c.String(200, c.RealIP())
	})

	e.GET("/request", func(c echo.Context) error {
		u := c.QueryParam("url")
		if u == "" {
			return c.String(200, "empty url")
		}
		h := http.Client{}
		res, err := h.Get(u)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return c.String(200, "Result: \n"+string(b))
	})

	e.Logger.Fatal(e.Start(":4000"))
}
