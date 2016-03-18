package main

import (
	"log"
)

func main() {
	db := DBHandle("./games.db")
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
				UpdateGame(existingGame, db)
			}
		}
	}
}
