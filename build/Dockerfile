FROM ubuntu:20.04

# install git for downloading go dependencies
RUN apt-get update && apt-get install -y git gcc openjdk-8-jre make

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
COPY rewinder/app/src/main/java/com ./rewinder/app/src/main/java/com

COPY makefile ./makefile
COPY swagger.yaml ./swagger.yaml
# build the frontend

RUN mkdir third_party && mv /third_party/swagger-codegen-cli-3.0.20.jar third_party/swagger-codegen-cli-3.0.20.jar
RUN make models
RUN cd pkg/www/retext && npm run build

# build the backend
RUN mkdir -p bin && go build -o bin/server cmd/server/main.go

FROM ubuntu:20.04 as deploy

COPY --from=staging /retext/bin/server /retext/bin/server
COPY --from=staging /retext/pkg/www/retext/dist /retext/pkg/www/retext/dist

WORKDIR /retext

CMD ["./bin/server"]
