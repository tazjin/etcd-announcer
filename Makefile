cntrpath = /gopath/src/etcd-announcer

all: gobuild docker

gobuild:
	docker run --rm -t -v $(CURDIR):$(cntrpath) -w $(cntrpath) \
	google/golang sh $(cntrpath)/build.sh

docker: gobuild
	docker build -t tazjin/etcd-announcer $(CURDIR)

clean:
	docker run --rm -t -v $(CURDIR):/gopath/src/etcd-announcer \
	 -w /gopath/src/etcd-announcer google/golang go clean
