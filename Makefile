.PHONY: test

test:
	set -e;\
	docker-compose -f ./integration-test/docker-compose.yml up -d;\
	docker-compose -f docker-compose.test.yml run integration_test bash -c "cd ./integration-test && go test --mod=vendor";\
	docker-compose -f docker-compose.test.yml down;