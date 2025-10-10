package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"crud_app/dto"
)

const id uint = 1

var (
	testUsersSort = []dto.User{
		{ID: 1, Name: "Jane", Age: 11},
		{ID: 2, Name: "John", Age: 10},
	}
	testUsersNoSort = []dto.User{
		{ID: 2, Name: "John", Age: 10},
		{ID: 1, Name: "Jane", Age: 11},
	}
	testUsersAfterDelete = []dto.User{
		{ID: 2, Name: "John", Age: 10},
	}
	testUser = &dto.User{
		Name: "John",
		Age:  10,
	}
	testUserWithID = &dto.User{
		ID:   id,
		Name: "John",
		Age:  10,
	}
	testNewUser = &dto.User{
		Name: "Jane",
		Age:  11,
	}
)

var (
	errDB = errors.New("sql: database is closed")
)

var (
	dsn = "host=localhost user=user password=password dbname=mydb port=5432 sslmode=disable"
)

func TestUserRepo_List(t *testing.T) {
	type testCase struct {
		name          string
		setupDB       func(*gorm.DB)
		expectedUsers []dto.User
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name: "successful list",
			setupDB: func(db *gorm.DB) {
				for _, user := range testUsersSort {
					db.Create(&dto.User{
						Name: user.Name,
						Age:  user.Age,
					})
				}
			},
			expectedUsers: testUsersSort,
		}, {
			name: "successful list with order by",
			setupDB: func(db *gorm.DB) {
				for _, user := range testUsersNoSort {
					db.Create(&dto.User{
						ID:   user.ID,
						Name: user.Name,
						Age:  user.Age,
					})
				}
			},
			expectedUsers: testUsersSort,
		}, {
			name: "successful list with where deleted_at",
			setupDB: func(db *gorm.DB) {
				for _, user := range testUsersSort {
					db.Create(&dto.User{
						ID:   user.ID,
						Name: user.Name,
						Age:  user.Age,
					})
				}
				db.Delete(&dto.User{}, id)
			},
			expectedUsers: testUsersAfterDelete,
		}, {
			name: "error database is closed",
			setupDB: func(db *gorm.DB) {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			},
			wantError:     true,
			expectedError: errDB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			require.NoError(t, err)

			resetDB(db)
			if tc.setupDB != nil {
				tc.setupDB(db)
			}

			repo := NewUserRepo(db)
			resultRepo, errRepo := repo.List(context.Background())

			if tc.wantError {
				require.Error(t, errRepo)
				require.Equal(t, tc.expectedError, errRepo)
				require.Nil(t, resultRepo)

				return
			}

			var result []dto.User
			query := `
				SELECT id, name, age, created_at, updated_at
				FROM users
				WHERE deleted_at IS NULL
				ORDER BY id
			`
			err = db.Raw(query).Scan(&result).Error

			require.NoError(t, errRepo)
			require.NoError(t, err)
			require.NotNil(t, result, resultRepo)

			for i, expected := range tc.expectedUsers {
				resRepo := resultRepo[i]
				res := result[i]

				require.Equal(t, resRepo.ID, res.ID)
				require.Equal(t, expected.Name, resRepo.Name, res.Name)
				require.Equal(t, expected.Age, resRepo.Age, res.Age)
			}
		})
	}
}

func TestUserRepo_Create(t *testing.T) {
	type testCase struct {
		name          string
		inputUser     *dto.User
		setupDB       func(*gorm.DB)
		expectedUser  *dto.User
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name:      "successful creation",
			inputUser: testUser,
			expectedUser:  testUserWithID,
		}, {
			name:      "error database is closed",
			inputUser: testUser,
			setupDB: func(db *gorm.DB) {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			},
			wantError:     true,
			expectedError: errDB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			require.NoError(t, err)

			resetDB(db)
			if tc.setupDB != nil {
				tc.setupDB(db)
			}

			repo := NewUserRepo(db)
			resultRepo, errRepo := repo.Create(context.Background(), tc.inputUser)

			if tc.wantError {
				require.Error(t, errRepo)
				require.Equal(t, tc.expectedError, errRepo)
				require.Nil(t, resultRepo)

				return
			}

			var result dto.User
			query := `
				SELECT id, name, age, created_at, updated_at
				FROM users
				WHERE id = $1
			`
			err = db.Raw(query, resultRepo.ID).Scan(&result).Error

			require.NoError(t, errRepo)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, tc.expectedUser.Name, result.Name, resultRepo.Name)
			require.Equal(t, tc.expectedUser.Age, result.Age, resultRepo.Age)
			require.Equal(t, tc.expectedUser.ID, resultRepo.ID)
			require.Greater(t, result.ID, uint(0))
		})
	}
}

func TestUserRepo_Update(t *testing.T) {
	type testCase struct {
		name          string
		inputUser     *dto.User
		id            uint
		setupDB       func(*gorm.DB)
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name:      "successful updation",
			inputUser: testNewUser,
			id:        testUserWithID.ID,
			setupDB: func(db *gorm.DB) {
				db.Create(&dto.User{
					ID:   testUserWithID.ID,
					Name: testUserWithID.Name,
					Age:  testUserWithID.Age,
				})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			require.NoError(t, err)

			resetDB(db)
			if tc.setupDB != nil {
				tc.setupDB(db)
			}

			repo := NewUserRepo(db)
			errRepo := repo.Update(context.Background(), tc.inputUser, tc.id)

			if tc.wantError {
				require.Error(t, errRepo)
				require.Equal(t, tc.expectedError, errRepo)

				return
			}

			var result dto.User
			query := `
				SELECT id, name, age, created_at, updated_at
				FROM users
				WHERE id = $1
			`
			err = db.Raw(query, tc.id).Scan(&result).Error

			require.NoError(t, errRepo)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, tc.inputUser.Name, result.Name)
			require.Equal(t, tc.inputUser.Age, result.Age)
		})
	}
}

func TestUserRepo_Delete(t *testing.T) {
	type testCase struct {
		name          string
		id            uint
		setupDB       func(*gorm.DB)
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name: "successful deletion",
			id:   testUserWithID.ID,
			setupDB: func(db *gorm.DB) {
				db.Create(&dto.User{
					ID:   testUserWithID.ID,
					Name: testUserWithID.Name,
					Age:  testUserWithID.Age,
				})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			require.NoError(t, err)

			resetDB(db)
			if tc.setupDB != nil {
				tc.setupDB(db)
			}

			repo := NewUserRepo(db)
			errRepo := repo.Delete(context.Background(), tc.id)

			if tc.wantError {
				require.Error(t, errRepo)
				require.Equal(t, tc.expectedError, errRepo)

				return
			}

			var result dto.User
			query := `
				SELECT id, name, age, created_at, updated_at, deleted_at
				FROM users
				WHERE id = $1
			`
			err = db.Raw(query, tc.id).Scan(&result).Error

			require.NoError(t, errRepo)
			require.NoError(t, err)
			require.True(t, result.DeletedAt.Valid, "Record should be soft deleted")
			require.NotNil(t, result.DeletedAt.Time)
		})
	}
}

func TestUserRepo_Exists(t *testing.T) {
	type testCase struct {
		name          string
		id            uint
		setupDB       func(*gorm.DB)
		expectedValue bool
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name: "successful check exists",
			id:   testUserWithID.ID,
			setupDB: func(db *gorm.DB) {
				db.Create(&dto.User{
					ID:   testUserWithID.ID,
					Name: testUserWithID.Name,
					Age:  testUserWithID.Age,
				})
			},
			expectedValue: true,
		}, {
			name:          "successful check not exists",
			id:            testUserWithID.ID,
		}, {
			name: "successful check not exists user deleted",
			id:   testUserWithID.ID,
			setupDB: func(db *gorm.DB) {
				db.Create(&dto.User{
					ID:   testUserWithID.ID,
					Name: testUserWithID.Name,
					Age:  testUserWithID.Age,
				})
				db.Delete(&dto.User{}, testUserWithID.ID)
			},
		}, {
			name: "error database is closed",
			id:   testUserWithID.ID,
			setupDB: func(db *gorm.DB) {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			},
			wantError:     true,
			expectedError: errDB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			require.NoError(t, err)

			resetDB(db)
			if tc.setupDB != nil {
				tc.setupDB(db)
			}

			repo := NewUserRepo(db)
			resultRepo, errRepo := repo.Exists(context.Background(), tc.id)

			if tc.wantError {
				require.Error(t, errRepo)
				require.Equal(t, tc.expectedError, errRepo)
				require.False(t, resultRepo)
				require.Equal(t, tc.expectedValue, resultRepo)

				return
			}

			var count int64
			query := `
				SELECT COUNT(*)
				FROM users
				WHERE id = $1
				AND deleted_at IS NULL
			`
			err = db.Raw(query, tc.id).Scan(&count).Error
			result := count > 0

			require.NoError(t, errRepo)
			require.NoError(t, err)
			require.Equal(t, tc.expectedValue, resultRepo, result)
			if tc.expectedValue {
				require.True(t, resultRepo)
			} else {
				require.False(t, resultRepo)
			}
		})
	}
}

func resetDB(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE users")
	db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
}
