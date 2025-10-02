package db

import (
	"context"
	"testing"
	"time"

	"github.com/ShubhKanodia/GoBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)   //will fail if there is an error
	require.NotEmpty(t, user) //will fail if the user is empty

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	//will fail if the id is zero(id needs to be auto created by the database)
	require.True(t, user.PasswordChangedAt.IsZero()) //filled with default val of zero timestamp
	require.NotZero(t, user.CreatedAt)               //will fail if the created at is zero(created at needs to be auto created by the database)

	return user

}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)    //will fail if there is an error
	require.NotEmpty(t, user2) //will fail if the user is empty
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordChangedAt, user2.PasswordChangedAt)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second) //will fail if the created at is not within 1 second
}
