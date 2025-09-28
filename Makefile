build:
	bun run build
	templ generate
	go build -tags goexperiment.jsonv2 -ldflags="-s -w" -o ./cmd/quiz-me/quiz-me.so ./cmd/quiz-me

run:
	bun run build
	templ generate
	go run -tags goexperiment.jsonv2 ./cmd/quiz-me

run-dev:
	go run -tags 'goexperiment.jsonv2 dev' ./cmd/quiz-me

save:
	docker build -t quiz-me .
	docker save -o quiz-me.tar quiz-me
