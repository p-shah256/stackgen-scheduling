package main

// If a player completes 5 matches in an hour, send a "dedication" notification
// If a player levels up twice in a day, send a "fast learner" notification
// If a player unlocks 3 achievements in a session, send a "on fire" notification

// {
//   "event_id": "ev-12345",
//   "event_type": "match_complete",
//   "timestamp": 1709575342,
//   "payload": {
//     // Various fields depending on event type
//   }
// }

type eventInfo struct {
	eventId   string
	playerId  string
	eventType string
	timestamp uint64
	payload   map[string]string // for now
}

type players map[string]playerMetadata
type playerMetadata struct {
	// playerID : {
	//	eventLog: {
	//			eventID_1: [{eventId, playerId, ...}],
	//			eventID_2: [{eventId, playerId, ...}],
	//	},
	//	badgeInfo: {bagdewindowstart: 213123, badgeRecords: 1},
	// }
	// eventType_2: [{eventId, playerId, ...}]}
	eventLog map[string][]eventInfo
	// {badge1: {badgeName, windowstart, windowend, record...}, badge2: {}}
	badgeInfo map[string]badge
}

type badge struct {
	name        string
	windowStart uint64 // check if badgeWindowStart < timestamp < badgeWindowStart + badgeWindowSize
	records     int
	windowSize  uint64
	gained      bool
}

/**
writing a brute force approach first
1. get the playerId from each event
2. update the players state

a. check if the badge window was still open and the no.ofEvents > badge requirement
b. notify
*/

// keep it realtime with something like webhooks instead of async to just processEvent()
// this will change later, will get playerinfo from db instead
func processEvent(event eventInfo, pl players) string {
	// Ensure player exists in map
	if _, ok := pl[event.playerId]; !ok {
		pl[event.playerId] = playerMetadata{
			eventLog:  make(map[string][]eventInfo),
			badgeInfo: make(map[string]*badge), // Store badge pointers
		}
	}

	// Append event to log
	pl[event.playerId].eventLog[event.eventType] = append(pl[event.playerId].eventLog[event.eventType], event)

	// Check badge conditions for "match_complete"
	if event.eventType == "match_complete" {
		// Get badge map
		badgeInfo := pl[event.playerId].badgeInfo

		// If "dedication" badge does not exist, initialize it
		if _, exists := badgeInfo["dedication"]; !exists {
			badgeInfo["dedication"] = &badge{
				windowStart: event.timestamp,
				windowSize:  60, // mins TODO: needs converting
				records:     0,
			}
		}

		// Directly modify the badge struct via pointer
		dedicationBadge := badgeInfo["dedication"]
		dedicationBadge.records++

		// Check if badge is gained
		if dedicationBadge.records >= 5 {
			dedicationBadge.gained = true
			return "dedication"
		}
	}

	return ""
}

