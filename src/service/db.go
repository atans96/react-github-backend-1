package service

import (
	"context"
	"fmt"
	fiberSession "github.com/gofiber/fiber/v2/middleware/session"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"server_go/src/server/server/routes/types"
)

type sessionStore struct {
	DB *fiberSession.Store
}

var SessionStore sessionStore

func (s *sessionStore) NewSessionStore() *sessionStore {
	s.DB = fiberSession.New()
	s.DB.CookieSecure = true
	s.DB.CookieHTTPOnly = true
	s.DB.CookiePath = "/"
	s.DB.CookieDomain = "127.0.0.1"
	s.DB.CookieSameSite = "strict"
	sess := new(sessionStore)
	return sess
}

type mongo struct {
	DB      *mongoDB.Collection
	Session *mongoDB.Client
}

var Mongo mongo //global variable to remember internal state of struct
func (m *mongo) NewDatastore() *mongo {
	m.connect()
	if m.DB != nil && m.Session != nil {
		mongoDataStore := new(mongo)
		mongoDataStore.DB = m.DB
		mongoDataStore.Session = m.Session
		return mongoDataStore
	} else {
		panic("nil DB mongo")
	}
	return nil
}

func (m *mongo) connect() {
	fmt.Println("connected to MongoDB")
	m.DB, m.Session = m.connectToMongo()
	go watch(m.DB)
}

func (m *mongo) connectToMongo() (a *mongoDB.Collection, b *mongoDB.Client) {

	var err error
	session, err := mongoDB.NewClient(options.Client().ApplyURI(os.Getenv("DATABASE")))
	if err != nil {
		panic(err)
	}
	err = session.Connect(context.Background())
	if err != nil {
		return nil, nil
	}
	if err != nil {
		panic(err)
	}

	var DB = session.Database("github").Collection("users")
	return DB, session
}

func watch(collection *mongoDB.Collection) {
	events := make(chan *types.ChangeEvent)
	go func() {
		cs, err := collection.Watch(context.TODO(), mongoDB.Pipeline{})
		if err != nil {
			panic(err)
		}
		// Whenever there is a new change event, decode the change event and print some information about it
		for cs.Next(context.TODO()) {
			event := &types.ChangeEvent{}
			err := cs.Decode(event)
			if err != nil {
				panic(err)
			}
			switch event.OperationType {
			case "update", "replace":
				events <- event
			case "insert":
				events <- event
				fmt.Println("new user created! need to trigger something?")
			}
		}
	}()
	select {
	case <-events:
		KafkaProducer(events)
	}
}
