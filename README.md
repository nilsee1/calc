# 🚕 Taxi Calculator

Ein schlanker Web-Taschenrechner in Go – Single Binary mit eingebettetem Frontend.

## Features

- Serverseitige Auswertung mathematischer Ausdrücke (Shunting-Yard-Algorithmus)
- Grundrechenarten (`+`, `-`, `*`, `/`) und Klammern
- Dark-Mode UI mit Verlaufsanzeige
- Vollständige Tastaturbedienung
- Single Binary Deployment via `go:embed`
- Docker-ready mit Multi-Stage Build

## Schnellstart

### Mit Go

```bash
go run .
# → http://localhost:8080
```

### Mit Docker

```bash
docker build -t taxi-calculator .
docker run -p 8080:8080 taxi-calculator
```

### Mit Docker Compose

```bash
docker compose up --build
# → http://localhost:8080
```

## Tests

```bash
go test ./... -v
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
├── docker-compose.yml   # Compose-Konfiguration
├── go.mod               # Go-Modul
└── README.md
```

## Lizenz

MIT

