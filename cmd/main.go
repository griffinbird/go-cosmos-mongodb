package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	//"os"
	//"reflect"
	"time"
	//"net/http"
	"encoding/json"
	//"io"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func newClientFromEnviroment() (*mongo.Client, error) {
	
	connString := os.Getenv("MONGODB_CONNECTION_STRING")
	if connString == "" {
		return nil, errors.New("MONGODB_CONNECTION_STRING is empty")
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)
	//defer cancel()

	// client option
	options := options.Client().ApplyURI(connString).SetDirect(true)

	// connect to mongodb
	client, err := mongo.Connect(ctx, options)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Successfully connected amd pinged MongoDB")

	return client, nil
}

func run() error {
	client, err := newClientFromEnviroment()
	if err != nil {
		return err
	}

	prompt := `Azure Cosmos DB Mongo API using Golang
-----------------------------------------
[a]   Query for single customer
[b]   Point read for single customer
[c]   List all product categories
[d]   Query products by category id
[e]   Update product category name
[f]   Query orders by customer id
[g]   Query for customer and all orders
[h]   Create new order and update order total
[i]   Delete order and update order total
[j]   Query top 10 customers
-------------------------------------------
[l]   Create databases and containers
[l]   Upload data to containers
[m]   Delete databases and containers
-------------------------------------------
[x]   Exit

> `

out:
	for {
		fmt.Print(prompt)
		result := ""
		fmt.Scanln((&result))
		fmt.Printf("\nYour selection is: %v\n", result)

		switch result {
		case "a":
			databaseName := "mongobird"
			collectionName := "customer"
			pk := "3200B375-08DC-4A36-9F6F-E325A0B46550"

			err := queryCustomer(client, databaseName, collectionName, pk)
			if err != nil {
				return err
			}
			
		case "b":
			databaseName := "mongobird"
			collectionName := "customer"
			pk := "3200B375-08DC-4A36-9F6F-E325A0B46550"

			customer, err := getCustomer(client, databaseName, collectionName, pk)
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(customer, "", "    ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", b)

			statistics, err := GetLastRequestStats(client, databaseName)
			if err != nil {
				return err
			}

			log.Printf("RequestCharge: %.2f\n", statistics.RequestCharge)


		case "l":
		/* 	if err := CreateDatabase(client); err != nil {
			return err
		} */
		//createContainer(databaseName, containerName, partitionKey)
		/*case "m":
			imports := []struct {
				URL       string
				PK        string
				Database  string
				Collection string
		}{
		/* 				{
			URL:       "https://raw.githubusercontent.com/MicrosoftDocs/mslearn-cosmosdb-modules-central/main/data/fullset/database-v2/customer",
			PK:        "id",
			Database:  "database-v2",
			Container: "customer",
		}, */
		/* 				{
			URL:       "https://raw.githubusercontent.com/MicrosoftDocs/mslearn-cosmosdb-modules-central/main/data/fullset/database-v2/productCategory",
			PK:        "type",
			Database:  "database-v2",
			Container: "productCategory",
		}, */
		/* 				{
			URL:       "https://raw.githubusercontent.com/MicrosoftDocs/mslearn-cosmosdb-modules-central/main/data/fullset/database-v3/product",
			PK:        "categoryId",
			Database:  "database-v3",
			Container: "product",
		}, */
		/*
			{
				URL:       "https://raw.githubusercontent.com/MicrosoftDocs/mslearn-cosmosdb-modules-central/main/data/fullset/database-v3/productCategory",
				PK:        "type",
				Database:  "database-v3",
				Container: "productCategory",
			}, */
		/* 		{
			URL:       "https://raw.githubusercontent.com/MicrosoftDocs/mslearn-cosmosdb-modules-central/main/data/fullset/database-v4/customer",
			PK:        "customerId",
			Database:  "database-v4",
			Collection: "customer",
		}, */
		/* {
			URL:       "https://raw.githubusercontent.com/MicrosoftDocs/mslearn-cosmosdb-modules-central/main/data/fullset/database-v4/product",
			PK:        "categoryId",
			Database:  "database-v4",
			Container: "product",
		}, */
		/* 				{
				URL:       "https://raw.githubusercontent.com/MicrosoftDocs/mslearn-cosmosdb-modules-central/main/data/fullset/database-v4/productMeta",
				PK:        "type",
				Database:  "database-v4",
				Container: "productMeta",
			},
		}*/
		/* 		}
		   		for _, item := range imports {
		   			// deleteContainer
		   			// createContainer + handle errors...
		   			log.Printf("importing Container %s from URL %s", item.Collection, item.URL)
		   			//err := ImportJSON(client, item.URL, item.PK, item.Database, item.Collection)
		   			if err != nil {
		   				return err
		   			}
		   		} */
		case "x":
			fmt.Println("exiting...")
			break out
		default:
			return errors.New("command doesn't exist. exiting")
		}
	}
	return nil
}
func queryCustomer(client *mongo.Client, databaseName, collectionName, pk string) error {
	ctx := context.TODO()
	// returns handle to database
	database := client.Database(databaseName)
	// returns handle to collection
	collection := database.Collection(collectionName)
	//fmt.Printf("%v\n", collection.Name())
	log.Printf("%s\n", collectionName)
	cursor, err := collection.Find(ctx, bson.M{"_id": pk})
	if err != nil {
		return err
	}
	for cursor.Next(context.TODO()) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			return err
		}
		b, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", b)
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	//fmt.Printf("\nQuerying customer id: [%v] in %v\\%v\n", //response["_id"], databaseName, collectionName)

	requestStatistics := struct {
		ActivityID                                  string
		CommandName                                 string
		EstimateDelayFromRateLimitingInMilliseconds int
		RequestCharge                               float64
		RequestDurationInMilliseconds               int
		RetiredDueToRateLimiting                    bool
		OK                                          int
	}{}
	map1 := map[string]interface{}{"getLastRequestStatistics": "1"}
	err = database.RunCommand(ctx, map1, nil).Decode(&requestStatistics)
	if err != nil {
		return err
	}

	log.Printf("RequestCharge: %.2f\n", requestStatistics.RequestCharge)

	// b, err := json.MarshalIndent(requestStatistics, "", "    ")
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("%s\n", b)

	//log.Printf("Query page received with %d items. Status %d. ActivityId %s. Consuming %v RU\n", len(queryResponse.Items), queryResponse.RawResponse.StatusCode, queryResponse.ActivityID, queryResponse.RequestCharge)

	return nil
}

func GetLastRequestStats(client *mongo.Client, databaseName string) (*RequestStatistics, error) {
	ctx := context.TODO()
	// returns handle to database
	database := client.Database(databaseName)

	statistics := RequestStatistics{}
	map1 := map[string]interface{}{"getLastRequestStatistics": "1"}
	err := database.RunCommand(ctx, map1, nil).Decode(&statistics)
	if err != nil {
		return nil, err
	}
	return &statistics, nil

}

/* func ImportJSON(client *options.ClientOptions, url1, pk, databaseName, collectionName string) error {

	res, err := http.Get(url1)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	items := []map[string]interface{}{}
	err = json.Unmarshal(b, &items)
	if err != nil {
		return err
	}
	db := client.Database(databaseName).Collection(collectionName)
	if err != nil {
		return err
	}

	container, err := db.Collection(collectionName)
	if err != nil {
		return err
	}

	ctx := context.Background()

	ruSum := 0.0
	start := time.Now()

	for _, item := range items {
		// pretty print as we insert
		b, err := json.MarshalIndent(item, "", "    ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", b)

		// insert the item
		id, ok := item[pk]
		if !ok {
			return fmt.Errorf("item does not have member %s", pk)
		}
		val, ok := id.(string)
		if !ok {
			return fmt.Errorf("item member %s should be a string", pk)
		}
		//pk := azcosmos.NewPartitionKeyString(val)
		res, err := container.InsertOne(ctx, pk, b, nil)
		if err != nil {
			return err
		}
		ruSum = ruSum + float64(res.RequestCharge)
	}

	elapsed := time.Since(start)
	log.Printf("Total RUs consumed: %f in %f seconds\n", ruSum, elapsed.Seconds())

	return nil
} */

/*func queryCustomer(client options.ClientOptions, containerName, databaseName, partitionKey string) error {
	//Querying for a single customer
	pk := azcosmos.NewPartitionKeyString(partitionKey)
	_ = pk
	fmt.Printf("\nQuerying customer id: [%v] in %v\\%v\n", pk, databaseName, containerName)

 	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return err
	}

 	queryPager := container.NewQueryItemsPager("select * from customer c", pk, nil)

	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.Background())
		if err != nil {
			return err
		}
		for _, item := range queryResponse.Items {
			map1 := map[string]interface{}{}
			err := json.Unmarshal(item, &map1)
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(map1, "", "    ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", b)
		}
		log.Printf("Query page received with %d items. Status %d. ActivityId %s. Consuming %v RU\n", len(queryResponse.Items), queryResponse.RawResponse.StatusCode, queryResponse.ActivityID, queryResponse.RequestCharge)
	}
	return nil
}*/

func getCustomer(client *mongo.Client, databaseName, collectionName, pk string) (*CustomerV1, error) {
	ctx := context.TODO()
	// returns handle to database
	database := client.Database(databaseName)
	// returns handle to collection
	collection := database.Collection(collectionName)
	//fmt.Printf("%v\n", collection.Name())
	result := collection.FindOne(ctx, bson.M{"_id": pk})
	if err := result.Err(); err != nil {
		return nil, err
	}
	customer := CustomerV1{}
	err := result.Decode(&customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

/* func CreateDatabase(client options.ClientOptions) error {
	schemaVersionStart := 1
	schemaVersionEnd := 4
	schemaVersion := 0
	if !(schemaVersion == 0) {
		schemaVersionStart = schemaVersion
		schemaVersionEnd = schemaVersion
	} else {
		schemaVersionStart = 1
		schemaVersionEnd = 4
	}
	for schemaVersionCounter := schemaVersionStart; schemaVersionCounter <= schemaVersionEnd; schemaVersionCounter++ {
		fmt.Printf("Create started for schema %v\n", schemaVersionCounter)
		err := CreateDatabaseAndContainers(client, "database-v"+strconv.Itoa(schemaVersionCounter), schemaVersionCounter)
		if err != nil {
			return err
		}
	}
	return nil
} */

func CreateDatabaseAndContainers(client options.ClientOptions, databaseName string, schema int) error {
	// if schema >= 1 && schema <= 4 {
	// 	//throughput := azcosmos.NewManualThroughputProperties(400)
	// 	//databaseProperties := azcosmos.DatabaseProperties{ID: databaseName}
	// 	//databaseOptions := &azcosmos.CreateDatabaseOptions{}
	// 	//databaseResp, err := client.CreateDatabase(context.Background(), databaseProperties, databaseOptions)
	// 	if err != nil {
	// 		var responseErr *azcore.ResponseError
	// 		errors.As(err, &responseErr)
	// 		if responseErr.ErrorCode == "Conflict" {
	// 			log.Printf("Database [%v] already exists\n", databaseName)
	// 		} else {
	// 			return err
	// 		}
	// 	} else {
	// 		fmt.Printf("Database [%v] created. ActivityId %s\n", databaseName, databaseResp.ActivityID)
	// 	}
	// }
	return nil
}
