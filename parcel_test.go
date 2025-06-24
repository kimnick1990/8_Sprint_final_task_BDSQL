package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "file:tracker.db?cache=shared")
	require.NoError(t, err)
	defer db.Close() // настройте подключение к БД
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0) // add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	gotParcel, err := store.Get(id)
	require.NoError(t, err)

	require.Equal(t, parcel.Client, gotParcel.Client)
	require.Equal(t, parcel.Status, gotParcel.Status)
	require.Equal(t, parcel.Address, gotParcel.Address)
	require.Equal(t, parcel.CreatedAt, gotParcel.CreatedAt)

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err) // delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "file:tracker.db?cache=shared")
	require.NoError(t, err)
	defer db.Close() // настройте подключение к БД

	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	gotParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, gotParcel.Address) // check
	// получите добавленную посылку и убедитесь, что адрес обновился
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "file:tracker.db?cache=shared")
	require.NoError(t, err)
	defer db.Close() // настройте подключение к БД
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	newStatus := ParcelStatusDelivered
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err) // set status
	// обновите статус, убедитесь в отсутствии ошибки

	gotParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newStatus, gotParcel.Status) // check
	// получите добавленную посылку и убедитесь, что статус обновился
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "file:tracker.db?cache=shared")
	require.NoError(t, err)
	defer db.Close() // настройте подключение к БД
	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	for _, parcel := range parcels {
		id, err := store.Add(parcel)
		require.NoError(t, err)
		require.Greater(t, id, 0)

		parcel.Number = id
		parcelMap[id] = parcel
	} // add
	//for i := 0; i < len(parcels); i++ {
	//id, err := // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	// обновляем идентификатор добавленной у посылки
	//parcels[i].Number = id

	// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
	//	parcelMap[id] = parcels[i]

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels)) // получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных

	// check
	for _, parcel := range storedParcels {
		gotParcel, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		require.Equal(t, gotParcel, parcel)
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
	}

}
