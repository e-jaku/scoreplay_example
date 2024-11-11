package repository

import (
	"context"
	"scoreplay/internal/config"
	"scoreplay/internal/db"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestContainer(ctx context.Context, t *testing.T) (testcontainers.Container, *config.DBConfig) {
	// config PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * time.Second),
	}

	// start the container
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	require.NoError(t, err)

	// get container host and port
	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432/tcp")

	// setup config based on host and port
	dbConfig := config.DBConfig{
		Host:       host,
		Port:       port.Int(),
		User:       "testuser",
		Password:   "testpassword",
		DBName:     "testdb",
		SSLMode:    "disable",
		Migrations: "../../migrations",
	}

	return postgresC, &dbConfig
}

// TestCRUD is an integration test over the CRUD functionality of both Repositories with a postgres test container
func TestCRUD(t *testing.T) {
	ctx := context.Background()

	testC, dbConfig := setupTestContainer(ctx, t)
	defer testC.Terminate(ctx)

	testDB, err := db.NewPostgresConnection(dbConfig)
	require.NoError(t, err)
	defer testDB.Close()

	mediaRepo := NewPostgresMediaRepository(testDB)
	tagRepo := NewPostgresTagRepository(testDB)

	// Invalid Tag id
	_, err = mediaRepo.ListMediaByTagId(ctx, "non-existing")
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid input syntax for type uuid")

	// Unknown tag id
	medias, err := mediaRepo.ListMediaByTagId(ctx, uuid.NewString())
	require.NoError(t, err)
	require.Empty(t, medias)

	// No tags in db yet
	tags, err := tagRepo.ListTags(ctx)
	require.NoError(t, err)
	require.Empty(t, tags)

	// Add 1 tag
	tag, err := tagRepo.CreateTag(ctx, "Some tag")
	require.NoError(t, err)
	require.Equal(t, "Some tag", tag.Name)

	//Create media with tag
	media, err := mediaRepo.CreateMedia(ctx, "Some fancy picture", []string{tag.ID.String()}, "test")
	require.NoError(t, err)
	require.Equal(t, "Some fancy picture", media.Name)

	// List tags should contain 1
	tags, err = tagRepo.ListTags(ctx)
	require.NoError(t, err)
	require.Len(t, tags, 1)

	// Add a second tag
	tag2, err := tagRepo.CreateTag(ctx, "Some other tag")
	require.NoError(t, err)

	// Create a second media with both tags
	media2, err := mediaRepo.CreateMedia(ctx, "Some other fancy picture", []string{tag2.ID.String(), tag.ID.String()}, "test")
	require.NoError(t, err)

	// List based on tag 1 should be 2 medias
	medias, err = mediaRepo.ListMediaByTagId(ctx, tag.ID.String())
	require.NoError(t, err)
	require.Len(t, medias, 2)

	// List based on tag 2 should be only the second media
	medias, err = mediaRepo.ListMediaByTagId(ctx, tag2.ID.String())
	require.NoError(t, err)
	require.Len(t, medias, 1)
	require.Equal(t, media2.Name, medias[0].Name)

	// Get the tags provided in the array
	tags, err = tagRepo.GetTags(ctx, []string{tag.ID.String()})
	require.NoError(t, err)
	require.Equal(t, tag.Name, tags[0].Name)

	// 1 tag exist the other does not, should throw an error
	_, err = tagRepo.GetTags(ctx, []string{tag.ID.String(), uuid.NewString()})
	require.Error(t, err)
}
