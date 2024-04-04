package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))
	err := row.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		parcel := Parcel{}
		err := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, parcel)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	parsel := Parcel{}
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&parsel.Status)
	if err != nil {
		return err
	}
	if parsel.Status == "registered" {
		_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number)) 
	if err != nil {
		return err
	}
	return nil
}
