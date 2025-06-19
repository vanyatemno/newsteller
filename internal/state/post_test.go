package state_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"newsteller/internal/models"
	"newsteller/internal/repositories"
	"newsteller/internal/state"
)

var dbClient *mongo.Client
var postCollection *mongo.Collection

const testDBName = "newsteller_test"
const postCollectionName = "posts"

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

	if err = pool.Retry(func() error {
		var err error
		dbClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
		if err != nil {
			return err
		}
		return dbClient.Ping(context.TODO(), nil)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	postCollection = dbClient.Database(testDBName).Collection(postCollectionName)

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func clearPostCollection(t *testing.T) {
	_, err := postCollection.DeleteMany(context.Background(), primitive.M{})
	require.NoError(t, err, "Failed to clear post collection")
}

func TestPostState_NewPostState(t *testing.T) {
	postState := state.NewPostState(postCollection)
	assert.NotNil(t, postState, "NewPostState should not return nil")
}

func TestPostState_InsertAndFindByID(t *testing.T) {
	clearPostCollection(t)
	postState := state.NewPostState(postCollection)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Millisecond)
	post := &models.Post{
		Title:     "Test Post 1",
		Content:   "Content of test post 1",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := postState.Insert(ctx, post)
	require.NoError(t, err, "Insert should not return an error")
	require.NotNil(t, post.ID, "Post ID should be set after insert")
	require.False(t, post.ID.IsZero(), "Post ID should not be zero after insert")

	foundPost, err := postState.FindByID(ctx, post.ID.Hex())
	require.NoError(t, err, "FindByID should not return an error")
	require.NotNil(t, foundPost, "FindByID should return a post")

	assert.Equal(t, post.ID, foundPost.ID)
	assert.Equal(t, "Test Post 1", foundPost.Title)
	assert.Equal(t, "Content of test post 1", foundPost.Content)
	assert.WithinDuration(t, post.CreatedAt, foundPost.CreatedAt, time.Millisecond, "CreatedAt should match")
	assert.WithinDuration(t, post.UpdatedAt, foundPost.UpdatedAt, time.Millisecond, "UpdatedAt should match")

	nonExistentID := primitive.NewObjectID().Hex()
	_, err = postState.FindByID(ctx, nonExistentID)
	assert.Error(t, err, "FindByID for non-existent post should return an error")
	assert.Equal(t, mongo.ErrNoDocuments, err, "Error should be mongo.ErrNoDocuments")
}

func TestPostState_FindAll(t *testing.T) {
	clearPostCollection(t)
	postState := state.NewPostState(postCollection)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Millisecond)
	post1 := &models.Post{Title: "Post A", Content: "Content A", CreatedAt: now.Add(-2 * time.Hour), UpdatedAt: now.Add(-2 * time.Hour)}
	post2 := &models.Post{Title: "Post B", Content: "Content B", CreatedAt: now.Add(-1 * time.Hour), UpdatedAt: now.Add(-1 * time.Hour)}

	err := postState.Insert(ctx, post1)
	require.NoError(t, err)
	err = postState.Insert(ctx, post2)
	require.NoError(t, err)

	allPosts, err := postState.FindAll(ctx)
	require.NoError(t, err, "FindAll should not return an error")
	require.Len(t, allPosts, 2, "FindAll should return 2 posts")

	titles := []string{allPosts[0].Title, allPosts[1].Title}
	assert.Contains(t, titles, "Post A")
	assert.Contains(t, titles, "Post B")
}

func TestPostState_FindPaginated(t *testing.T) {
	clearPostCollection(t)
	postState := state.NewPostState(postCollection)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Millisecond)
	for i := 0; i < 5; i++ {
		post := &models.Post{
			Title:     fmt.Sprintf("Paginated Post %d", i+1),
			Content:   fmt.Sprintf("Content %d", i+1),
			CreatedAt: now.Add(time.Duration(i) * time.Minute),
			UpdatedAt: now.Add(time.Duration(i) * time.Minute),
		}
		err := postState.Insert(ctx, post)
		require.NoError(t, err)
	}

	query1 := &repositories.PaginatedSearchQuery{Page: 1, Limit: 2, Keyword: ""}
	paginatedPosts1, total1, err1 := postState.FindPaginated(ctx, query1)
	require.NoError(t, err1, "FindPaginated (page 1) should not return an error")
	assert.Equal(t, int64(5), total1, "Total count should be 5")
	require.Len(t, paginatedPosts1, 2, "Should return 2 posts for page 1")
	assert.Equal(t, "Paginated Post 5", paginatedPosts1[0].Title) // Latest
	assert.Equal(t, "Paginated Post 4", paginatedPosts1[1].Title)

	query2 := &repositories.PaginatedSearchQuery{Page: 2, Limit: 2, Keyword: ""}
	paginatedPosts2, total2, err2 := postState.FindPaginated(ctx, query2)
	require.NoError(t, err2, "FindPaginated (page 2) should not return an error")
	assert.Equal(t, int64(5), total2, "Total count should be 5")
	require.Len(t, paginatedPosts2, 2, "Should return 2 posts for page 2")
	assert.Equal(t, "Paginated Post 3", paginatedPosts2[0].Title)
	assert.Equal(t, "Paginated Post 2", paginatedPosts2[1].Title)

	query3 := &repositories.PaginatedSearchQuery{Page: 3, Limit: 2, Keyword: ""}
	paginatedPosts3, total3, err3 := postState.FindPaginated(ctx, query3)
	require.NoError(t, err3, "FindPaginated (page 3) should not return an error")
	assert.Equal(t, int64(5), total3, "Total count should be 5")
	require.Len(t, paginatedPosts3, 1, "Should return 1 post for page 3")
	assert.Equal(t, "Paginated Post 1", paginatedPosts3[0].Title)

	query4 := &repositories.PaginatedSearchQuery{Page: 1, Limit: 2, Keyword: "Paginated Post 3"}
	paginatedPosts4, total4, err4 := postState.FindPaginated(ctx, query4)
	require.NoError(t, err4, "FindPaginated (with keyword) should not return an error")
	assert.Equal(t, int64(1), total4, "Total count should be 0 with keyword (due to repo issue)")
	require.Len(t, paginatedPosts4, 1, "Should return 0 posts for page 1 with keyword (due to repo issue)")
	assert.Equal(t, "Paginated Post 3", paginatedPosts4[0].Title)
}

func TestPostState_Update(t *testing.T) {
	clearPostCollection(t)
	postState := state.NewPostState(postCollection)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Millisecond)
	post := &models.Post{
		Title:     "Original Title",
		Content:   "Original Content",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := postState.Insert(ctx, post)
	require.NoError(t, err)
	require.NotEmpty(t, post.ID)

	updatedTime := now.Add(3 * time.Hour).Truncate(time.Millisecond)

	updateData := &models.Post{
		ID:        post.ID,
		Title:     "Updated Title",
		Content:   "Updated Content",
		CreatedAt: post.CreatedAt,
		UpdatedAt: updatedTime,
	}

	err = postState.Update(ctx, updateData)
	require.NoError(t, err, "Update should not return an error")

	foundPostState, err := postState.FindByID(ctx, updateData.ID.Hex())
	require.NoError(t, err)
	require.NotNil(t, foundPostState)
	assert.Equal(t, "Updated Title", foundPostState.Title)
	assert.Equal(t, "Updated Content", foundPostState.Content)
	assert.WithinDuration(t, updatedTime, foundPostState.UpdatedAt, time.Millisecond, "UpdatedAt in state cache should match")

	repo := repositories.NewPostRepository(postCollection)
	foundPostDB, err := repo.FindByID(ctx, updateData.ID.Hex())
	require.NoError(t, err)
	require.NotNil(t, foundPostDB)
	assert.Equal(t, "Updated Title", foundPostDB.Title)
	assert.Equal(t, "Updated Content", foundPostDB.Content)
	assert.WithinDuration(t, updatedTime, foundPostDB.UpdatedAt, time.Millisecond, "UpdatedAt in DB should match")
}

func TestPostState_Delete(t *testing.T) {
	clearPostCollection(t)
	postState := state.NewPostState(postCollection)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Millisecond)
	post := &models.Post{
		Title:     "To Be Deleted",
		Content:   "Delete me",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := postState.Insert(ctx, post)
	require.NoError(t, err)
	require.NotNil(t, post.ID)
	postIDHex := post.ID.Hex()

	_, err = postState.FindByID(ctx, postIDHex)
	require.NoError(t, err, "Post should be findable before delete")

	err = postState.Delete(ctx, postIDHex)
	require.NoError(t, err, "Delete should not return an error")

	_, err = postState.FindByID(ctx, postIDHex)
	assert.Error(t, err, "FindByID after delete should return an error")
	assert.Equal(t, mongo.ErrNoDocuments, err, "Error should be mongo.ErrNoDocuments after delete")

	repo := repositories.NewPostRepository(postCollection)
	_, err = repo.FindByID(ctx, postIDHex)
	assert.Error(t, err, "FindByID from repo after delete should return an error")
	assert.Equal(t, mongo.ErrNoDocuments, err, "Error from repo should be mongo.ErrNoDocuments after delete")
}

func TestPostState_Insert_Sorting(t *testing.T) {
	clearPostCollection(t)
	postState := state.NewPostState(postCollection).(*state.Post)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Millisecond)
	post1 := &models.Post{Title: "Post 1", CreatedAt: now.Add(2 * time.Hour)}
	post2 := &models.Post{Title: "Post 2", CreatedAt: now}
	post3 := &models.Post{Title: "Post 3", CreatedAt: now.Add(1 * time.Hour)}

	// Insert in non-sorted order
	err := postState.Insert(ctx, post1)
	require.NoError(t, err)
	err = postState.Insert(ctx, post2)
	require.NoError(t, err)
	err = postState.Insert(ctx, post3)
	require.NoError(t, err)

	allPosts, err := postState.FindAll(ctx)
	require.NoError(t, err)
	require.Len(t, allPosts, 3)

	query := &repositories.PaginatedSearchQuery{Page: 1, Limit: 3, Keyword: ""}
	sortedPostsFromDB, _, err := postState.FindPaginated(ctx, query)
	require.NoError(t, err)
	require.Len(t, sortedPostsFromDB, 3)

	assert.Equal(t, "Post 1", sortedPostsFromDB[0].Title)
	assert.Equal(t, "Post 3", sortedPostsFromDB[1].Title)
	assert.Equal(t, "Post 2", sortedPostsFromDB[2].Title)
}
