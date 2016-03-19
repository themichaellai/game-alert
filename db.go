package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"regexp"
)

type Game struct {
	Home, Away, Status, PositionId, Url string
	Id                                  int
}

type Event struct {
	Status     string
	GameId, Id int
	Datetime   int64
}

type EventGame struct {
	Event *Event
	Game  *Game
}

func GameResponseToGame(gameResponse *GameResponse) *Game {
	return &Game{
		Home:       gameResponse.Home.Names.Short,
		Away:       gameResponse.Away.Names.Short,
		Status:     parseStatus(gameResponse.FinalMessage),
		PositionId: gameResponse.BracketPositionId,
		Url:        gameResponse.Url,
	}
}

func parseStatus(rawStatus string) string {
	halfRe, _ := regexp.Compile(`^(\d\w\w)( Half)?`)
	matched := halfRe.MatchString(rawStatus)
	if matched {
		matches := halfRe.FindStringSubmatch(rawStatus)
		return fmt.Sprintf("%s Half", matches[1])
	} else {
		return rawStatus
	}
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
		&game.PositionId,
		&game.Url,
	)
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
	stmt, err := tx.Prepare("update games set status = ?, home = ?, away = ?, url = ? where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(game.Status, game.Home, game.Away, game.Url, game.Id)
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
	stmt, err := tx.Prepare("insert into games(home, away, status, positionId, url) values(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(game.Home, game.Away, game.Status, game.PositionId, game.Url)
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

func GetLatestEvents(db *sql.DB) ([]*EventGame, error) {
	const limit = 50
	rows, err := db.Query(fmt.Sprintf("select gameId, datetime, home, away, games.status, positionId, url from events inner join games on games.id = gameId order by datetime desc limit %d;", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventgames = make([]*EventGame, limit)
	i := 0
	for rows.Next() {
		var g Game
		var e Event
		err := rows.Scan(
			&g.Id,
			&e.Datetime,
			&g.Home,
			&g.Away,
			&g.Status,
			&g.PositionId,
			&g.Url,
		)
		eventgames[i] = &EventGame{
			Game:  &g,
			Event: &e}
		i++
		if err != nil {
			return nil, err
		}
	}
	return eventgames[:i], nil
}
