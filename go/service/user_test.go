package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"crud_app/dto"
	mock_repository "crud_app/repository/mocks_repository"
	mock_service "crud_app/service/mocks_service"
)

const id uint = 1

var (
	testUsers = []dto.User{
		{ID: 1, Name: "John", Age: 10},
		{ID: 2, Name: "Jane", Age: 10},
	}
	testUser = &dto.User{
		Name: "John",
		Age:  10,
	}
	testUserWithID = &dto.User{
		ID:   1,
		Name: "John",
		Age:  10,
	}
	testUserNameEmpty = &dto.User{
		Name: "",
		Age:  10,
	}
	testUserNameShort = &dto.User{
		Name: "Jo",
		Age:  10,
	}
	testUserNameLong = &dto.User{
		Name: strings.Repeat("J", 101),
		Age:  10,
	}
	testUserAgeNeg = &dto.User{
		Name: "John",
		Age:  0,
	}
	testUserAgeUnreal = &dto.User{
		Name: "John",
		Age:  151,
	}
)

var (
	errRepo         = errors.New("repository error")
	errNameEmpty    = errors.New("name is required")
	errNameTooShort = errors.New("name must be at least 2 characters long")
	errNameTooLong  = errors.New("name cannot exceed 100 characters")
	errAgeNeg       = errors.New("age must be positive")
	errAgeUnreal    = errors.New("age seems unrealistic")
	errUserNil      = errors.New("user object cannot be nil")
	errUserNotFound = errors.New("user with ID not found")
)

func TestUser_List(t *testing.T) {
	type testCase struct {
		name          string
		setupMocks    func(*mock_repository.MockUserRepo)
		expectedUsers []dto.User
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name: "successful list with users",
			setupMocks: func(mockRepo *mock_repository.MockUserRepo) {
				expectedUsers := testUsers
				mockRepo.EXPECT().
					List(gomock.Any()).
					Return(expectedUsers, nil)
			},
			expectedUsers: testUsers,
			wantError:     false,
			expectedError: nil,
		}, {
			name: "error repository list",
			setupMocks: func(mr *mock_repository.MockUserRepo) {
				expectedError := errRepo
				mr.EXPECT().
					List(gomock.Any()).
					Return(nil, expectedError)
			},
			expectedUsers: nil,
			wantError:     true,
			expectedError: errRepo,
		}, {
			name: "empty result",
			setupMocks: func(mr *mock_repository.MockUserRepo) {
				mr.EXPECT().
					List(gomock.Any()).
					Return([]dto.User{}, nil)
			},
			expectedUsers: []dto.User{},
			wantError:     false,
			expectedError: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockValidator := mock_service.NewMockUserValidator(ctrl)
			mockRepo := mock_repository.NewMockUserRepo(ctrl)

			tc.setupMocks(mockRepo)

			service := NewUser(mockValidator, mockRepo)
			result, err := service.List(context.Background())

			if tc.wantError {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUsers, result)
			}
		})
	}
}

func TestUser_Create(t *testing.T) {
	type testCase struct {
		name          string
		input         *dto.User
		setupMocks    func(*mock_service.MockUserValidator, *mock_repository.MockUserRepo)
		expectedUser  *dto.User
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name:  "successful creation",
			input: testUser,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedUser := testUserWithID
				mockValidator.EXPECT().
					Create(gomock.Any(), testUser).
					Return(nil)
				mockRepo.EXPECT().
					Create(gomock.Any(), testUser).
					Return(expectedUser, nil)
			},
			expectedUser:  testUserWithID,
			wantError:     false,
			expectedError: nil,
		}, {
			name:  "error repository create",
			input: testUser,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				mockValidator.EXPECT().
					Create(gomock.Any(), testUser).
					Return(nil)
				expectedError := errRepo
				mockRepo.EXPECT().
					Create(gomock.Any(), testUser).
					Return(nil, expectedError)
			},
			expectedUser:  nil,
			wantError:     true,
			expectedError: errRepo,
		}, {
			name:  "error name is required",
			input: testUserNameEmpty,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errNameEmpty
				mockValidator.EXPECT().
					Create(gomock.Any(), testUserNameEmpty).
					Return(expectedError)
			},
			expectedUser:  nil,
			wantError:     true,
			expectedError: errNameEmpty,
		}, {
			name:  "error name must be at least 2 characters long",
			input: testUserNameShort,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errNameTooShort
				mockValidator.EXPECT().
					Create(gomock.Any(), testUserNameShort).
					Return(expectedError)
			},
			expectedUser:  nil,
			wantError:     true,
			expectedError: errNameTooShort,
		}, {
			name:  "error name cannot exceed 100 characters",
			input: testUserNameLong,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errNameTooLong
				mockValidator.EXPECT().
					Create(gomock.Any(), testUserNameLong).
					Return(expectedError)
			},
			expectedUser:  nil,
			wantError:     true,
			expectedError: errNameTooLong,
		}, {
			name:  "error age must be positive",
			input: testUserAgeNeg,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errAgeNeg
				mockValidator.EXPECT().
					Create(gomock.Any(), testUserAgeNeg).
					Return(expectedError)
			},
			expectedUser:  nil,
			wantError:     true,
			expectedError: errAgeNeg,
		}, {
			name:  "error age seems unrealistic",
			input: testUserAgeUnreal,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errAgeUnreal
				mockValidator.EXPECT().
					Create(gomock.Any(), testUserAgeUnreal).
					Return(expectedError)
			},
			expectedUser:  nil,
			wantError:     true,
			expectedError: errAgeUnreal,
		}, {
			name:  "error user object cannot be nil",
			input: nil,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errUserNil
				mockValidator.EXPECT().
					Create(gomock.Any(), nil).
					Return(expectedError)
			},
			expectedUser:  nil,
			wantError:     true,
			expectedError: errUserNil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockValidator := mock_service.NewMockUserValidator(ctrl)
			mockRepo := mock_repository.NewMockUserRepo(ctrl)

			tc.setupMocks(mockValidator, mockRepo)

			service := NewUser(mockValidator, mockRepo)
			result, err := service.Create(context.Background(), tc.input)

			if tc.wantError {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUser, result)
			}
		})
	}
}

func TestUser_Update(t *testing.T) {
	type testCase struct {
		name          string
		user          *dto.User
		id            uint
		setupMocks    func(*mock_service.MockUserValidator, *mock_repository.MockUserRepo)
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name: "successful updation",
			user: testUser,
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				mockValidator.EXPECT().
					Update(gomock.Any(), testUser, id).
					Return(nil)
				mockRepo.EXPECT().
					Update(gomock.Any(), testUser, id).
					Return(nil)
			},
			wantError:     false,
			expectedError: nil,
		}, {
			name: "error repository update",
			user: testUser,
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				mockValidator.EXPECT().
					Update(gomock.Any(), testUser, id).
					Return(nil)
				expectedError := errRepo
				mockRepo.EXPECT().
					Update(gomock.Any(), testUser, id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errRepo,
		}, {
			name: "error name is required",
			user: testUserNameEmpty,
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errNameEmpty
				mockValidator.EXPECT().
					Update(gomock.Any(), testUserNameEmpty, id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errNameEmpty,
		}, {
			name: "error name must be at least 2 characters long",
			user: testUserNameShort,
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errNameTooShort
				mockValidator.EXPECT().
					Update(gomock.Any(), testUserNameShort, id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errNameTooShort,
		}, {
			name: "error name cannot exceed 100 characters",
			user: testUserNameLong,
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errNameTooLong
				mockValidator.EXPECT().
					Update(gomock.Any(), testUserNameLong, id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errNameTooLong,
		}, {
			name: "error age must be positive",
			user: testUserAgeNeg,
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errAgeNeg
				mockValidator.EXPECT().
					Update(gomock.Any(), testUserAgeNeg, id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errAgeNeg,
		}, {
			name: "error age seems unrealistic",
			user: testUserAgeUnreal,
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errAgeUnreal
				mockValidator.EXPECT().
					Update(gomock.Any(), testUserAgeUnreal, id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errAgeUnreal,
		}, {
			name: "error user object cannot be nil",
			user: nil,
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errUserNil
				mockValidator.EXPECT().
					Update(gomock.Any(), nil, id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errUserNil,
		}, {
			name: "error repository exists",
			user: testUser,
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errRepo
				mockValidator.EXPECT().
					Update(gomock.Any(), testUser, id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errRepo,
		}, {
			name: "error user with ID not found",
			user: testUser,
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errUserNotFound
				mockValidator.EXPECT().
					Update(gomock.Any(), testUser, id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errUserNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockValidator := mock_service.NewMockUserValidator(ctrl)
			mockRepo := mock_repository.NewMockUserRepo(ctrl)

			tc.setupMocks(mockValidator, mockRepo)

			service := NewUser(mockValidator, mockRepo)
			err := service.Update(context.Background(), tc.user, tc.id)

			if tc.wantError {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
	type testCase struct {
		name          string
		id            uint
		setupMocks    func(*mock_service.MockUserValidator, *mock_repository.MockUserRepo)
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name: "successful deletion",
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				mockValidator.EXPECT().
					Delete(gomock.Any(), id).
					Return(nil)
				mockRepo.EXPECT().
					Delete(gomock.Any(), id).
					Return(nil)
			},
			wantError:     false,
			expectedError: nil,
		}, {
			name: "error repository delete",
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				mockValidator.EXPECT().
					Delete(gomock.Any(), id).
					Return(nil)
				expectedError := errRepo
				mockRepo.EXPECT().
					Delete(gomock.Any(), id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errRepo,
		}, {
			name: "error repository exists",
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errRepo
				mockValidator.EXPECT().
					Delete(gomock.Any(), id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errRepo,
		}, {
			name: "error user with ID not found",
			id:   id,
			setupMocks: func(mockValidator *mock_service.MockUserValidator, mockRepo *mock_repository.MockUserRepo) {
				expectedError := errUserNotFound
				mockValidator.EXPECT().
					Delete(gomock.Any(), id).
					Return(expectedError)
			},
			wantError:     true,
			expectedError: errUserNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockValidator := mock_service.NewMockUserValidator(ctrl)
			mockRepo := mock_repository.NewMockUserRepo(ctrl)

			tc.setupMocks(mockValidator, mockRepo)

			service := NewUser(mockValidator, mockRepo)
			err := service.Delete(context.Background(), tc.id)

			if tc.wantError {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
