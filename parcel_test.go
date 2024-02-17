package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

type ParcelTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (suite *ParcelTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite", "tracker.db")
	suite.Require().NoError(err, "unable to open DBMS connection")
	suite.db = db
}

func (suite *ParcelTestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *ParcelTestSuite) TearDownTest() {
	_, err := suite.db.Exec("DELETE FROM parcel")
	suite.Require().NoError(err, "unable to clean up parcel records after test")
}

func TestParcelTestSuite(t *testing.T) {
	suite.Run(t, new(ParcelTestSuite))
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func (suite *ParcelTestSuite) TestAddGetDelete() { //t *testing.T) {
	// prepare
	//db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	//require.NoErrorf(t, err, "unable to open DBMS connection, err=%v", err)
	//defer db.Close()

	store := NewParcelStore(suite.db)
	expectedParcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(expectedParcel)
	suite.Require().NoErrorf(err, "unable to add new parcel, err=%v", err)
	suite.Require().Greaterf(id, 0, "invalid parcel id=%d", id)

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной expectedParcel
	actualParcel, err := store.Get(id)
	suite.Require().NoErrorf(err, "unable to get parcel with id=%d, err=%v", id, err)
	suite.Assert().Equal(expectedParcel.Client, actualParcel.Client)
	suite.Assert().Equal(expectedParcel.Status, actualParcel.Status)
	suite.Assert().Equal(expectedParcel.Address, actualParcel.Address)
	suite.Assert().Equal(expectedParcel.CreatedAt, actualParcel.CreatedAt)

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	err = store.Delete(id)
	suite.Require().NoErrorf(err, "unable to delete parcel with id=%d, err=%v", id, err)

	// проверьте, что посылку больше нельзя получить из БД
	_, err = store.Get(id)
	suite.Require().Errorf(err, "got parcel id=%d after deletion", id)
}

// TestSetAddress проверяет обновление адреса
func (suite *ParcelTestSuite) TestSetAddress() { //t *testing.T) {
	// prepare
	//db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	//require.NoErrorf(t, err, "unable to open DBMS connection, err=%v", err)
	//defer db.Close()

	store := NewParcelStore(suite.db)
	expectedParcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(expectedParcel)
	suite.Require().NoErrorf(err, "unable to add new parcel, err=%v", err)
	suite.Require().Greaterf(id, 0, "invalid parcel id=%d", id)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	suite.Require().NoErrorf(err, "unable to set address, err=%v", err)

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	actualParsel, err := store.Get(id)
	suite.Require().NoErrorf(err, "unable to get parcel with id=%d, err=%v", id, err)
	suite.Assert().Equal(newAddress, actualParsel.Address)
}

// TestSetStatus проверяет обновление статуса
func (suite *ParcelTestSuite) TestSetStatus() { //t *testing.T) {
	// prepare
	//db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	//require.NoErrorf(t, err, "unable to open DBMS connection, err=%v", err)
	//defer db.Close()

	store := NewParcelStore(suite.db)
	expectedParcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(expectedParcel)
	suite.Require().NoErrorf(err, "unable to add new parcel, err=%v", err)
	suite.Require().Greaterf(id, 0, "invalid parcel id=%d", id)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	err = store.SetStatus(id, ParcelStatusSent)
	suite.Require().NoErrorf(err, "unable to set status, err=%v", err)

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	actualParsel, err := store.Get(id)
	suite.Require().NoErrorf(err, "unable to get parcel with id=%d, err=%v", id, err)
	suite.Assert().Equal(ParcelStatusSent, actualParsel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func (suite *ParcelTestSuite) TestGetByClient() { //t *testing.T) {
	// prepare
	//db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	//require.NoErrorf(t, err, "unable to open DBMS connection, err=%v", err)
	//defer db.Close()

	store := NewParcelStore(suite.db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		id, err := store.Add(parcels[i])
		suite.Require().NoErrorf(err, "unable to add new parcel, err=%v", err)
		suite.Require().Greaterf(id, 0, "invalid parcel id=%d", id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	storedParcels, err := store.GetByClient(client)
	fmt.Printf("storedParcels: len=%d, %v\n", len(storedParcels), storedParcels)

	suite.Require().NoErrorf(err, "unable to get parcel by client id=%id, err=%v", client, err)
	suite.Assert().Equal(len(storedParcels), len(parcelMap))

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		p, ok := parcelMap[parcel.Number]
		suite.Require().True(ok)
		// убедитесь, что значения полей полученных посылок заполнены верно
		suite.Assert().Equal(p.Client, parcel.Client)
		suite.Assert().Equal(p.Status, parcel.Status)
		suite.Assert().Equal(p.Address, parcel.Address)
		suite.Assert().Equal(p.CreatedAt, parcel.CreatedAt)
	}
}
