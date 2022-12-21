FROM ubuntu:20.04 as staging
# skip all interactive prompts for installs
ARG DEBIAN_FRONTEND=noninteractive
# install git for downloading go dependencies
RUN apt-get update && apt-get install -y git gcc openjdk-11-jre make

# major dependencies are vendored in the project and extracted here
COPY third_party /third_party
RUN tar xf /third_party/node/node-v14.5.0-linux-x64.tar.gz
RUN tar xf /third_party/go/go1.19.4.linux-amd64.tar.gz

# so we can find npm + go
ENV PATH=$PATH:/node-v14.5.0-linux-x64/bin:/go/bin

RUN npm install -g yarn

# create project dir
RUN mkdir vweb

WORKDIR vweb

# copy lock dependency and lock files
COPY go.mod ./
COPY go.sum ./
COPY client/package.json client/package.json
COPY client/yarn.lock client/yarn.lock

# install dependencies for frontend and backend
RUN go mod download
RUN cd client && yarn install

# grab the source
# backend server source
COPY pkg ./pkg
COPY cmd ./cmd
COPY internal ./internal
# client source
COPY client ./client
# app source

COPY makefile ./makefile
COPY api ./api
# build the frontend

RUN mkdir third_party && mv /third_party/openapi third_party/openapi
RUN make models
RUN make client/dist
RUN make bin/server

FROM ubuntu:20.04 as deploy
COPY --from=staging /vweb/bin/server /vweb/bin/server
COPY --from=staging /vweb/client/dist /vweb/client/dist

WORKDIR /vweb

CMD ["./bin/server"]
