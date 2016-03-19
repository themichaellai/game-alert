package main

import (
	"database/sql"
	"log"
	"time"
)

func createEvent(game *Game, db *sql.DB) (*Event, error) {
	event := &Event{
		GameId:   game.Id,
		Datetime: time.Now().Unix(),
		Status:   game.Status,
	}
	err := InsertEvent(event, db)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func runScrape(dbFilename string) {
	log.Println("scraping...")
	db := DBHandle(dbFilename)
	defer db.Close()
	gameResponses, err := Scrape("http://data.ncaa.com/carmen/brackets/championships/basketball-men/d1/2016/data.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, gameResponse := range gameResponses {
		game := GameResponseToGame(gameResponse)
		existingGame, err := SelectByPositionId(game.PositionId, db)
		if err != nil { // game does not exist, insert
			log.Printf("game %s does not exist\n", game.PositionId)
			err = InsertGame(game, db)
			if err != nil {
				log.Fatal(err)
			}
		} else { // game exists, update
			if existingGame.Status != game.Status {
				log.Printf("game %d status setting to %s\n", existingGame.Id, game.Status)
				existingGame.Status = game.Status
				existingGame.Home = game.Home
				existingGame.Away = game.Away
				err := UpdateGame(existingGame, db)
				if err != nil {
					log.Fatal(err)
				}

				_, err = createEvent(existingGame, db)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func main() {
	const dbFilename = "./games.db"
	ticker := time.NewTicker(time.Minute)
	go func() {
		for _ = range ticker.C {
			runScrape(dbFilename)
		}
	}()
	ServeFeed(dbFilename)
}
