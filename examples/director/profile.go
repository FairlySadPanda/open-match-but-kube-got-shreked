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
	"github.com/FairlySadPanda/open-match-but-kube-got-shreked/examples"
	"github.com/FairlySadPanda/open-match-but-kube-got-shreked/pkg/pb"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/types/known/anypb"
)

// generateProfiles generates test profiles.
func generateProfiles() []*pb.MatchProfile {
	var profiles []*pb.MatchProfile
	modes := []string{examples.EightBall.ToString(), examples.NineBall.ToString(), examples.KoreanCarom.ToString(), examples.JapaneseCarom.ToString()}
	regions := []string{examples.Africa.ToString(), examples.Asia.ToString(), examples.Australasia.ToString(), examples.Europe.ToString(), examples.NorthAmerica.ToString(), examples.SouthAmerica.ToString()}
	extension := &any.Any{}
	extension.MarshalFrom(&examples.EloRange{
		Range: 200,
	})

	for _, mode := range modes {
		for _, region := range regions {
			profiles = append(profiles, &pb.MatchProfile{
				Name: "mode_based_profile",
				Pools: []*pb.Pool{
					{
						Name: mode,
						TagPresentFilters: []*pb.TagPresentFilter{
							{
								Tag: mode,
							},
							{
								Tag: region,
							},
						},
					},
				},
				Extensions: map[string]*anypb.Any{"Elo": extension},
			},
			)
		}

	}

	return profiles
}
