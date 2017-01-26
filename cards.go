package main
import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	
)


/*
Various structs we are using. 
*/

type ShuffleResp struct {
	Success bool `json:"success"`
	DeckID string `json:"deck_id"`
	Shuffled bool `json:"shuffled"`
	Remaining int `json:"remaining"`
}


type DrawResp struct {
	Success bool `json:"success"`
	Cards []struct {
		Image string `json:"image"`
		Value string `json:"value"`
		Suit string
		Code string
	} `json:"cards"`
	DeckID string `json:"deck_id"`
	Remaining int `json:"remaining"`
}

type CurHand struct {
	Value string
}

/*
Generic bit of code we use to simulate Curl GET requests to any particular URL. 
*/

func curlJson(url string) (stuff []byte){
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	} else {
		defer r.Body.Close()
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		stuff :=  contents
		return stuff
	}
}


/*
Create a new deck. 
*/

func newDeck()(deckid string) {
	jsonstring := curlJson("https://deckofcardsapi.com/api/deck/new/shuffle/?deck_count=1")
	var new ShuffleResp
	err := json.Unmarshal(jsonstring, &new)
	if err != nil {
		panic(err)
	}
	if new.Success != true {
		panic("Shuffling did not return true")
	} else {
		deckid = new.DeckID
	}
	return 
}

/*
Draw from existing deck, provided you have an ID, return value of card drawn. 
*/

func drawFromExisting(deckID string) (value string) {
	createurl := "https://deckofcardsapi.com/api/deck/" + deckID + "/draw/?count=1"
	var hand DrawResp
	jsonstring := curlJson(createurl)
	err := json.Unmarshal(jsonstring, &hand)
	if err != nil {
		panic(err)
	}
	var curdrawval string 
	for _, element := range hand.Cards {
		curdrawval = element.Value
	}
	return curdrawval
}

/*
Reshuffle existing deck
*/
func reshuffleExisting(deckID string){
	createurl := "https://deckofcardsapi.com/api/deck/" + deckID + "/shuffle/"
	jsonstring := curlJson(createurl)
	var old ShuffleResp
	err := json.Unmarshal(jsonstring, &old)
	if err != nil {
		panic(err)
	}
	if old.Success != true {
		panic("shuffling did not return true")
	}
	return 
}



func main() {
/*
 Hand map to keep a running tally       
*/
	hand := make(map[CurHand]int)
/*
 Declare deckid here so we can update it from within the main game loop 
*/
	var deckid string
/*
Main game loop
*/
	for {
/*
Rudimentary input handling
*/
		
		var input string
		fmt.Println(`
Enter one of the following
[1] start new game
[2] draw a card
[3] reshuffle hand into deck
`)
		fmt.Scanln(&input)
/*
Horrible monstrous case statement that reads the input and does things that you select. 
*/
		switch  {
		case input  == "1":
			deckid = newDeck()
			fmt.Println("New game started: your current Deck ID is:", deckid)
		case input == "2":
/*
Error checking just in case you selected an option without starting a game first. 
*/
			if deckid == "" {
				fmt.Println("No deck selected, please start a new game")
			} else {
				value := drawFromExisting(deckid)
				hand[CurHand{value}]++
				fmt.Println("You drew:", value)
				fmt.Println("Current hand:")
				for k := range hand {
					fmt.Println(k.Value + " : " + strconv.Itoa(hand[k]))
				}
			}
		case input == "3":
			if deckid == "" {
				fmt.Println("No deck selected, please start a new game")
			} else {
				reshuffleExisting(deckid)
				hand = make(map[CurHand]int)
			}
		default:
			fmt.Println("I didn't catch that, please enter a valid option")
		}
	}
}

/*
Notes: 
No bonuses attempted due to the speed with which this was written. 
This is simultaneously severely over and under-engineered
*/
