// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"math/rand"
	"time"

	"github.com/FairlySadPanda/open-match-but-kube-got-shreked/examples"
	"github.com/FairlySadPanda/open-match-but-kube-got-shreked/pkg/pb"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/types/known/anypb"
)

// Ticket generates a Ticket with a mode search field that has one of the
// randomly selected modes.
func makeTicket() *pb.Ticket {
	tags := []string{examples.Region(rand.Intn(4)).ToString()}
	tags = append(tags, gameModes()...)

	extension := &any.Any{}
	extension.MarshalFrom(&examples.MatchUsingElo{
		Elo: float32(generateRandomNormalizedElo()),
	})

	ticket := &pb.Ticket{
		// A ticket has an ID
		// An assignment if it's assigned to a game server
		// assignments are a string and a metadata map
		// there's a metadata "Extensions" map
		// SearchFields provides the search fields
		SearchFields: &pb.SearchFields{
			// StringArgs is a map of string to string that matches on equality
			// tags is a search critera that is just string tags
			// search criteria for doubles - assessed via ranges
			DoubleArgs: map[string]float64{
				examples.TimeEnteringQueueKey: enterQueueTime(),
			},
			Tags: tags,
		},
		Extensions: map[string]*anypb.Any{"Elo": extension},
	}

	return ticket
}

// gameModes selects up to two unique game modes for this Ticket.
func gameModes() []string {
	if rand.Intn(2) == 1 {
		return []string{examples.GameMode(rand.Intn(4)).ToString()}
	}

	modes := map[string]struct{}{}
	keys := make([]string, 2)

	for i := 0; i < 2; {
		mode := examples.GameMode(rand.Intn(4)).ToString()
		if _, ok := modes[mode]; !ok {
			modes[mode] = struct{}{}
			keys[i] = mode
			i++
		}
	}

	return keys
}

// enterQueueTime returns the time at which the Ticket entered matchmaking queue. To simulate
// variability, we pick any random second interval in the past 5 seconds as the players queue
// entry time.
func enterQueueTime() float64 {
	return float64(time.Now().Add(-time.Duration(rand.Intn(5)) * time.Second).UnixNano())
}

// generate a believable range of Elo values based on the mean value - the average player - being rated 1500.
func generateRandomNormalizedElo() float64 {
	return rand.NormFloat64()*350 + 1500
}
