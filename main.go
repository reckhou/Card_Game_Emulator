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
  SAME_CARD_COUNT = 1
  DECK_SIZE = (CARD_MAX - CARD_MIN + 1)*SAME_CARD_COUNT
  HAND_SIZE = 12
  MAX_DRAW_AMOUNT = 2
  TOTAL_BATTLES = 100000
)

var gDeck []int
var gHand []int
var gDamage int
var gBomb int
var gTonghuaLine int
var gTriple int
var gDouble int
var gSingle int

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
  gBomb = 0
  gTonghuaLine = 0
  gTriple = 0
  gDouble = 0
  gSingle = 0
  round_stat := make([]int, 20)
  
  for battle <= TOTAL_BATTLES {
    initDeck()
    printCardSet(gDeck, "Deck")
    initHand()
    gDamage = 0
    round = 1
    for len(gDeck) !=0 && gDamage < 1600 {
      fmt.Printf("====== 第%d场第%d轮 ======\n", battle, round)
      drawFromDeck(gHand, MAX_DRAW_AMOUNT)
      //gHand = []int{9, 13, 13, 17, 21, 25, 27, 29}
      //printCardSet(gHand, "Hand_debug")
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
    round_stat[round]++
  }

  fmt.Printf("====== 每场战斗平均伤害:%d MAX:%d MIN:%d ======\n====== 每场战斗平均经过回合:%d MAX:%d MIN:%d ======\n", totalDamage/TOTAL_BATTLES, maxDamage, minDamage, totalRounds/TOTAL_BATTLES, maxRound, minRound)
  fmt.Printf("====== 回合数分布: ======\n")
  fmt.Println(round_stat)
  round_stat_percent := make([]float32, len(round_stat))
  for i, val := range round_stat {
    var percent float32
    percent = float32(val)/float32(TOTAL_BATTLES)
    round_stat_percent[i] = percent * 100.0
  }
  fmt.Println(round_stat_percent)

  fmt.Printf("====== 总回合数:%d ======\n", totalRounds)
  fmt.Printf("====== 各种组合百分比: ======\n")
  fmt.Printf("====== 炸弹: %f 同花顺:%f 顺子:%f 三带一:%f 三带二:%f 连对:%f 对子: %f 单张: %f ======\n", float32(gBomb)/float32(totalRounds)*100, float32(gTonghuaLine)/float32(totalRounds)*100, 0, float32(gTriple)/float32(totalRounds)*100, 0, 0, float32(gDouble)/float32(totalRounds)*100, float32(gSingle)/float32(totalRounds)*100)
}

func initDeck() () {
  gDeck = make([]int, DECK_SIZE)
  for i := 0; i != DECK_SIZE; i++ {
    gDeck[i] = CARD_MIN+i/SAME_CARD_COUNT
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
    //fmt.Printf("Hand exceed! Draw %d instead.", amount)
  }

  if len(gDeck) < amount {
    amount = len(gDeck)
    //fmt.Printf("Deck Empty! Draw %d instead.", amount)
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
  if desc != "Attack_card" {
    return
  }
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
  r := rand.New(rand.NewSource(time.Now().UnixNano()*time.Now().Unix()))
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
  tonghua_line := []int{}
  //line := []int{}
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

  // find line & tonghua_line
  // 1. made an array without same card
  tmp_diffcard := []int{}
  for i := 0; i != len(tmp_diffcard); i++ {
    tmp_diffcard[i] = -1
  }
  j := 0
  for i := 0; i != len(tmp); i++ {
      if i > 0 {
        if tmp_diffcard[j-1] != tmp[i] {
          tmp_diffcard = append(tmp_diffcard, tmp[i])
          j++
        } else {
          continue
        }
      } else {
        tmp_diffcard = append(tmp_diffcard, tmp[i])
        j++
      }
    }
  printCardSet(tmp_diffcard, "Diffcard")

  // 2. made array with same color to find tonghua_line
  tmp_club := []int{}
  tmp_diamond := []int{}
  tmp_spade := []int{}
  tmp_heart := []int{}
  c := 0
  d := 0
  s := 0
  h := 0
  for i := 0; i != len(tmp_diffcard); i++ {
    switch tmp_diffcard[i]%4 {
    case CLUB:
      tmp_club = append(tmp_club, tmp_diffcard[i])
      c++
    case DIAMOND:
      tmp_diamond = append(tmp_diamond, tmp_diffcard[i])
      d++
    case SPADE:
      tmp_spade = append(tmp_spade, tmp_diffcard[i])
      s++
    case HEART:
      tmp_heart = append(tmp_heart, tmp_diffcard[i])
      h++
    }
  }
  printCardSet(tmp_club, "tmp_club")
  printCardSet(tmp_diamond, "tmp_diamond")
  printCardSet(tmp_spade, "tmp_spade")
  printCardSet(tmp_heart, "tmp_heart")
  findTonghuaLine(tmp_club, tonghua_line)
  findTonghuaLine(tmp_heart, tonghua_line)
  findTonghuaLine(tmp_spade, tonghua_line)
  findTonghuaLine(tmp_diamond, tonghua_line)

  printCardSet(tonghua_line, "Tonghua")

  // remove tonghua cards from tmp_diffcard & tmp
  tmp_diffcard_after := make([]int, len(tmp_diffcard))
  copy(tmp_diffcard_after, tmp_diffcard[0:])
  for i := 0; i < len(tonghua_line); i++ {
    for j := 0; j < len(tmp_diffcard_after); j++ {
      if tonghua_line[i] == tmp_diffcard_after[j] {
        if j < len(tmp_diffcard_after) {
          tmp_diffcard_after = append(tmp_diffcard_after[:j], tmp_diffcard_after[j+1:]...)
        } 
        break
      }
    }
  }

  for i := 0; i < len(tonghua_line); i++ {
    for j := 0; j < len(tmp); j++ {
      if tonghua_line[i] == tmp[j] {
        if j < len(tmp) {
          tmp = append(tmp[:j], tmp[j+1:]...)
        } 
        break
      }
    }
  }

  printCardSet(tmp_diffcard_after, "Diffcard_after")


  
  // find cards with line
  /*for i := 0; i != len(tmp_diffcard)-5; i++ {
    
  }*/

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
      gBomb++
      start_idx := len(bomb)-4
      attack_card = append(attack_card, bomb[start_idx:]...)
      break
    }

    if (len(tonghua_line) != 0) {
      gTonghuaLine++
      attack_card = append(attack_card, tonghua_line[0:]...)
      break
    }

    if (len(triple) != 0) {
      gTriple++
      start_idx := len(triple)-3
      attack_card = append(attack_card, triple[start_idx:]...)
      break
    }

    if (len(double) != 0) {
      gDouble++
      start_idx := len(double)-2
      attack_card = append(attack_card, double[start_idx:]...)
      break
    }

    if (len(tmp) != 0) {
      gSingle++
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
      break
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

func findTonghuaLine(target, add_to []int) {
  if len(target) >= 5 {
    j := 0
    for i := 0; i < len(target)-5-j; i++ {
      if target[i] == target[i+1]-4 && target[i] == target[i+2]-8 && target[i] == target[i+3]-12 && target[i] == target[i+4]-16 {
        // find line, check if further combo exists
        for j < len(target)-5-i {
          if target[i+j]-4*j == target[i] {
            j++
          } else {
            break
          }
        }
        add_to = append(add_to, target[i:i+5+j]...)
      }
    }
  }

}
