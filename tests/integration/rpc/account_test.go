package tests

import (
	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/docker/distribution/context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

// Users

var (
	testUser = account.SignUpRequest{
		Name:     "user",
		Password: "userPassword",
		Email:    "user@amp.io",
	}
)

func createUser(t *testing.T, user *account.SignUpRequest) context.Context {
	// SignUp
	_, err := accountClient.SignUp(ctx, user)
	assert.NoError(t, err)

	// Create a verify token
	token, err := auth.CreateToken(user.Name, auth.TokenTypeVerify, time.Hour)
	assert.NoError(t, err)

	// Verify
	_, err = accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, err)

	// Create a login token
	token, err = auth.CreateToken(user.Name, auth.TokenTypeLogin, time.Hour)
	return metadata.NewContext(ctx, metadata.Pairs(auth.TokenKey, token))
}

func TestUserShouldSignUpAndVerify(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// Create a token
	token, err := auth.CreateToken(testUser.Name, auth.TokenTypeVerify, time.Hour)
	assert.NoError(t, err)

	// Verify
	_, err = accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, err)
}

func TestUserSignUpInvalidNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Name = "UpperCaseIsNotAllowed"
	_, err := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpInvalidEmailShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Email = "this is not an email"
	_, err := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Password = ""
	_, err := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpAlreadyExistsShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// SignUp
	_, err = accountClient.SignUp(ctx, &testUser)
	assert.Error(t, err)
}

func TestUserVerifyNotATokenShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// Verify
	_, err = accountClient.Verify(ctx, &account.VerificationRequest{Token: "this is not a token"})
	assert.Error(t, err)
}

func TestUserVerifyNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a verify token
	token, err := auth.CreateToken("nonexistinguser", auth.TokenTypeVerify, time.Hour)
	assert.NoError(t, err)

	// Verify
	_, err = accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.Error(t, err)
}

// TODO: Check expired token

func TestUserLogin(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Login
	_, err := accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: testUser.Password,
	})
	assert.NoError(t, err)
}

func TestUserLoginNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Login
	_, err := accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: testUser.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginNonVerifiedUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: testUser.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginInvalidNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Login
	_, err := accountClient.Login(ctx, &account.LogInRequest{
		Name:     "not the right user name",
		Password: testUser.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Login
	_, err := accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: "not the right password",
	})
	assert.Error(t, err)
}

func TestUserPasswordReset(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)
}

func TestUserPasswordResetMalformedRequestShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: "this is not a valid user name"})
	assert.Error(t, err)
}

func TestUserPasswordResetNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: "nonexistinguser"})
	assert.Error(t, err)
}

func TestUserPasswordSet(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Set
	token, _ := auth.CreateToken(testUser.Name, auth.TokenTypePassword, time.Hour)
	_, err := accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    token,
		Password: "newPassword",
	})
	assert.NoError(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: "newPassword",
	})
	assert.NoError(t, err)
}

func TestUserPasswordSetInvalidTokenShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)

	// Password Set
	_, err = accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    "this is an invalid token",
		Password: "newPassword",
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: "newPassword",
	})
	assert.Error(t, err)
}

func TestUserPasswordSetNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)

	// Password Set
	token, _ := auth.CreateToken("nonexistinguser", auth.TokenTypePassword, time.Hour)
	_, err = accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    token,
		Password: "newPassword",
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: "newPassword",
	})
	assert.Error(t, err)
}

func TestUserPasswordSetInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)

	// Password Set
	token, _ := auth.CreateToken(testUser.Name, auth.TokenTypePassword, time.Hour)
	_, err = accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    token,
		Password: "",
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: "",
	})
	assert.Error(t, err)
}

func TestUserPasswordChange(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := "newPassword"
	_, err := accountClient.PasswordChange(ownerCtx, &account.PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.NoError(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.NoError(t, err)
}

func TestUserPasswordChangeInvalidExistingPassword(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := "newPassword"
	_, err := accountClient.PasswordChange(ownerCtx, &account.PasswordChangeRequest{
		ExistingPassword: "this is not the right password",
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserPasswordChangeEmptyNewPassword(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := ""
	_, err := accountClient.PasswordChange(ownerCtx, &account.PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserPasswordChangeInvalidNewPassword(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := "aze"
	_, err := accountClient.PasswordChange(ownerCtx, &account.PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserForgotLogin(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// ForgotLogin
	_, err = accountClient.ForgotLogin(ctx, &account.ForgotLoginRequest{
		Email: testUser.Email,
	})
	assert.NoError(t, err)
}

func TestUserForgotLoginMalformedEmailShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// ForgotLogin
	_, err = accountClient.ForgotLogin(ctx, &account.ForgotLoginRequest{
		Email: "this is not a valid email",
	})
	assert.Error(t, err)
}

func TestUserForgotLoginNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// ForgotLogin
	_, err := accountClient.ForgotLogin(ctx, &account.ForgotLoginRequest{
		Email: testUser.Email,
	})
	assert.Error(t, err)
}

func TestUserGet(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Get
	getReply, err := accountClient.GetUser(ctx, &account.GetUserRequest{
		Name: testUser.Name,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, getReply)
	assert.Equal(t, getReply.User.Name, testUser.Name)
	assert.Equal(t, getReply.User.Email, testUser.Email)
	assert.NotEmpty(t, getReply.User.CreateDt)
}

func TestUserGetMalformedUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Get
	_, err := accountClient.GetUser(ctx, &account.GetUserRequest{
		Name: "this user is malformed",
	})
	assert.Error(t, err)
}

func TestUserGetNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Get
	_, err := accountClient.GetUser(ctx, &account.GetUserRequest{
		Name: testUser.Name,
	})
	assert.Error(t, err)
}

func TestUserList(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// List
	listReply, err := accountClient.ListUsers(ctx, &account.ListUsersRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, listReply)
	assert.Len(t, listReply.Users, 1)
	assert.Equal(t, listReply.Users[0].Name, testUser.Name)
	assert.Equal(t, listReply.Users[0].Email, testUser.Email)
	assert.NotEmpty(t, listReply.Users[0].CreateDt)
}

// Organizations

var (
	testOrg = account.CreateOrganizationRequest{
		Name:  "organization",
		Email: "organization@amp.io",
	}
	testMember = account.SignUpRequest{
		Name:     "organization-member",
		Password: "organizationMemberPassword",
		Email:    "organization.member@amp.io",
	}
)

func createOrganization(t *testing.T, org *account.CreateOrganizationRequest, owner *account.SignUpRequest) context.Context {
	// Create a user
	ownerCtx := createUser(t, owner)

	// CreateOrganization
	_, err := accountClient.CreateOrganization(ownerCtx, org)
	assert.NoError(t, err)

	return ownerCtx
}

func addUserToOrganization(t *testing.T, org *account.CreateOrganizationRequest, ownerCtx context.Context, user *account.SignUpRequest) context.Context {
	// Create a user
	userCtx := createUser(t, user)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: org.Name,
		UserName:         user.Name,
	})
	assert.NoError(t, err)
	return userCtx
}

func TestOrganizationCreate(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// CreateOrganization
	_, err := accountClient.CreateOrganization(ownerCtx, &testOrg)
	assert.NoError(t, err)
}

func TestOrganizationCreateInvalidNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// CreateOrganization
	invalidRequest := testOrg
	invalidRequest.Name = "this is not a valid name"
	_, err := accountClient.CreateOrganization(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestOrganizationCreateInvalidEmailShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// CreateOrganization
	invalidRequest := testOrg
	invalidRequest.Email = "this is not a valid email"
	_, err := accountClient.CreateOrganization(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestOrganizationCreateAlreadyExistsShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// CreateOrganization again
	_, err := accountClient.CreateOrganization(ownerCtx, &testOrg)
	assert.Error(t, err)
}

func TestOrganizationAddUser(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationAddUserInvalidOrganizationNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserInvalidUserNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserToNonExistingOrganizationShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create owner
	ownerCtx := createUser(t, &testUser)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserNotOwnerShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	createOrganization(t, &testOrg, &testUser)

	// Create member
	memberCtx := createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(memberCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddNonValidatedUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// SignUp member
	_, err := accountClient.SignUp(ctx, &testMember)
	assert.NoError(t, err)

	// AddUserToOrganization
	_, err = accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddSameUserTwiceShouldSucceed(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// AddUserToOrganization
	_, err = accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationRemoveUser(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationRemoveUserInvalidOrganizationNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserInvalidUserNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserFromNonExistingOrganizationShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create user
	ownerCtx := createUser(t, &testUser)

	// Create member
	createUser(t, &testMember)

	// RemoveUserFromOrganization
	_, err := accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserNotOwnerShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	memberCtx := createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(memberCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// RemoveUserFromOrganization
	_, err := accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveSameUserTwiceShouldSucceed(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationRemoveAllOwnersShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testUser.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationList(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createOrganization(t, &testOrg, &testUser)

	// List
	listReply, err := accountClient.ListOrganizations(ctx, &account.ListOrganizationsRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, listReply)
	assert.Len(t, listReply.Organizations, 1)
	assert.Equal(t, listReply.Organizations[0].Name, testOrg.Name)
	assert.Equal(t, listReply.Organizations[0].Email, testOrg.Email)
	assert.NotEmpty(t, listReply.Organizations[0].CreateDt)
	assert.NotEmpty(t, listReply.Organizations[0].Members)
	assert.Equal(t, listReply.Organizations[0].Members[0].Name, testUser.Name)
	assert.Equal(t, listReply.Organizations[0].Members[0].Role, schema.OrganizationRole_ORGANIZATION_OWNER)
}

// Teams

var (
	testTeam = account.CreateTeamRequest{
		OrganizationName: testOrg.Name,
		TeamName:         "team",
	}
)

func createTeam(t *testing.T, org *account.CreateOrganizationRequest, owner *account.SignUpRequest, team *account.CreateTeamRequest) context.Context {
	// Create a user
	ownerCtx := createOrganization(t, org, owner)

	// CreateTeam
	_, err := accountClient.CreateTeam(ownerCtx, team)
	assert.NoError(t, err)

	return ownerCtx
}

func TestTeamCreate(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// CreateTeam
	_, err := accountClient.CreateTeam(ownerCtx, &testTeam)
	assert.NoError(t, err)
}

func TestTeamCreateInvalidOrganizationNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// CreateTeam
	invalidRequest := testTeam
	invalidRequest.OrganizationName = "this is not a valid name"
	_, err := accountClient.CreateTeam(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestTeamCreateInvalidTeamNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// CreateTeam
	invalidRequest := testTeam
	invalidRequest.TeamName = "this is not a valid name"
	_, err := accountClient.CreateTeam(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestTeamCreateNonExistingOrganizationShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// CreateTeam
	invalidRequest := testTeam
	invalidRequest.OrganizationName = "nonExistingOrg"
	_, err := accountClient.CreateTeam(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestTeamCreateNotOrgOwnerShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	createOrganization(t, &testOrg, &testUser)

	// Create a user not part of the organization
	notOrgOwnerCtx := createUser(t, &testMember)

	// CreateTeam
	_, err := accountClient.CreateTeam(notOrgOwnerCtx, &testTeam)
	assert.Error(t, err)
}

func TestTeamCreateAlreadyExistsShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// CreateTeam again
	_, err := accountClient.CreateTeam(ownerCtx, &testTeam)
	assert.Error(t, err)
}

func TestTeamAddUser(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)
	addUserToOrganization(t, &testOrg, ownerCtx, &testMember)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestTeamAddUserInvalidOrganizationNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: "this is not a valid name",
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserInvalidTeamNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserInvalidUserNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestTeamAddUserToNonExistingOrganizationShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createUser(t, &testUser)
	createUser(t, &testMember)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserToNonExistingTeamShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createOrganization(t, &testOrg, &testUser)
	createUser(t, &testMember)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddNonExistingUserToTeamShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserNotOrganizationOwnerShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	createTeam(t, &testOrg, &testUser, &testTeam)
	memberCtx := createUser(t, &testMember)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(memberCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddNonValidatedUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// SignUp member
	_, err := accountClient.SignUp(ctx, &testMember)
	assert.NoError(t, err)

	// AddUserToTeam
	_, err = accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddSameUserTwiceShouldSucceed(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member
	createUser(t, &testMember)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// AddUserToTeam again
	_, err = accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestTeamRemoveUser(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member
	createUser(t, &testMember)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = accountClient.RemoveUserFromTeam(ownerCtx, &account.RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestTeamRemoveUserInvalidOrganizationNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member
	createUser(t, &testMember)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = accountClient.RemoveUserFromTeam(ownerCtx, &account.RemoveUserFromTeamRequest{
		OrganizationName: "this is not a valid name",
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserInvalidTeamNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member
	createUser(t, &testMember)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = accountClient.RemoveUserFromTeam(ownerCtx, &account.RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserInvalidUserNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member
	createUser(t, &testMember)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = accountClient.RemoveUserFromTeam(ownerCtx, &account.RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserFromNonExistingOrganizationShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create user
	ownerCtx := createUser(t, &testUser)

	// Create member
	createUser(t, &testMember)

	// RemoveUserFromTeam
	_, err := accountClient.RemoveUserFromTeam(ownerCtx, &account.RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserFromNonExistingTeamShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create user
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// RemoveUserFromTeam
	_, err := accountClient.RemoveUserFromTeam(ownerCtx, &account.RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserNotOwnerShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	/// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member
	memberCtx := createUser(t, &testMember)

	// AddUserToTeam
	_, err := accountClient.AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = accountClient.RemoveUserFromTeam(memberCtx, &account.RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// RemoveUserFromTeam
	_, err := accountClient.RemoveUserFromTeam(ownerCtx, &account.RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserNotPartOfTheTeamShouldSucceed(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	/// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member
	createUser(t, &testMember)

	// RemoveUserFromTeam
	_, err := accountClient.RemoveUserFromTeam(ownerCtx, &account.RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestTeamList(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createTeam(t, &testOrg, &testUser, &testTeam)

	// List
	listReply, err := accountClient.ListTeams(ctx, &account.ListTeamsRequest{
		OrganizationName: testOrg.Name,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, listReply)
	assert.Len(t, listReply.Teams, 1)
	assert.Equal(t, listReply.Teams[0].Name, testTeam.TeamName)
	assert.NotEmpty(t, listReply.Teams[0].CreateDt)
	assert.NotEmpty(t, listReply.Teams[0].Members)
	assert.Equal(t, listReply.Teams[0].Members[0].Name, testUser.Name)
	assert.Equal(t, listReply.Teams[0].Members[0].Role, schema.TeamRole_TEAM_OWNER)
}

func TestTeamListInvalidOrganizationNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createTeam(t, &testOrg, &testUser, &testTeam)

	// List
	_, err := accountClient.ListTeams(ctx, &account.ListTeamsRequest{
		OrganizationName: "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestTeamListNonExistingOrganizationNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// List
	_, err := accountClient.ListTeams(ctx, &account.ListTeamsRequest{
		OrganizationName: testTeam.OrganizationName,
	})
	assert.Error(t, err)
}