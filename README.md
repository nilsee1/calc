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

Image von der Codeberg Container Registry pullen und starten:

```bash
docker compose pull
docker compose up -d
# → http://localhost:8080
```

### Mit Docker (manuell)

```bash
docker pull codeberg.org/nilsee/calc:latest
docker run -p 8080:8080 codeberg.org/nilsee/calc:latest
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

Bei Codeberg einloggen, Image bauen und in die Registry pushen:

```bash
# 1. Bei der Codeberg Container Registry anmelden
docker login codeberg.org

# 2. Image bauen
docker build -t codeberg.org/nilsee/calc:latest .

# 3. Image pushen
docker push codeberg.org/nilsee/calc:latest
```

> **Tipp:** Als Passwort beim `docker login` ein
> [Access Token](https://codeberg.org/user/settings/applications)
> mit `package:write`-Berechtigung verwenden.

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
