FROM scratch
MAINTAINER Vincent Ambo <dev@tazj.in>

ADD etcd-announcer /etcd-announcer

CMD []
ENTRYPOINT ["/etcd-announcer"]
