package repositories

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"newsteller/internal/models"
)

var dbClient *mongo.Client
var postRepo *Post
var collection *mongo.Collection

const testCollectionName = "posts_test"

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "latest",
		Env: []string{
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=password",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	mongoURI := fmt.Sprintf("mongodb://root:password@%s", resource.GetHostPort("27017/tcp"))

	if err := pool.Retry(func() error {
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		dbClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			return err
		}
		return dbClient.Ping(ctx, nil)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	collection = dbClient.Database("newsteller_test").Collection(testCollectionName)
	postRepo = NewPostRepository(collection)

	code := m.Run()

	if err := dbClient.Disconnect(context.Background()); err != nil {
		log.Printf("Could not disconnect from mongo: %s", err)
	}

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func clearCollection(ctx context.Context) {
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to clear collection: %v", err)
	}
}

func TestNewPostRepository(t *testing.T) {
	assert.NotNil(t, postRepo)
	assert.Equal(t, collection, postRepo.c)
}

func TestPost_Create(t *testing.T) {
	ctx := context.Background()
	defer clearCollection(ctx)

	t.Run("Positive: Create a new post", func(t *testing.T) {
		post := &models.Post{
			Title:     "Test Post 1",
			Content:   "Content for test post 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		id, err := postRepo.Create(ctx, post)
		require.NoError(t, err)
		require.NotNil(t, id)
		assert.False(t, id.IsZero())

		var createdPost models.Post
		err = collection.FindOne(ctx, bson.M{"_id": *id}).Decode(&createdPost)
		require.NoError(t, err)
		assert.Equal(t, post.Title, createdPost.Title)
		assert.Equal(t, post.Content, createdPost.Content)
	})
}

func TestPost_All(t *testing.T) {
	ctx := context.Background()
	defer clearCollection(ctx)

	t.Run("Positive: Retrieve all posts", func(t *testing.T) {
		post1 := models.Post{ID: primitive.NewObjectID(), Title: "Post A", Content: "Content A", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		post2 := models.Post{ID: primitive.NewObjectID(), Title: "Post B", Content: "Content B", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		_, err := collection.InsertMany(ctx, []interface{}{post1, post2})
		require.NoError(t, err)

		posts, err := postRepo.All(ctx)
		require.NoError(t, err)
		assert.Len(t, posts, 2)
	})

	t.Run("Positive: Retrieve empty list when no posts", func(t *testing.T) {
		clearCollection(ctx) // Ensure collection is empty
		posts, err := postRepo.All(ctx)
		require.NoError(t, err)
		assert.Len(t, posts, 0)
	})
}

func TestPost_FindByID(t *testing.T) {
	ctx := context.Background()
	defer clearCollection(ctx)

	t.Run("Positive: Find an existing post", func(t *testing.T) {
		post := models.Post{
			Title:     "Find Me",
			Content:   "Content to find",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		res, err := collection.InsertOne(ctx, post)
		require.NoError(t, err)
		insertedID := res.InsertedID.(primitive.ObjectID)

		foundPost, err := postRepo.FindByID(ctx, insertedID.Hex())
		require.NoError(t, err)
		require.NotNil(t, foundPost)
		assert.Equal(t, insertedID, foundPost.ID)
		assert.Equal(t, "Find Me", foundPost.Title)
	})

	t.Run("Negative: Post not found", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		foundPost, err := postRepo.FindByID(ctx, nonExistentID)
		require.Error(t, err)
		assert.Nil(t, foundPost)
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})

	t.Run("Negative: Invalid ID format", func(t *testing.T) {
		invalidID := "this-is-not-an-object-id"
		foundPost, err := postRepo.FindByID(ctx, invalidID)
		require.Error(t, err)
		assert.Nil(t, foundPost)
		// Check for the specific error from primitive.ObjectIDFromHex
		_, idErr := primitive.ObjectIDFromHex(invalidID)
		assert.EqualError(t, err, idErr.Error())
	})
}

func TestPost_Update(t *testing.T) {
	ctx := context.Background()
	defer clearCollection(ctx)

	t.Run("Positive: Update an existing post", func(t *testing.T) {
		originalPost := models.Post{
			Title:     "Original Title",
			Content:   "Original Content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		res, err := collection.InsertOne(ctx, originalPost)
		require.NoError(t, err)
		insertedID := res.InsertedID.(primitive.ObjectID)

		updatedPostData := &models.Post{
			ID:        insertedID,
			Title:     "Updated Title",
			Content:   "Updated Content",
			UpdatedAt: time.Now().Add(time.Hour), // Ensure UpdatedAt is different
		}

		err = postRepo.Update(ctx, updatedPostData)
		require.NoError(t, err)

		var fetchedPost models.Post
		err = collection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&fetchedPost)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", fetchedPost.Title)
		assert.Equal(t, "Updated Content", fetchedPost.Content)
		// Check if UpdatedAt was actually updated (precision might be an issue with direct comparison)
		assert.True(t, fetchedPost.UpdatedAt.After(originalPost.UpdatedAt))
	})

	t.Run("Negative: Post not found for update", func(t *testing.T) {
		nonExistentPost := &models.Post{
			ID:        primitive.NewObjectID(),
			Title:     "Non Existent",
			Content:   "This post does not exist",
			UpdatedAt: time.Now(),
		}
		err := postRepo.Update(ctx, nonExistentPost)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no post found with ID")
	})

	t.Run("Negative: Invalid ID format for update", func(t *testing.T) {
		// This case is tricky because ObjectIDFromHex is called on post.ID.Hex()
		// If post.ID is not a valid ObjectID, Hex() might panic or return an empty string.
		// Let's assume post.ID is somehow malformed or zero, leading to an invalid hex.
		// A zero ObjectID will result in "000000000000000000000000" which is valid.
		// The function primitive.ObjectIDFromHex(post.ID.Hex()) will likely not fail if post.ID is a primitive.ObjectID.
		// The error for invalid ID format is primarily caught by FindByID or Delete.
		// For Update, the ID is taken from an existing models.Post struct.
		// If we want to test the ObjectIDFromHex failure within Update, we'd need to manipulate post.ID.Hex() to be invalid,
		// which is not straightforward with the current struct.
		// The more likely scenario for an ID error in Update is if the ID itself is valid hex but doesn't exist.
		// The current implementation of Update calls primitive.ObjectIDFromHex(post.ID.Hex()).
		// If post.ID is already a primitive.ObjectID, post.ID.Hex() will return its hex string.
		// primitive.ObjectIDFromHex will then convert it back. This path should not fail for format.
		// The test for "Post not found for update" covers the case where the ID is valid format but not in DB.
		// Let's consider if the input `post *models.Post` could have an `ID` that is not a valid `primitive.ObjectID`.
		// The `models.Post` struct defines `ID primitive.ObjectID`. So it will always be a valid `primitive.ObjectID` or zero.
		// The `Hex()` method on a zero `primitive.ObjectID` returns "000000000000000000000000".
		// `primitive.ObjectIDFromHex("000000000000000000000000")` is valid.
		// So, the ObjectIDFromHex call in Update is unlikely to fail due to format if the input is a valid models.Post.
		// The primary failure point for ID in Update is "not found".
		// If the intent is to test the internal ObjectIDFromHex, it's better tested in FindByID/Delete.
		// For completeness, if we assume a scenario where post.ID.Hex() could return an invalid string (e.g., if models.Post.ID was string type):
		// This test is more conceptual for Update as structured.
		// For now, we'll skip a direct "invalid ID format" for Update as the structure makes it hard to trigger that specific error path.
		// The "post not found" covers the practical ID-related error for Update.
		t.Skip("Skipping invalid ID format for Update as models.Post.ID is primitive.ObjectID, making Hex() always produce valid hex.")
	})
}

func TestPost_Delete(t *testing.T) {
	ctx := context.Background()
	defer clearCollection(ctx)

	t.Run("Positive: Delete an existing post", func(t *testing.T) {
		post := models.Post{
			Title:     "To Be Deleted",
			Content:   "Delete this content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		res, err := collection.InsertOne(ctx, post)
		require.NoError(t, err)
		insertedID := res.InsertedID.(primitive.ObjectID)

		err = postRepo.Delete(ctx, insertedID.Hex())
		require.NoError(t, err)

		// Verify it's deleted
		var deletedPost models.Post
		err = collection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&deletedPost)
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})

	t.Run("Negative: Delete a non-existent post", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		err := postRepo.Delete(ctx, nonExistentID)
		// DeleteOne doesn't return an error if no document matches.
		// It returns a result with DeletedCount: 0. The repository method also doesn't return an error.
		// This behavior might need adjustment in the repository if an error is desired for "not found".
		// For now, test current behavior: no error.
		require.NoError(t, err)
	})

	t.Run("Negative: Invalid ID format for delete", func(t *testing.T) {
		invalidID := "not-a-valid-hex-id"
		err := postRepo.Delete(ctx, invalidID)
		require.Error(t, err)
		_, idErr := primitive.ObjectIDFromHex(invalidID)
		assert.EqualError(t, err, idErr.Error())
	})
}

func TestPost_FindPaginated(t *testing.T) {
	ctx := context.Background()
	defer clearCollection(ctx)

	// Seed data
	var postsToSeed []interface{}
	for i := 0; i < 25; i++ {
		postsToSeed = append(postsToSeed, models.Post{
			ID:        primitive.NewObjectID(),
			Title:     fmt.Sprintf("Paginated Post %d", i+1),
			Content:   fmt.Sprintf("Content for paginated post %d. Some have 'keyword'.", i+1),
			CreatedAt: time.Now().Add(time.Duration(i) * time.Minute), // Vary CreatedAt for sorting
			UpdatedAt: time.Now(),
		})
	}
	_, err := collection.InsertMany(ctx, postsToSeed)
	require.NoError(t, err)

	// Add specific posts for keyword search
	keywordPost1 := models.Post{ID: primitive.NewObjectID(), Title: "Special Keyword Post", Content: "This has the keyword.", CreatedAt: time.Now().Add(100 * time.Minute), UpdatedAt: time.Now()}
	keywordPost2 := models.Post{ID: primitive.NewObjectID(), Title: "Another Post", Content: "This also contains the special keyword.", CreatedAt: time.Now().Add(101 * time.Minute), UpdatedAt: time.Now()}
	_, err = collection.InsertMany(ctx, []interface{}{keywordPost1, keywordPost2})
	require.NoError(t, err)

	t.Run("Positive: Retrieve first page, no keyword", func(t *testing.T) {
		query := &PaginatedSearchQuery{Page: 1, Limit: 10}
		posts, total, err := postRepo.FindPaginated(ctx, query)
		require.NoError(t, err)
		assert.Len(t, posts, 10)
		assert.EqualValues(t, 25+2, total) // 25 generic + 2 keyword posts
		// Check sorting (newest first)
		assert.True(t, posts[0].CreatedAt.After(posts[1].CreatedAt))
	})

	t.Run("Positive: Retrieve second page, no keyword", func(t *testing.T) {
		query := &PaginatedSearchQuery{Page: 2, Limit: 10}
		posts, total, err := postRepo.FindPaginated(ctx, query)
		require.NoError(t, err)
		assert.Len(t, posts, 10)
		assert.EqualValues(t, 27, total)
	})

	t.Run("Positive: Retrieve last page (partial), no keyword", func(t *testing.T) {
		query := &PaginatedSearchQuery{Page: 3, Limit: 10} // 27 total, 10 per page. Page 3 should have 7.
		posts, total, err := postRepo.FindPaginated(ctx, query)
		require.NoError(t, err)
		assert.Len(t, posts, 7)
		assert.EqualValues(t, 27, total)
	})

	t.Run("Positive: Retrieve with keyword", func(t *testing.T) {
		query := &PaginatedSearchQuery{Page: 1, Limit: 5, Keyword: "special keyword"}
		posts, total, err := postRepo.FindPaginated(ctx, query)
		require.NoError(t, err)
		assert.Len(t, posts, 2) // Only 2 posts match "special keyword"
		assert.EqualValues(t, 2, total)
		assert.Contains(
			t,
			strings.ToLower(posts[0].Content+posts[0].Title),
			"special keyword",
			"Content of post 0 does not contain keyword",
		)
		assert.Contains(
			t,
			strings.ToLower(posts[1].Content+posts[1].Title),
			"special keyword",
			"Content of post 1 does not contain keyword",
		)
	})

	t.Run("Positive: Retrieve with keyword, no results", func(t *testing.T) {
		query := &PaginatedSearchQuery{Page: 1, Limit: 5, Keyword: "nonexistentkeyword123"}
		posts, total, err := postRepo.FindPaginated(ctx, query)
		require.NoError(t, err)
		assert.Len(t, posts, 0)
		assert.EqualValues(t, 0, total)
	})

	t.Run("Positive: Page beyond total results, no keyword", func(t *testing.T) {
		query := &PaginatedSearchQuery{Page: 10, Limit: 10} // Total 27, page 10 should be empty
		posts, total, err := postRepo.FindPaginated(ctx, query)
		require.NoError(t, err)
		assert.Len(t, posts, 0)
		assert.EqualValues(t, 27, total)
	})
}
