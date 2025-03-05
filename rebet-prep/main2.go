package main

import (
	"fmt"
	"sort"
)

// ServerStatus represents the current state of a game server
type ServerStatus struct {
	ID             string
	CurrentPlayers int
	MaxCapacity    int
	Region         string
	Latency        int // in milliseconds
}

// PlayerGroup represents a group of players wanting to join a server
type PlayerGroup struct {
	GroupID     string
	PlayerCount int
	RegionPref  string
	MaxLatency  int // maximum acceptable latency
}

// Your solution function signature
func findOptimalServer(servers []ServerStatus, playerGroup PlayerGroup) string {
	// for each server in servers, find the best server with least latency that (n)
	// has capacity
	// has latency match
	// region pref

	outputServers := []ServerStatus{}
	for _, serv := range servers {
		spaceLeft := serv.MaxCapacity - serv.CurrentPlayers
		if spaceLeft > playerGroup.PlayerCount && serv.Latency <= playerGroup.MaxLatency && serv.Region == playerGroup.RegionPref {
			outputServers = append(outputServers, serv)
		}
	}

	if len(outputServers) == 0 {
		return "no_available_server"
	}
	sort.Slice(outputServers, func(i, j int) bool {
		iLoad := (outputServers[i].CurrentPlayers + playerGroup.PlayerCount) / outputServers[i].MaxCapacity
		jLoad := (outputServers[j].CurrentPlayers + playerGroup.PlayerCount) / outputServers[j].MaxCapacity
		return iLoad < jLoad
	})

	return outputServers[0].ID
}

// func main() {
// 	// Example data
// 	servers := []ServerStatus{
// 		{ID: "server-1", CurrentPlayers: 50, MaxCapacity: 100, Region: "us-west", Latency: 30},
// 		{ID: "server-2", CurrentPlayers: 80, MaxCapacity: 100, Region: "us-east", Latency: 45},
// 		{ID: "server-3", CurrentPlayers: 40, MaxCapacity: 80, Region: "us-west", Latency: 25},
// 		{ID: "server-4", CurrentPlayers: 90, MaxCapacity: 100, Region: "eu-central", Latency: 120},
// 		{ID: "server-5", CurrentPlayers: 20, MaxCapacity: 50, Region: "asia-east", Latency: 200},
// 	}
//
// 	playerGroup := PlayerGroup{
// 		GroupID:     "group-1",
// 		PlayerCount: 15,
// 		RegionPref:  "us-west",
// 		MaxLatency:  50,
// 	}
//
// 	optimalServerID := findOptimalServer(servers, playerGroup)
// 	fmt.Printf("Optimal server for group %s: %s\n", playerGroup.GroupID, optimalServerID)
// }
