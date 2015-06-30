FROM ubuntu:14.04

# This part is taken from the official docker image --------------------

RUN apt-get update && apt-get install -y \
		build-essential ca-certificates curl \
		bzr git mercurial openssh-client\
		--no-install-recommends

ENV GOLANG_VERSION 1.3

RUN curl -sSL http://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
	| tar -v -C /usr/src -xz
WORKDIR /usr/src/go

RUN cd src && ./make.bash --no-clean 2>&1

ENV PATH /usr/src/go/bin:$PATH

RUN mkdir -p /go/src
ENV GOPATH /go
ENV PATH /go/bin:$PATH
WORKDIR /go

# ----------------------------------------------------------------------


RUN useradd -m git

ENV GITEA_PATH $GOPATH/src/github.com/go-gitea/gitea
ENV GITEA_CUSTOM_CONF_PATH $GITEA_PATH/custom/conf
ENV GITEA_CUSTOM_CONF $GITEA_CUSTOM_CONF_PATH/app.ini

RUN go get -u -d github.com/go-gitea/gitea
# WORKDIR $GITEA_PATH
WORKDIR /go/src/github.com/go-gitea/gitea
RUN go build github.com/go-gitea/gitea
RUN chown -R git $GITEA_PATH

ADD init_gitea.sh /tmp/
RUN chown git /tmp/init_gitea.sh
RUN chmod +x /tmp/init_gitea.sh

USER git
ENV HOME /home/git
ENV USER git
ENV PATH $GITEA_PATH:$PATH

RUN git config --global user.name "Gitea" && git config --global user.email "gitea@gitea.io"

ENTRYPOINT ["/tmp/init_gitea.sh"]
CMD ["gitea", "web"]
