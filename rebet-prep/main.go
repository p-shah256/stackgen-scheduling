/*
Here's your task: I have a stream of player activity events that look like this:
{
  "player_id": "p123",
  "event_type": "item_purchase",
  "timestamp": 1646161776,
  "data": {
    "item_id": "sword_42",
    "currency_type": "gems",
    "amount": 500
  }
}

These events come in as a batch, and I need you to:

1. Group them by player_id
2. For each player, calculate their total spend across all currencies
3. Identify which item each player spent the most on
4. Return a transformed data structure with this summarized information
*/

package main

import "time"

// INPUT
type playerActivity struct {
	playerId   string
	event_type string
	timestamp  time.Time
	data       dataType
}

type dataType struct {
	item_id       string
	currency_type string
	amount        int
}

// OUTPUT
// {pid1: {total_spend: 8082, itemMost: "gems"}, pid2: {}, pid3: {}}
type playerInfoExternal struct {
	totalSpent int
	itemMost   string
	itemSpent  int
}

// {
//
//	  pid1: playerInfoInternal{
//	    itemSpendLog: {"swords": 800, "gems": 280}
//	    externalInfo: playerInfoExternal{
//	      totalSpent: 8082,
//	      itemMost: "gems",
//	      itemSpent: 900,
//	    },
//	  },
//	  pid2: playerInfoInternal{},
//	}
type playerInfoInternal struct {
	itemSpendLog map[string]int
	exInfo       playerInfoExternal
}

func transfrom(activities []playerActivity) (map[string]playerInfoExternal, error) {
	internal := make(map[string]playerInfoInternal)
	external := make(map[string]playerInfoExternal)

	for _, activity := range activities {
		itemId := activity.data.item_id
		if _, ok := internal[activity.playerId]; !ok {
			internal[activity.playerId] = playerInfoInternal{}
		}
		if _, ok := external[activity.playerId]; !ok {
			external[activity.playerId] = playerInfoExternal{}
		}

		playerIn := internal[activity.playerId]

		if _, ok := playerIn.itemSpendLog[itemId]; !ok {
      playerIn.itemSpendLog[itemId] = 0
		}
		playerIn.itemSpendLog[itemId] += activity.data.amount
		playerIn.exInfo.totalSpent += activity.data.amount

		// update external info
		if playerIn.itemSpendLog[itemId] > playerIn.exInfo.itemSpent {
      // update internal 
			playerIn.exInfo.itemMost = itemId
			playerIn.exInfo.itemSpent = playerIn.itemSpendLog[itemId]
		}

    external[activity.playerId] = playerIn.exInfo
	}
	return external, nil
}
