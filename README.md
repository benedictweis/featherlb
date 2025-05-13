# featherlb - A lightweight Load Balancer

Zur Bewertung:

- Fertiges Projekt mit allen wichtigen Anforderungen erfüllt ist eine 1,5.
- Fertiges Projekt mit nicht erfüllten wichtigen Anforderungen ist eine 2,5.
- Besonders Gutes verbessert die Note um 0,25 bis 0,5.
- Wenn es nur mit extremem Aufwand funktioniert o.ä., verschlechtert das die Note um bis zu 0,5.

## Leute

- Benedict Weis

## Beschreibung

Ein simpler Load Balancer welche eingenende Pakete an mehrere Backends verteilt und ausgehende Pakete wieder zurück an den Client sendet.

## Estimation-Eintrag

- 13h: Implementierung des Load Balancers

## Planungsnotizen

Besondere Eigenschaft:

- Einfaches Deployment/einfache Ausführbarkeit (siehe unten)

Deliverables:

- Docker compose setup, das den Test laufen lässt
  - Backend mehrfach deployed
  - Client als Testprogramm, das Last erzeugt
  - Load Balancer
    - Liest config ein
    - Kennt verschiedene Algorithmen
    - Leitet TCP Anfragen entsprechend um

## Quellen

- <https://github.com/benedictweis/featherlb>
- <https://github.com/benedictweis/featherlb/tree/ebpf-playground>

## Metriken

- Erreichte Ziele (übererfüllt?)
  - [ ] …
  - [ ] …
- Besonders gut (bis 0,5 besser)? Wählen Sie *eins* aus:
  - [ ] Passende Tests,
  - [ ] Verständlichen Quellcode,
  - [ ] Dokumentation/Kommentare für komplizierte Stellen,
  - [x] Einfaches Deployment/einfache Ausführbarkeit,
  - [ ] Besonders Effizient,
  - [ ] Besonders Gut Skalierbar,
  - [ ] Angriffe in threat model analysiert und erkannte verhindert?
  - [ ] Besonders gute UX,
  - [ ] …
- [x] Quellcode verständlich? Komplizierte Stellen dokumentiert/kommentiert?
- [x] Ausführung funktioniert lokal / wurde als Video gezeigt?

## Deployment

Requirements:

- Bash <https://www.gnu.org/software/bash/>
- GNU Make <https://www.gnu.org/software/make/#download>
- Docker <https://docs.docker.com/get-started/get-docker/>

```bash
make e2e
```

Die Ausgabe in der Konsole dient nur zu Debugging-Zwecken. Die Resultate des e2e Tests werden in eine Datei unter `./test/e2e/runs/<datetime>.log` gespeichert. Dort werden die gemessenen Werte niedergeschrieben. Bei den Angaben zu backend1 und backend2 handelt es sich um die gemessene Anzahl an requests im access log des jeweiligen nginx servers. Das ganze Setup ist hochgradig konfigurierbar. In aktueller form wird mittels des [wrk](https://github.com/wg/wrk)-tools eine HTTP Last gegen den Load Balancer gesendet. Es werden dabei für 10s mit 12 Threads und 400 Connections so viele Anfragen gesendet wie der Load Balancer verarbeiten kann.

Abweichungen in den Tests sind zu erwarten und kommen stark auf die Platform an auf der der Test ausgeführt wird.

## eBPF

Ursprünglich war geplant, einen Load Balancer mit eBPF zu bauen. Einige Tests in diese Richtung sind unter <https://github.com/benedictweis/featherlb/tree/ebpf-playground> dokumentiert.

## Local Development

Starten des Load Balancers:

```bash
go run ./cmd/featherlb --debug --config featherlb.yaml
```

Starten zweier Web Server:

```bash
python3 -m http.server 8081
```

```bash
python3 -m http.server 8081
```

Anfragen senden:

```bash
curl localhost:8080
```
