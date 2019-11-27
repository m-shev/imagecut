.PHONY: build run qa prod test

build:
	bash -c "mkdir -p app-data/log && mkdir -p app-data/cache && mkdir -p app-data/images"
	go build -ldflags="-w -s" -mod=vendor -o build/imagecut cmd/main.go

run: build
	build/imagecut

qa:
	docker run -p=3000:3000 --rm -e IMAGECUT_ENV=QA -v imagecut-data-qa:/opt/imagecut/app-data mshev/imagecut

prod:
	docker run -p=3000:3000 --rm -e IMAGECUT_ENV=PROD -v imagecut-data:/opt/imagecut/app-data -d mshev/imagecut

test:
	docker-compose -f ./integration-test/docker-compose.yml up -d;\
	docker-compose -f ./integration-test/docker-compose.yml run integration_test bash -c "cd ./integration-test && go test --mod=vendor";\
	docker-compose -f ./integration-test/docker-compose.yml down -v;