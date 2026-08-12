.PHONY: run web tunnel

run:
	@go run ./cmd/app

web:
	@echo "Mini App: http://localhost:5173"
	@python3 -m http.server 5173 --directory web >/dev/null 2>&1

# Прод/стабильный URL после push в main:
# https://hronovv.github.io/FlatStalker/
