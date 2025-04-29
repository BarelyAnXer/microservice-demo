package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"    // Redis client
	"github.com/jackc/pgx/v4/pgxpool" // PostgreSQL driver and connection pool
)

var ctx = context.Background() // Use a background context for long-running operations

func main() {
	log.Println("Starting Gardener Bot worker...")

	// --- Get Configuration from Environment Variables ---
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisQueueName := os.Getenv("REDIS_QUEUE_NAME")

	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgDB := os.Getenv("POSTGRES_DB")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")

	// Basic validation
	if redisHost == "" || redisPort == "" || redisQueueName == "" || pgHost == "" || pgPort == "" || pgDB == "" || pgUser == "" || pgPassword == "" {
		log.Fatal("Missing required environment variables for Redis or PostgreSQL connection.")
	}

	log.Printf("Connecting to Redis at %s:%s", redisHost, redisPort)
	log.Printf("Listening on Redis queue: %s", redisQueueName)
	log.Printf("Connecting to PostgreSQL at %s:%s, Database: %s", pgHost, pgPort, pgDB)

	// --- Connect to Redis ---
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   0, // use default DB
	})

	// Ping Redis to ensure connection is working
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis.")

	// --- Connect to PostgreSQL ---
	// Build the connection string
	pgConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pgUser, pgPassword, pgHost, pgPort, pgDB)

	// Use a connection pool
	dbpool, err := pgxpool.Connect(ctx, pgConnStr)
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL database: %v", err)
	}
	defer dbpool.Close() // Close the pool when the main function exits

	// Ping DB to ensure connection is working
	err = dbpool.Ping(ctx)
	if err != nil {
		log.Fatalf("Could not ping PostgreSQL database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL.")

	// --- Main Worker Loop ---
	log.Println("Gardener Bot is listening for actions...")
	for {
		// BLPOP blocks until an element is available in the queue
		// 0 means block indefinitely
		result, err := rdb.BLPop(ctx, 0, redisQueueName).Result()
		if err != nil {
			// Handle potential errors, but BLPOP blocking means errors are often connection issues
			log.Printf("Error from Redis BLPOP: %v. Attempting to reconnect...", err)
			time.Sleep(5 * time.Second) // Wait before retrying
			continue                    // Go to the next iteration of the loop
		}

		// result is a slice: [queue_name, element]
		if len(result) != 2 {
			log.Printf("Received unexpected result from BLPOP: %v", result)
			continue // Skip processing this unexpected message
		}

		actionType := result[1] // The actual message (e.g., "water", "sunshine")
		log.Printf("Received action from queue: %s", actionType)

		// --- Process the Action (Update DB) ---
		err = processAction(dbpool, actionType)
		if err != nil {
			log.Printf("Error processing action '%s': %v", actionType, err)
			// In a real application, you might want to log this error and potentially
			// move the message to a dead-letter queue in Redis.
		} else {
			log.Printf("Successfully processed action '%s'.", actionType)
		}
	}
}

// processAction updates the database based on the received action type
func processAction(dbpool *pgxpool.Pool, actionType string) error {
	// Start a transaction
	tx, err := dbpool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	// Defer a rollback in case anything goes wrong
	defer tx.Rollback(ctx) // Rollback is a no-op if the transaction is committed

	// SQL to increment the count for a specific metric_name
	// This uses INSERT ... ON CONFLICT (UPSERT) which is a clean way
	// to handle both inserting the first count and incrementing subsequent counts.
	sqlStatement := `
		INSERT INTO garden_state (metric_name, count)
		VALUES ($1, 1)
		ON CONFLICT (metric_name)
		DO UPDATE SET count = garden_state.count + 1;
	`

	_, err = tx.Exec(ctx, sqlStatement, actionType)
	if err != nil {
		return fmt.Errorf("unable to execute UPSERT statement: %w", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil // Success
}
