# Get base centos image
FROM centos:latest

#Install epel-release, debug tools, nginx, supervisor, boto packages.
#Create required directories
RUN \
  yum install -y epel-release bison python-setuptools bzip2 wget make gcc git&& \
  yum clean all && \
  rm -f /etc/localtime && \
  ln -s /usr/share/zoneinfo/Asia/Kolkata /etc/localtime

#Install Go
RUN \
  cd /tmp && \
  wget https://storage.googleapis.com/golang/go1.12.6.linux-amd64.tar.gz && \
  tar -C /usr/local -xzf go1.12.6.linux-amd64.tar.gz && \
  ln -s /usr/local/go/bin/go /bin/go && \
  ln -s /usr/local/go/bin/gofmt /bin/gofmt

#Set environment variables
ENV PATH=$PATH:/usr/local/go/bin:/usr/local/project/tastysearch/bin
ENV GO111MODULE=on CGO_ENABLED=0

EXPOSE 80

WORKDIR /usr/local/project/tastysearch

#Copy source directory
COPY ./ /usr/local/project/tastysearch

RUN go mod download

#Build
RUN \
    mkdir bin && \
    make build

ENTRYPOINT ["/usr/local/project/tastysearch/bin/server"]