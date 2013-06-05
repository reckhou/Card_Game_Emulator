package main

import (
  "fmt"
  "math/rand"
  "time"
  "sort"
)

const (
  CLUB = 0
  DIAMOND = 1
  SPADE = 2
  HEART = 3
  CARD_MIN = 8
  CARD_MAX = 39 // 3-10
  DECK_SIZE = CARD_MAX - CARD_MIN + 1
  HAND_SIZE = 10
  MAX_DRAW_AMOUNT = 4
)

var gDeck []int
var gHand []int

func main() {
  round := 1
  for {
    fmt.Printf("====== 第%d轮 ======\n", round)
    initDeck()
    initHand()
    attack()
    round++
  }
  initDeck()
  printCardSet(gDeck, "Deck")
  initHand()


  for len(gDeck) !=0 {
    fmt.Printf("====== 第%d轮 ======\n", round)
    drawFromDeck(gHand, MAX_DRAW_AMOUNT)
    attack()
    round++
  }
}

func initDeck() () {
  gDeck = make([]int, DECK_SIZE)
  for i := 0; i != DECK_SIZE; i++ {
    gDeck[i] = CARD_MIN+i
  }
  shuffleDeck_Fisher_Yates(gDeck)
}

func initHand() {
  gHand = []int{}
  drawFromDeck(gHand, 10)
}

func drawFromDeck(hand []int, amount int) {
  if len(hand) + amount > HAND_SIZE {
    amount = HAND_SIZE - len(hand)
    fmt.Printf("Hand exceed! Draw %d instead.", amount)
  }

  if len(gDeck) < amount {
    amount = len(gDeck)
    fmt.Printf("Deck Empty! Draw %d instead.", amount)
  }

  drawn := gDeck[:amount]
  hand_new := append(hand, drawn...)
  //gHand = []int{12, 13, 14, 15, 33}
  gHand = hand_new
  gDeck = gDeck[amount:]
  sort.Sort(sort.IntSlice(gHand))
  printCardSet(gDeck, "Deck")
  printCardSet(gHand, "Hand")
}

func printCardSet(card_set []int, desc string) {
  fmt.Printf("\n%s(%d left):[", desc, len(card_set))
  for i, num := range card_set {
    if i%4 == 0 {
      fmt.Printf("\n")
    }
    var card_type string
    switch card_set[i]%4 {
      case 0:
        card_type = "草花"
      case 1:
        card_type = "方块"
      case 2:
        card_type = "黑桃"
      case 3:
        card_type = "红桃"
      default:
        card_type = "空"
    }
    fmt.Printf("%s%d(%d), ", card_type, num/4+1, num)
  }
  fmt.Printf("]\n")
}

func shuffleDeck_Fisher_Yates(deck []int) {
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
}

func attack() {
  // 基本策略：1. 优先级炸弹>同花顺>顺子>三带二>三带一>连对>对子>单张
  // 2. 点数越高越早出

  tmp := make([]int, len(gHand))
  copy(tmp, gHand)
  // 1. Split hand
  bomb := []int{}
  /*tonghua_line := []int{}
  line := []int{}*/
  triple := []int{}
  double := []int{}
  for i := 0; i != len(tmp); i++ {
    num := tmp[i]
    // find bomb
    if num%4 == 0 {
      if i+4 <= len(tmp) {
        if tmp[i+1] == num+1 && tmp[i+2] == num+2 && tmp[i+3] == num+3 {
          fmt.Printf("i:%d\n", i)
          bomb = append(bomb, tmp[i:i+4]...)
          tmp = append(tmp[:i], tmp[i+4:]...)
          printCardSet(bomb, "bomb")
          printCardSet(tmp, "tmp")
          i = -1
          continue
        }
      }
    }
  }

  for i := 0; i != len(tmp); i++ {
    // find triple
    if i+3 <= len(tmp) {
      if tmp[i+1]/4 == tmp[i]/4 && tmp[i+2]/4 == tmp[i]/4 {
        fmt.Printf("tmp[%d]: %d\n", i, tmp[i])
        triple = append(triple, tmp[i:i+3]...)
        tmp = append(tmp[:i], tmp[i+3:]...)
        printCardSet(triple, "triple")
        printCardSet(tmp, "tmp")
        i = -1
        continue
      }
    }
  }

  for i := 0; i != len(tmp); i++ {
    // find double
    if i+2 <= len(tmp) {
      if tmp[i+1]/4 == tmp[i]/4 {
        fmt.Printf("tmp[%d]: %d\n", i, tmp[i])
        double = append(double, tmp[i:i+2]...)
        tmp = append(tmp[:i], tmp[i+2:]...)
        printCardSet(double, "double")
        printCardSet(tmp, "tmp")
        i = -1
        continue
      }
    }
  }

  printCardSet(bomb, "Bomb")
  printCardSet(triple, "Triple")
  printCardSet(double, "Double")
  printCardSet(tmp, "Single")
}
