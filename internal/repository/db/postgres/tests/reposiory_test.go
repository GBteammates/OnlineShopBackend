//go:build integration

package postgres_test

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	store  *repository.PGres
	logger = zap.NewExample().Sugar()
)

func init() {
	var err error
	store, err = repository.NewPgxStorage(context.Background(), logger, "postgresql://localhost:5432/shop?user=shopteam&password=123")
	if err != nil {
		panic("can't connect to db")
	}
}

func TestCategoryCreate(t *testing.T) {
	var err error
	cat := repository.NewCategoryRepo(store, logger)
	id, err := cat.CreateCategory(context.Background(), &models.Category{
		Name:        "testCat",
		Description: "Description",
	})
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories WHERE id=$1`, id)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestGetCatList(t *testing.T) {
	var err error
	cat := repository.NewCategoryRepo(store, logger)
	store.GetPool().Exec(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des'), ('2', '2des'), ('3', '3des')`)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	res, err := cat.GetCategoryList(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, res)
	for c := range res {
		require.Contains(t, []string{"1", "2", "3"}, c.Name)
		require.Contains(t, []string{"1des", "2des", "3des"}, c.Description)
	}
}

func TestUserCreate(t *testing.T) {
	var err error
	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ($1, $2) RETURNING id`,
		user.Rights.Name, user.Rights.Rules)
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)
	u := repository.NewUser(store, logger)
	res, err := u.Create(context.Background(), &user)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, res.ID)
}

func TestGetUSer(t *testing.T) {
	var err error
	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	_, err = store.GetPool().Exec(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)
	u := repository.NewUser(store, logger)
	res, err := u.GetUserByEmail(context.Background(), "123@mail.ru")
	assert.NoError(t, err)
	require.Equal(t, user.Firstname, res.Firstname)
	require.Equal(t, user.Password, res.Password)
	require.Equal(t, user.Email, res.Email)
}

func TestItemCreate(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item := repository.NewItemRepo(store, logger)
	id, err := item.CreateItem(context.Background(), &models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	})
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestItemUpdate(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item.Title,
		item.Category.Id,
		item.Description,
		item.Price,
		item.Vendor,
	)
	row.Scan(&item.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	newItem := models.Item{
		Id:          item.Id,
		Title:       "NewName",
		Description: "desc",
		Price:       500,
		Category:    cat,
	}

	itemEx := repository.NewItemRepo(store, logger)
	err = itemEx.UpdateItem(context.Background(), &newItem)
	assert.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `SELECT name, price FROM items`)
	row.Scan(&item.Title, &item.Price)
	require.Equal(t, newItem.Title, item.Title)
	require.Equal(t, newItem.Price, item.Price)
}

func TestItemGet(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
		Image:       "imageUrl",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description, picture) VALUES
	('1', '1des', 'imageUrl') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
		Vendor:      "vendor",
		Images:      []string{"1.jpg"},
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor, pictures)
	values ($1, $2, $3, $4, $5, $6) RETURNING id`,
		item.Title,
		item.Category.Id,
		item.Description,
		item.Price,
		item.Vendor,
		item.Images,
	)
	row.Scan(&item.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	itm := repository.NewItemRepo(store, logger)
	res, err := itm.GetItem(context.TODO(), item.Id)
	require.NoError(t, err)
	require.Equal(t, item.Id, res.Id)
	require.Equal(t, item.Title, res.Title)
}

func TestItemSearchLine(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	itm := repository.NewItemRepo(store, logger)
	ch, err := itm.SearchLine(context.Background(), "test")
	assert.NoError(t, err)
	for r := range ch {
		require.Equal(t, item1.Title, r.Title)
		require.Equal(t, item1.Id, r.Id)
	}

}

func TestItemItemsList(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	itm := repository.NewItemRepo(store, logger)
	ch, err := itm.ItemsList(context.Background())
	assert.NoError(t, err)
	for r := range ch {
		assert.Contains(t, item1.Title, r.Title)
		assert.Equal(t, item1.Description, r.Description)
	}

}

func TestCartCreate(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)

	err = row.Scan(&user.ID)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)

	cartMdl := models.Cart{
		UserId:   user.ID,
		Items:    []models.ItemWithQuantity{{Item: item1, Quantity: 1}, {Item: item2, Quantity: 1}},
		ExpireAt: time.Now().Add(time.Hour * 2),
	}

	crt := repository.NewCartStore(store, logger)
	res, err := crt.Create(context.Background(), cartMdl.UserId)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM carts`)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, res)
}

func TestCartAddItem(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)

	err = row.Scan(&user.ID)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)

	cartMdl := models.Cart{
		UserId:   user.ID,
		Items:    []models.ItemWithQuantity{{Item: item1, Quantity: 1}},
		ExpireAt: time.Now().Add(time.Hour * 2),
	}

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO carts (user_id, expire_at) VALUES ($1, $2) RETURNING id`,
		cartMdl.UserId, cartMdl.ExpireAt)
	err = row.Scan(&cartMdl.Id)
	require.NoError(t, err)
	defer store.GetPool().Exec(context.Background(), `DELETE from carts`)
	crt := repository.NewCartStore(store, logger)
	err = crt.AddItemToCart(context.Background(), cartMdl.Id, item2.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE from cart_items`)
	require.NoError(t, err)
	row = store.GetPool().QueryRow(context.Background(), `SELECT COUNT(cart_id) FROM cart_items`)
	var count int
	err = row.Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestCartDelete(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)

	err = row.Scan(&user.ID)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)

	cartMdl := models.Cart{
		UserId:   user.ID,
		Items:    []models.ItemWithQuantity{{Item: item1, Quantity: 1}},
		ExpireAt: time.Now().Add(time.Hour * 2),
	}

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO carts (user_id, expire_at) VALUES ($1, $2) RETURNING id`,
		cartMdl.UserId, cartMdl.ExpireAt)
	err = row.Scan(&cartMdl.Id)
	require.NoError(t, err)
	store.GetPool().Exec(context.Background(), `INSERT INTO cart_items (cart_id, item_id) VALUES ($1, $2)`, cartMdl.Id, item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE from carts`)
	crt := repository.NewCartStore(store, logger)
	err = crt.DeleteCart(context.Background(), cartMdl.Id)
	require.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `SELECT COUNT(cart_id) FROM cart_items`)
	var count int
	err = row.Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestCartDeleteItem(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)

	err = row.Scan(&user.ID)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)

	cartMdl := models.Cart{
		UserId:   user.ID,
		Items:    []models.ItemWithQuantity{{Item: item1, Quantity: 1}},
		ExpireAt: time.Now().Add(time.Hour * 2),
	}

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO carts (user_id, expire_at) VALUES ($1, $2) RETURNING id`,
		cartMdl.UserId, cartMdl.ExpireAt)
	err = row.Scan(&cartMdl.Id)
	require.NoError(t, err)
	defer store.GetPool().Exec(context.Background(), `DELETE from carts`)
	store.GetPool().Exec(context.Background(), `INSERT INTO cart_items (cart_id, item_id, item_quantity) VALUES ($1, $2, $3)`, cartMdl.Id, item1.Id, cartMdl.Items[0].Quantity)
	crt := repository.NewCartStore(store, logger)
	err = crt.DeleteItemFromCart(context.Background(), cartMdl.Id, item1.Id)
	require.NoError(t, err)
	row = store.GetPool().QueryRow(context.Background(), `SELECT COUNT(cart_id) FROM cart_items`)
	var count int
	err = row.Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

}

func TestOrderCreate(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)

	err = row.Scan(&user.ID)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)

	ordr := repository.NewOrderRepo(store, logger)
	order := models.Order{
		ShipmentTime: time.Now().Add(2 * time.Hour),
		User:         user,
		Address:      user.Address,
		Status:       models.StatusProcessed,
		Items:        []models.ItemWithQuantity{{Item: item1, Quantity: 1}, {Item: item2, Quantity: 1}},
	}
	res, err := ordr.Create(context.Background(), &order)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM orders`)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM order_items`)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, res.ID)
}

func TestOrderDelete(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)

	err = row.Scan(&user.ID)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)

	order := models.Order{
		ShipmentTime: time.Now().Add(2 * time.Hour),
		User:         user,
		Address:      user.Address,
		Status:       models.StatusProcessed,
		Items:        []models.ItemWithQuantity{{Item: item1, Quantity: 1}, {Item: item2, Quantity: 1}},
	}

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO orders (shipment_time, user_id, status, address) 
	VALUES ($1, $2, $3, $4) RETURNING id`, order.ShipmentTime, order.User.ID, order.Status,
		fmt.Sprintf("%s %s %s %s", order.User.Address.Zipcode, order.User.Address.Country, order.User.Address.City, order.User.Address.Street))
	row.Scan(&order.ID)

	row = store.GetPool().QueryRow(context.Background(),
		`INSERT INTO order_items (order_id, item_id) VALUES ($1, $2), ($1, $3)`, order.ID, order.Items[0].Id, order.ID, order.Items[1].Id)

	rdrRp := repository.NewOrderRepo(store, logger)
	err = rdrRp.DeleteOrder(context.Background(), &order)
	require.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `SELECT COUNT(id) FROM orders`)
	var count int
	row.Scan(&count)
	require.Equal(t, 0, count)
}

func TestOrderChangeAddres(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)

	err = row.Scan(&user.ID)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)

	order := models.Order{
		ShipmentTime: time.Now().Add(2 * time.Hour),
		User:         user,
		Address:      user.Address,
		Status:       models.StatusProcessed,
		Items:        []models.ItemWithQuantity{{Item: item1, Quantity: 1}, {Item: item2, Quantity: 1}},
	}

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO orders (shipment_time, user_id, status, address) 
	VALUES ($1, $2, $3, $4) RETURNING id`, order.ShipmentTime, order.User.ID, order.Status,
		fmt.Sprintf("%s %s %s %s", order.User.Address.Zipcode, order.User.Address.Country, order.User.Address.City, order.User.Address.Street))
	row.Scan(&order.ID)

	row = store.GetPool().QueryRow(context.Background(),
		`INSERT INTO order_items (order_id, item_id) VALUES ($1, $2), ($1, $3)`, order.ID, order.Items[0].Id, order.ID, order.Items[1].Id)

	rdrRp := repository.NewOrderRepo(store, logger)
	err = rdrRp.ChangeAddress(context.Background(), &order, models.UserAddress{
		Zipcode: "123456",
		Country: "Kyrgizstan",
		Street:  "Baitik Batyr, 10",
		City:    "Bishkek",
	})
	defer store.GetPool().Exec(context.Background(), `DELETE FROM orders`)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM order_items`)
	require.NoError(t, err)
	row = store.GetPool().QueryRow(context.Background(), `SELECT address FROM orders`)
	var addr string
	row.Scan(&addr)
	assert.Contains(t, addr, "Bishkek")

}

func TestOrderChangeStatus(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)

	err = row.Scan(&user.ID)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)

	order := models.Order{
		ShipmentTime: time.Now().Add(2 * time.Hour),
		User:         user,
		Address:      user.Address,
		Status:       models.StatusProcessed,
		Items:        []models.ItemWithQuantity{{Item: item1, Quantity: 1}, {Item: item2, Quantity: 1}},
	}

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO orders (shipment_time, user_id, status, address) 
	VALUES ($1, $2, $3, $4) RETURNING id`, order.ShipmentTime, order.User.ID, order.Status,
		fmt.Sprintf("%s %s %s %s", order.User.Address.Zipcode, order.User.Address.Country, order.User.Address.City, order.User.Address.Street))
	row.Scan(&order.ID)

	row = store.GetPool().QueryRow(context.Background(),
		`INSERT INTO order_items (order_id, item_id) VALUES ($1, $2), ($1, $3)`, order.ID, order.Items[0].Id, order.ID, order.Items[1].Id)

	rdrRp := repository.NewOrderRepo(store, logger)
	err = rdrRp.ChangeStatus(context.Background(), &order, models.StatusCourier)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM orders`)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM order_items`)
	require.NoError(t, err)
	row = store.GetPool().QueryRow(context.Background(), `SELECT status FROM orders`)
	var status models.Status
	row.Scan(&status)
	assert.Equal(t, models.StatusCourier, status)

}

func TestOrdersGetOrderByID(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)

	err = row.Scan(&user.ID)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)

	order := models.Order{
		ShipmentTime: time.Now().Add(2 * time.Hour),
		User:         user,
		Address:      user.Address,
		Status:       models.StatusProcessed,
		Items:        []models.ItemWithQuantity{{Item: item1, Quantity: 1}, {Item: item2, Quantity: 1}},
	}

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO orders (shipment_time, user_id, status, address) 
	VALUES ($1, $2, $3, $4) RETURNING id`, order.ShipmentTime, order.User.ID, order.Status,
		fmt.Sprintf("%s -> %s -> %s -> %s", order.User.Address.Zipcode, order.User.Address.Country, order.User.Address.City, order.User.Address.Street))
	row.Scan(&order.ID)
	fmt.Printf("order id %s: \n", order.ID.String())
	_, err = store.GetPool().Exec(context.Background(),
		`INSERT INTO order_items (order_id, item_id, item_quantity) VALUES ($1, $2, $3), ($1, $4, $5)`, order.ID, order.Items[0].Id, order.Items[0].Quantity, order.Items[1].Id, order.Items[1].Quantity)
	require.NoError(t, err)
	rdrRp := repository.NewOrderRepo(store, logger)
	res, err := rdrRp.GetOrderByID(context.Background(), order.ID)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM orders`)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM order_items`)
	require.NoError(t, err)
	require.Equal(t, order.Items[0].Title, res.Items[0].Title)
	require.Equal(t, order.ID, res.ID)
	require.Equal(t, order.Address, res.Address)
}

func TestOrdersGetOrders(t *testing.T) {
	var err error

	cat := models.Category{
		Name:        "1",
		Description: "1des",
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des') RETURNING id`)
	err = row.Scan(&cat.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	assert.NoError(t, err)

	item1 := models.Item{
		Title:       "testItem",
		Description: "desc",
		Price:       300,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item1.Title,
		item1.Category.Id,
		item1.Description,
		item1.Price,
		item1.Vendor,
	)
	row.Scan(&item1.Id)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM items`)

	item2 := models.Item{
		Title:       "Item",
		Description: "desc",
		Price:       400,
		Category:    cat,
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item2.Title,
		item2.Category.Id,
		item2.Description,
		item2.Price,
		item2.Vendor,
	)
	row.Scan(&item2.Id)

	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)

	err = row.Scan(&user.ID)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)

	order := models.Order{
		ShipmentTime: time.Now().Add(2 * time.Hour),
		User:         user,
		Address:      user.Address,
		Status:       models.StatusProcessed,
		Items:        []models.ItemWithQuantity{{Item: item1, Quantity: 1}, {Item: item2, Quantity: 1}},
	}

	order2 := models.Order{
		ShipmentTime: time.Now().Add(2 * time.Hour),
		User:         user,
		Address:      user.Address,
		Status:       models.StatusCourier,
		Items:        []models.ItemWithQuantity{{Item: item1, Quantity: 1}},
	}

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO orders (shipment_time, user_id, status, address) 
	VALUES ($1, $2, $3, $4) RETURNING id`, order.ShipmentTime, order.User.ID, order.Status,
		fmt.Sprintf("%s -> %s -> %s -> %s", order.User.Address.Zipcode, order.User.Address.Country, order.User.Address.City, order.User.Address.Street))
	row.Scan(&order.ID)

	row = store.GetPool().QueryRow(context.Background(), `INSERT INTO orders (shipment_time, user_id, status, address) 
	VALUES ($1, $2, $3, $4) RETURNING id`, order2.ShipmentTime, order2.User.ID, order2.Status,
		fmt.Sprintf("%s -> %s -> %s -> %s", order2.User.Address.Zipcode, order2.User.Address.Country, order2.User.Address.City, order2.User.Address.Street))
	row.Scan(&order2.ID)

	_, err = store.GetPool().Exec(context.Background(),
		`INSERT INTO order_items (order_id, item_id, item_quantity) VALUES ($1, $2, $3), ($1, $4, $5)`, order.ID, order.Items[0].Id, order.Items[0].Quantity, order.Items[1].Id, order.Items[1].Quantity)
	require.NoError(t, err)
	_, err = store.GetPool().Exec(context.Background(),
		`INSERT INTO order_items (order_id, item_id, item_quantity) VALUES ($1, $2, $3)`, order2.ID, order2.Items[0].Id, order2.Items[0].Quantity)
	require.NoError(t, err)
	rdrRp := repository.NewOrderRepo(store, logger)
	ch, err := rdrRp.GetOrdersForUser(context.Background(), &user)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM orders`)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM order_items`)
	require.NoError(t, err)
	res := make([]models.Order, 0, 2)
	for o := range ch {
		res = append(res, o)
	}
	require.Equal(t, order.Items[0].Title, res[0].Items[0].Title)
	// require.Equal(t, order.ID, res[0].ID)
	require.Equal(t, order.Address, res[0].Address)

	require.Equal(t, order2.Items[0].Title, res[1].Items[0].Title)
	// require.Equal(t, order2.ID, res[1].ID)
	require.Equal(t, order2.Address, res[1].Address)
}
