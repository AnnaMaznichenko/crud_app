package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"crud_app/dto"
	mock_repository "crud_app/repository/mocks_repository"
)

var (
	errFailedToCheck = fmt.Errorf("failed to check user existence: %w", errRepo)
)

var (
	testUserNameShortWithSpaces = &dto.User{
		Name: "  J  ",
		Age:  10,
	}
)

func TestUserValidator_Create(t *testing.T) {
	type testCase struct {
		name          string
		input         *dto.User
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name:          "successful creation",
			input:         testUser,
		}, {
			name:          "error user object cannot be nil",
			wantError:     true,
			expectedError: errUserNil,
		}, {
			name:          "error name is required",
			input:         testUserNameEmpty,
			wantError:     true,
			expectedError: errNameEmpty,
		}, {
			name:          "error name must be at least 2 characters long",
			input:         testUserNameShort,
			wantError:     true,
			expectedError: errNameTooShort,
		}, {
			name:          "error name must be at least 2 characters long after trimmed spaces",
			input:         testUserNameShortWithSpaces,
			wantError:     true,
			expectedError: errNameTooShort,
		}, {
			name:          "error name cannot exceed 100 characters",
			input:         testUserNameLong,
			wantError:     true,
			expectedError: errNameTooLong,
		}, {
			name:          "error age must be positive",
			input:         testUserAgeNeg,
			wantError:     true,
			expectedError: errAgeNeg,
		}, {
			name:          "error age seems unrealistic",
			input:         testUserAgeUnreal,
			wantError:     true,
			expectedError: errAgeUnreal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			v := NewUserValidator(nil)
			err := v.Create(context.Background(), tc.input)

			if tc.wantError {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)

				return
			}
			
			require.NoError(t, err)
		})
	}
}

func TestUserValidator_Update(t *testing.T) {
	type testCase struct {
		name          string
		input         *dto.User
		id            uint
		setupMocks    func(*mock_repository.MockUserRepo)
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name:  "successful updation",
			input: testUser,
			id:    id,
			setupMocks: func(mockRepo *mock_repository.MockUserRepo) {
				mockRepo.EXPECT().
					Exists(gomock.Any(), id).
					Return(true, nil)
			},
		}, {
			name:          "error user object cannot be nil",
			id:            id,
			wantError:     true,
			expectedError: errUserNil,
		}, {
			name:          "error name is required",
			input:         testUserNameEmpty,
			id:            id,
			wantError:     true,
			expectedError: errNameEmpty,
		}, {
			name:          "error name must be at least 2 characters long",
			input:         testUserNameShort,
			id:            id,
			wantError:     true,
			expectedError: errNameTooShort,
		}, {
			name:          "error name must be at least 2 characters long after trimmed spaces",
			input:         testUserNameShortWithSpaces,
			wantError:     true,
			expectedError: errNameTooShort,
		}, {
			name:          "error name cannot exceed 100 characters",
			input:         testUserNameLong,
			id:            id,
			wantError:     true,
			expectedError: errNameTooLong,
		}, {
			name:          "error age must be positive",
			input:         testUserAgeNeg,
			id:            id,
			wantError:     true,
			expectedError: errAgeNeg,
		}, {
			name:          "error age seems unrealistic",
			input:         testUserAgeUnreal,
			id:            id,
			wantError:     true,
			expectedError: errAgeUnreal,
		}, {
			name:  "error failed to check",
			input: testUser,
			id:    id,
			setupMocks: func(mockRepo *mock_repository.MockUserRepo) {
				mockRepo.EXPECT().
					Exists(gomock.Any(), id).
					Return(false, errRepo)
			},
			wantError:     true,
			expectedError: errFailedToCheck,
		}, {
			name:  "error user not found",
			input: testUser,
			id:    id,
			setupMocks: func(mockRepo *mock_repository.MockUserRepo) {
				mockRepo.EXPECT().
					Exists(gomock.Any(), id).
					Return(false, nil)
			},
			wantError:     true,
			expectedError: errUserNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockUserRepo(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(mockRepo)
			}

			v := NewUserValidator(mockRepo)
			err := v.Update(context.Background(), tc.input, tc.id)

			if tc.wantError {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)

				return
			}
				
			require.NoError(t, err)
		})
	}
}

func TestUserValidator_Delete(t *testing.T) {
	type testCase struct {
		name          string
		id            uint
		setupMocks    func(*mock_repository.MockUserRepo)
		wantError     bool
		expectedError error
	}

	cases := []testCase{
		{
			name: "successful deletion",
			id:   id,
			setupMocks: func(mockRepo *mock_repository.MockUserRepo) {
				mockRepo.EXPECT().
					Exists(gomock.Any(), id).
					Return(true, nil)
			},
		}, {
			name: "error failed to check",
			id:   id,
			setupMocks: func(mockRepo *mock_repository.MockUserRepo) {
				mockRepo.EXPECT().
					Exists(gomock.Any(), id).
					Return(false, errRepo)
			},
			wantError:     true,
			expectedError: errFailedToCheck,
		}, {
			name: "error user not found",
			id:   id,
			setupMocks: func(mockRepo *mock_repository.MockUserRepo) {
				mockRepo.EXPECT().
					Exists(gomock.Any(), id).
					Return(false, nil)
			},
			wantError:     true,
			expectedError: errUserNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockUserRepo(ctrl)
			tc.setupMocks(mockRepo)

			v := NewUserValidator(mockRepo)
			err := v.Delete(context.Background(), tc.id)

			if tc.wantError {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)

				return
			}
				
			require.NoError(t, err)
		})
	}
}
