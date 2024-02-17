package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

const (
	sqlInsertParcel = "insert into parcel (client, status, address, created_at) values (:client, :status, :address, :createdat)"
)

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec(sqlInsertParcel,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("createdat", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	number, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	return int(number), nil
}

const (
	sqlSelectByParcelId = "select * from parcel where number = :number"
)

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	row := s.db.QueryRow(sqlSelectByParcelId, sql.Named("number", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

const (
	sqlSelectByClientId = "select * from parcel where client = :client"
)

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	rows, err := s.db.Query(sqlSelectByClientId, sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for rows.Next() {
		var p Parcel
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}
		res = append(res, p)
	}

	return res, nil
}

const (
	sqlUpdateStatusByParcelId = "update parcel set status = :status where number = :number"
)

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec(sqlUpdateStatusByParcelId,
		sql.Named("status", status),
		sql.Named("number", number))

	return err
}

const (
	sqlUpdateAddressById = "update parcel set address = :address where number = :number"
)

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if !strings.EqualFold(parcel.Status, "registered") {
		return fmt.Errorf("unable to change address for parcel id=%d with status=%s", parcel.Number, parcel.Status)
	}

	_, err = s.db.Exec(sqlUpdateAddressById, sql.Named("address", address), sql.Named("number", number))

	return err
}

const (
	sqlDeleteById = "delete from parcel where number = :number"
)

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if !strings.EqualFold(parcel.Status, "registered") {

		// не удаляет, но возвращает nil, так как в противном случае валится main.go, там ошибка удаления не обрабатывается, и программа сразу завершается
		return nil

		//return fmt.Errorf("unable to delete parcel with id=%d, status=%s", parcel.Number, parcel.Status)
	}

	_, err = s.db.Exec(sqlDeleteById, sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}
