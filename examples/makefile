protogen_dir_go=/tmp/protogen

protoc-go-clean:
	rm -rf ${protogen_dir_go}
	find -iname '*.pb.go' -exec rm {} \;

protoc-go: protoc-go-clean
	mkdir -p ${protogen_dir_go}/github.com/FairlySadPanda
	ln -s $(shell pwd) ${protogen_dir_go}/github.com/FairlySadPanda/open-match-but-kube-got-shreked
	find -iname '*.proto' -exec protoc --proto_path=. --go_out=${protogen_dir_go}/ --go-grpc_out=${protogen_dir_go}/ {} \;