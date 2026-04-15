package product

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockRepo struct {
	items  map[uint]*Product
	nextID uint
}

func newMockRepo() *mockRepo {
	return &mockRepo{items: make(map[uint]*Product), nextID: 1}
}

func (m *mockRepo) Create(p *Product) error {
	p.ID = m.nextID
	m.nextID++
	m.items[p.ID] = p
	return nil
}
func (m *mockRepo) FindByID(id uint) (*Product, error) {
	if p, ok := m.items[id]; ok {
		return p, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *mockRepo) FindAll() ([]Product, error) {
	out := make([]Product, 0, len(m.items))
	for _, p := range m.items {
		out = append(out, *p)
	}
	return out, nil
}
func (m *mockRepo) Update(p *Product) error { m.items[p.ID] = p; return nil }
func (m *mockRepo) Delete(id uint) error    { delete(m.items, id); return nil }

func TestCreateProduct_Success(t *testing.T) {
	svc := NewService(newMockRepo())

	p, err := svc.Create(CreateProductRequest{
		Name:  "Teclado",
		Price: 120.50,
		Stock: 10,
	}, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Teclado", p.Name)
	assert.Equal(t, uint(1), p.UserID)
	assert.Equal(t, uint(1), p.ID)
}

func TestUpdateProduct_PartialFields(t *testing.T) {
	svc := NewService(newMockRepo())

	p, _ := svc.Create(CreateProductRequest{Name: "Mouse", Price: 50, Stock: 5}, 1)

	newPrice := 75.0
	updated, err := svc.Update(p.ID, UpdateProductRequest{Price: &newPrice})

	assert.NoError(t, err)
	assert.Equal(t, 75.0, updated.Price)
	assert.Equal(t, "Mouse", updated.Name, "name should remain unchanged")
	assert.Equal(t, 5, updated.Stock)
}

func TestGetProductByID_NotFound(t *testing.T) {
	svc := NewService(newMockRepo())
	_, err := svc.GetByID(999)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
