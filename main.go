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
  HAND_SIZE = 15
  MAX_DRAW_AMOUNT = 2
  TOTAL_BATTLES = 10000
)

var gDeck []int
var gHand []int
var gDamage int

func main() {
  round := 1
  battle := 1
  gDamage = 0
  totalDamage := 0
  maxDamage := 0
  minDamage := 10000
  totalRounds := 0
  maxRound := 0
  minRound := 100
  /*for {
    fmt.Printf("====== 第%d轮 ======\n", round)
    initDeck()
    initHand()
    attack()
    return
    round++
  }*/
  initDeck()
  printCardSet(gDeck, "Deck")
  initHand()

  for battle <= TOTAL_BATTLES {
    initDeck()
    printCardSet(gDeck, "Deck")
    initHand()
    gDamage = 0
    round = 1
    for len(gDeck) !=0 && gDamage < 1500 {
      fmt.Printf("====== 第%d场第%d轮 ======\n", battle, round)
      drawFromDeck(gHand, MAX_DRAW_AMOUNT)
      attack()
      round++
    }
    fmt.Printf("====== 本场战斗总伤害:%d ======\n", gDamage)
    totalDamage += gDamage
    if gDamage > maxDamage {
      maxDamage = gDamage
    }
    if gDamage < minDamage {
      minDamage = gDamage
    }

    if round > maxRound {
      maxRound = round
    }
    if round < minRound {
      minRound = round
    }

    battle++
    totalRounds += round
  }

  fmt.Printf("====== 每场战斗平均伤害:%d MAX:%d MIN:%d ======\n====== 每场战斗平均经过回合:%d MAX:%d MIN:%d ======\n", totalDamage/TOTAL_BATTLES, maxDamage, minDamage, totalRounds/TOTAL_BATTLES, maxRound, minRound)
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
  return
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
    j := r.Int() % (card_cnt);
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
          bomb = append(bomb, tmp[i:i+4]...)
          tmp = append(tmp[:i], tmp[i+4:]...)
          //printCardSet(bomb, "bomb")
          //printCardSet(tmp, "tmp")
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
        triple = append(triple, tmp[i:i+3]...)
        tmp = append(tmp[:i], tmp[i+3:]...)
        //printCardSet(triple, "triple")
        //printCardSet(tmp, "tmp")
        i = -1
        continue
      }
    }
  }

  for i := 0; i != len(tmp); i++ {
    // find double
    if i+2 <= len(tmp) {
      if tmp[i+1]/4 == tmp[i]/4 {
        double = append(double, tmp[i:i+2]...)
        tmp = append(tmp[:i], tmp[i+2:]...)
        //printCardSet(double, "double")
        //printCardSet(tmp, "tmp")
        i = -1
        continue
      }
    }
  }

  printCardSet(bomb, "Bomb")
  printCardSet(triple, "Triple")
  printCardSet(double, "Double")
  printCardSet(tmp, "Single")
  
  attack_card := []int{}
  
  for {
    if (len(bomb) != 0) {
      start_idx := len(bomb)-4
      attack_card = append(attack_card, bomb[start_idx:]...)
      break
    }

    if (len(triple) != 0) {
      start_idx := len(triple)-3
      attack_card = append(attack_card, triple[start_idx:]...)
      break
    }

    if (len(double) != 0) {
      start_idx := len(double)-2
      attack_card = append(attack_card, double[start_idx:]...)
      break
    }

    if (len(tmp) != 0) {
      start_idx := len(tmp)-1
      attack_card = append(attack_card, tmp[start_idx:]...)
      break
    }
    break
  }

  printCardSet(attack_card, "Attack_card")
  calculateDamage(attack_card)
  
  for j, val := range gHand {
    if val == attack_card[0] {
      gHand = append(gHand[:j], gHand[j+len(attack_card):]...)
    }
  }
}

func calculateDamage(attack_card []int) {
  sum := 0
  count := 0
  for _, val := range attack_card {
    sum += val
    count++
  }
  fmt.Printf("\nDamage: %d\n", sum * count)
  gDamage += sum * count
}
