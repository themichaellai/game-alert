package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Game struct {
	Home, Away, Status, PositionId string
	Id                             int
}

type Event struct {
	Status     string
	GameId, Id int
	Datetime   int64
}

func GameResponseToGame(gameResponse *GameResponse) *Game {
	return &Game{
		Home:       gameResponse.Home.Names.Short,
		Away:       gameResponse.Away.Names.Short,
		Status:     gameResponse.FinalMessage,
		PositionId: gameResponse.BracketPositionId}
}

func SelectByPositionId(positionId string, db *sql.DB) (*Game, error) {
	stmt, err := db.Prepare("select * from games where positionId = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var game Game
	err = stmt.QueryRow(positionId).Scan(
		&game.Id,
		&game.Home,
		&game.Away,
		&game.Status,
		&game.PositionId)
	return &game, err
}

func UpdateGame(game *Game, db *sql.DB) error {
	if game.Id == 0 {
		return errors.New("need a game id to update game")
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("update games set status = ?, home = ?, away = ? where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(game.Status, game.Home, game.Away, game.Id)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func InsertGame(game *Game, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into games(home, away, status, positionId) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(game.Home, game.Away, game.Status, game.PositionId)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func InsertEvent(event *Event, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into events(gameId, datetime, status) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(event.GameId, event.Datetime, event.Status)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func DBHandle(dbFilename string) *sql.DB {
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
