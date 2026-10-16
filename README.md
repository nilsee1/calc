# Calculator

Ein schlanker Web-Taschenrechner in Go – Single Binary mit eingebettetem Frontend.

## Features

- Serverseitige Auswertung mathematischer Ausdrücke (Shunting-Yard-Algorithmus)
- Grundrechenarten (`+`, `-`, `*`, `/`) und Klammern
- Dark-Mode UI mit Verlaufsanzeige
- Vollständige Tastaturbedienung
- Single Binary Deployment via `go:embed`
- Docker-ready mit Multi-Stage Build

## Schnellstart

### Mit Docker Compose (empfohlen)

Image von der GitHub Container Registry (GHCR) pullen und starten:

```bash
docker compose pull
docker compose up -d
# → http://localhost:8080
```

### Mit Docker (manuell)

```bash
docker pull ghcr.io/nilsee1/calc:latest
docker run -p 8080:8080 ghcr.io/nilsee1/calc:latest
```

### Mit Go (Entwicklung)

```bash
go run .
# → http://localhost:8080
```

## Tests

```bash
go test ./... -v
```

## Image bauen & pushen (Maintainer)

Images werden automatisch via **GitHub Actions** bei jedem Push auf den `main`-Branch gebaut und in die GitHub Container Registry (GHCR) gepusht.

Manuell bauen und pushen:

```bash
# 1. Bei GHCR anmelden (mit GitHub Personal Access Token mit write:packages Berechtigung)
echo $GITHUB_TOKEN | docker login ghcr.io -u nilsee1 --password-stdin

# 2. Image bauen
docker build -t ghcr.io/nilsee1/calc:latest .

# 3. Image pushen
docker push ghcr.io/nilsee1/calc:latest
```

## API

### `POST /api/calculate`

**Request:**
```json
{"expression": "12.5 * (3 + 4)"}
```

**Response (Erfolg):**
```json
{"result": 87.5}
```

**Response (Fehler):**
```json
{"error": "Division durch Null"}
```

### `GET /health`

Gibt `200 OK` zurück – für Docker Healthchecks.

## Tastaturkürzel

| Taste       | Funktion      |
|-------------|---------------|
| `0-9`, `.`  | Zifferneingabe|
| `+ - * /`   | Operatoren    |
| `( )`       | Klammern      |
| `Enter` / `=` | Berechnen   |
| `Backspace` | Letzte Eingabe löschen |
| `Escape`    | Alles löschen |

## Projektstruktur

```
├── main.go              # Server & Routing
├── calculator.go        # Rechenlogik (Shunting-Yard)
├── calculator_test.go   # Unit-Tests
├── static/
│   ├── index.html       # UI
│   ├── style.css        # Styling
│   └── app.js           # Frontend-Logik
├── Dockerfile           # Multi-Stage Build
├── docker-compose.yml   # Compose-Konfiguration (pull von Registry)
├── go.mod               # Go-Modul
└── README.md
```

## Lizenz

MIT
