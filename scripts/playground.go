package main

import (
	"Medscribe/reports"
	"Medscribe/user"
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tjarratt/babble"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Helper struct {
	userCol   *mongo.Collection
	reportCol *mongo.Collection
}

func NewHelper() (Helper, error) {
	if err := godotenv.Load(".env"); err != nil {
		return Helper{}, fmt.Errorf("error loading .env file: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return Helper{}, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return Helper{}, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	db := client.Database(os.Getenv("MONGODB_DB"))
	userCol := db.Collection(os.Getenv("MONGODB_USER_COLLECTION_TEST"))
	reportsColl := db.Collection(os.Getenv("MONGODB_REPORT_COLLECTION_TEST"))

	return Helper{userCol: userCol, reportCol: reportsColl}, nil

}

func createUser(name, email, password string) error {
	helper, err := NewHelper()
	if err != nil {
		return fmt.Errorf("error creating test helper: %w", err)
	}
	ctx := context.Background()
	userStore := user.NewUserStore(helper.userCol)
	_, err = userStore.Put(ctx, name, email, password)
	return err

}

func generateReports(providerID string, amount int) error {
	helper, err := NewHelper()
	if err != nil {
		return fmt.Errorf("error creating test helper: %w", err)
	}

	ctx := context.Background()

	babbler := babble.NewBabbler()
	babbler.Separator = " "

	for i := 0; i < amount; i++ {
		babbler.Count = 1
		name := babbler.Babble()

		pronouns := []string{"HE", "SHE", "THEY"}
		patientOrClient := []string{"Patient", "Client"}
		randomPatientOrClient := pronouns[rand.IntN(len(patientOrClient))]
		randomPronoun := pronouns[rand.IntN(len(pronouns))]

		babbler.Count = 20
		report := reports.Report{
			Name:               name,
			TimeStamp:          primitive.NewDateTimeFromTime(time.Now()),
			Duration:           rand.Float64() * 1000,
			ProviderID:         providerID,
			PatientOrClient:    randomPatientOrClient,
			FinishedGenerating: true,
			IsFollowUp:         rand.IntN(2) == 0,
			Pronouns:           randomPronoun,
			Subjective:         reports.ReportContent{Loading: true, Data: babbler.Babble()},
			Objective:          reports.ReportContent{Loading: true, Data: babbler.Babble()},
			Assessment:         reports.ReportContent{Loading: true, Data: babbler.Babble()},
			Planning:           reports.ReportContent{Loading: true, Data: babbler.Babble()},
			Summary:            reports.ReportContent{Loading: true, Data: babbler.Babble()},
			OneLinerSummary:    babbler.Babble(),
			ShortSummary:       babbler.Babble(),
		}

		insertResp, err := helper.reportCol.InsertOne(ctx, report)
		if err != nil {
			return fmt.Errorf("failed to insert report: %v", err)
		}

		insertID, ok := insertResp.InsertedID.(primitive.ObjectID)
		if !ok {
			return fmt.Errorf("unexpected type for InsertedID: %T", insertID)
		}

		fmt.Printf("Inserted report with ID: %v\n", insertID)
	}
	return nil
}

func wipeCollection(col *mongo.Collection) error {
	ctx := context.Background()
	result, err := col.DeleteMany(ctx, bson.M{})
	fmt.Println("deleted result: ", result)
	return err
}

func main() {
	helper, err := NewHelper()
	if err != nil {
		return
	}
	wipeCollection(helper.reportCol)
	wipeCollection(helper.userCol)
	// generateReports("67acf58a1a987f31bfe008f7", 1)
}
