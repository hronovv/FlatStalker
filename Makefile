.PHONY: run web tunnel logs

SERVER ?= $(shell awk -F= '/^SERVER=/{print $$2; exit}' .env)

run:
	@go run ./cmd/app

logs:
	@test -n "$(SERVER)" || (echo "SERVER=root@IP нет в .env" && exit 1)
	ssh -t $(SERVER) 'journalctl -u flatstalker -f -o cat'

web:
	@echo "Mini App: http://localhost:5173"
	@python3 -m http.server 5173 --directory web >/dev/null 2>&1
