package main

import (
  "fmt"
  "math/rand"
  "time"
)

const (
  CLUB = 0
  DIAMOND = 1
  SPADE = 2
  HEART = 3
  CARD_MIN = 8
  CARD_MAX = 39 // 3-10
)

func main() {
  cards := initDeck()
  printDeck(cards)
}

func initDeck() (cards []int) {
  cards = make([]int, 32)
  for i := 0; i <=CARD_MAX-CARD_MIN; i++ {
    cards[i] = CARD_MIN+i
  }
  cards = shuffleDeck_Fisher_Yates(cards)
  return cards
}

func printDeck(deck []int) {
  fmt.Printf("Deck:[")
  for i, num := range deck {
    if i%4 == 0 {
      fmt.Printf("\n")
    }
    var card_type string
    switch deck[i]%4 {
      case 0:
        card_type = "C"
      case 1:
        card_type = "D"
      case 2:
        card_type = "S"
      case 3:
        card_type = "H"
    }
    fmt.Printf("%s%d, ", card_type, num/4+1)
  }
  fmt.Printf("]\n")
}

func shuffleDeck_Fisher_Yates(deck []int) (result []int) {
  card_cnt := len(deck)
  var temp int
  if card_cnt == 0 {
    return
  }
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  for ;card_cnt !=0; card_cnt-- {
    j := r.Int() % (card_cnt+1);
    temp = deck[card_cnt-1];
    deck[card_cnt-1] = deck[j];
    deck[j] = temp;
  }
  result = deck
  return result
}
