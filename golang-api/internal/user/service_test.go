package user

import (
	"errors"
	"testing"

	"github.com/fixia/golang-api/config"
	"github.com/fixia/golang-api/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockRepo struct {
	users  map[uint]*User
	byMail map[string]*User
	nextID uint
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		users:  make(map[uint]*User),
		byMail: make(map[string]*User),
		nextID: 1,
	}
}

func (m *mockRepo) Create(u *User) error {
	if _, ok := m.byMail[u.Email]; ok {
		return errors.New("duplicate")
	}
	u.ID = m.nextID
	m.nextID++
	m.users[u.ID] = u
	m.byMail[u.Email] = u
	return nil
}

func (m *mockRepo) FindByID(id uint) (*User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepo) FindByEmail(email string) (*User, error) {
	if u, ok := m.byMail[email]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepo) FindAll() ([]User, error) {
	out := make([]User, 0, len(m.users))
	for _, u := range m.users {
		out = append(out, *u)
	}
	return out, nil
}

func (m *mockRepo) Update(u *User) error {
	m.users[u.ID] = u
	m.byMail[u.Email] = u
	return nil
}

func (m *mockRepo) Delete(id uint) error {
	if u, ok := m.users[id]; ok {
		delete(m.byMail, u.Email)
		delete(m.users, id)
	}
	return nil
}

func setupTestConfig() {
	config.AppConfig = &config.Config{
		JWTSecret: "test_secret",
		JWTExpHrs: 1,
	}
}

func TestRegister_Success(t *testing.T) {
	setupTestConfig()
	svc := NewService(newMockRepo())

	u, err := svc.Register(RegisterRequest{
		FirstName: "Ada",
		LastName:  "Lovelace",
		Email:     "ada@example.com",
		Password:  "supersecret",
	})

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, "ada@example.com", u.Email)
	assert.Equal(t, "user", u.Role)
	assert.NotEqual(t, "supersecret", u.Password, "password must be hashed")
}

func TestRegister_DuplicateEmail(t *testing.T) {
	setupTestConfig()
	svc := NewService(newMockRepo())

	req := RegisterRequest{
		FirstName: "Ada", LastName: "L",
		Email: "dup@example.com", Password: "supersecret",
	}
	_, err := svc.Register(req)
	assert.NoError(t, err)

	_, err = svc.Register(req)
	assert.ErrorIs(t, err, utils.ErrDuplicateEmail)
}

func TestLogin_Success(t *testing.T) {
	setupTestConfig()
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.Register(RegisterRequest{
		FirstName: "Ada", LastName: "L",
		Email: "ada@example.com", Password: "supersecret",
	})
	assert.NoError(t, err)

	resp, err := svc.Login(LoginRequest{Email: "ada@example.com", Password: "supersecret"})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "ada@example.com", resp.User.Email)
}

func TestLogin_WrongPassword(t *testing.T) {
	setupTestConfig()
	svc := NewService(newMockRepo())

	_, _ = svc.Register(RegisterRequest{
		FirstName: "Ada", LastName: "L",
		Email: "ada@example.com", Password: "supersecret",
	})

	_, err := svc.Login(LoginRequest{Email: "ada@example.com", Password: "wrongpass"})
	assert.ErrorIs(t, err, utils.ErrInvalidCredentials)
}

func TestLogin_UserNotFound(t *testing.T) {
	setupTestConfig()
	svc := NewService(newMockRepo())

	_, err := svc.Login(LoginRequest{Email: "ghost@example.com", Password: "whatever"})
	assert.ErrorIs(t, err, utils.ErrInvalidCredentials)
}

func TestPasswordIsBcryptHashed(t *testing.T) {
	setupTestConfig()
	svc := NewService(newMockRepo())

	u, err := svc.Register(RegisterRequest{
		FirstName: "Ada", LastName: "L",
		Email: "ada@example.com", Password: "supersecret",
	})
	assert.NoError(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("supersecret"))
	assert.NoError(t, err, "password must be a valid bcrypt hash of the original")
}
