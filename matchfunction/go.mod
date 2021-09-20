module github.com/FairlySadPanda/open-match-but-kube-got-shreked/tutorials/matchmaker101/matchfunction

go 1.14

require (
	google.golang.org/grpc v1.36.0
	open-match.dev/open-match v0.0.0-dev
	github.com/google/uuid v1.3.0
)

replace open-match.dev/open-match v0.0.0-dev => ./the_entire_project/
