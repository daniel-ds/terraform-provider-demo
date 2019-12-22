.PHONY: build test testacc testacc-default start-es stop-es
TEST?=$$(go list ./... |grep -v 'vendor')

default: build

start-es:
	docker run -p 9200:9200 -p 9300:9300 --rm --name elastic-demo -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.5.2

start-es-newport:
	docker run -p 9201:9200 -p 9300:9300 --rm --name elastic-demo -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.5.2

stop-es:
	docker stop elastic-demo

build:
	go build -o terraform-provider-demo

test:
	go test ./... -v -count=1

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -count=1 -timeout 5m

testacc-default:
	TEST_DEMO_URL="http://localhost:9200" TF_ACC=1 go test $(TEST) -v $(TESTARGS) -count=1 -timeout 5m
