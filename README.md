# featherlb - An eBPF based Load Balancer

Zur Bewertung:

- Fertiges Projekt mit allen wichtigen Anforderungen erfüllt ist eine 1,5.
- Fertiges Projekt mit nicht erfüllten wichtigen Anforderungen ist eine 2,5.
- Besonders Gutes verbessert die Note um 0,25 bis 0,5.
- Wenn es nur mit extremem Aufwand funktioniert o.ä., verschlechtert das die Note um bis zu 0,5.

## featherlb - An eBPF based Load Balancer

### Leute

- Benedict Weis

### Beschreibung

Ein simpler Load Balancer welche eingenende Pakete an mehrere Backends verteilt und ausgehende Pakete wieder zurück an den Client sendet.

eBPF wird verwdndet um Pakete auf Layer 3 abzufangen und den restlichen Network Stack zu umgeben. Zumal ist der Wechsel zwischen Kernel und User-space teuer und wird auch umgangen.

### Estimation-Eintrag

- 13h: Implementierung des Load Balancers

### Planungsnotizen

Wichtig: <https://ebpf.io/what-is-ebpf/>

Anforderungen:

- eBPF Anwendung welche auf Basis einer Map oder einem Array anfragen an Backends verteilt (C oder Rust)
- Userspace Anwendung welche eine Konfigurationsdatei einliest, welche die Backeds beinhalten und die Map oder einem Array befüllt (Go oder Rust)
- Optional
  - Session pinning
  - Metrics
  - Verschiedene Algorithmen (siehe <https://www.cloudflare.com/en-gb/learning/performance/types-of-load-balancing-algorithms/>)

Details:

- <https://cloudchirp.medium.com/go-c-rust-and-more-picking-the-right-ebpf-application-stack-7abd1c1ba9f4>
- <https://aya-rs.dev>
- <https://ebpf-go.dev>
- <https://github.com/eunomia-bpf/bpf-developer-tutorial/tree/main/src/42-xdp-loadbalancer>
- <https://medium.com/@oayirnil/lab-setting-up-a-rust-aya-ebpf-l4-load-balancer-dev-environment-184e643531f2>
- <https://github.com/lizrice/lb-from-scratch>
- <https://stackoverflow.com/questions/72120362/im-not-receiving-packets-using-xdp-tx>
- <https://www.youtube.com/watch?v=L3_AOFSNKK8>

Besondere Eigenschaft:

- Unglaublich Effizient da große Teile des Network Stacks umgangen werden

Deliverables:

- Docker compose setup, das den Test laufen lässt
  - Backend mehrfach deployed
  - Client als Testprogramm, das Last erzeugt
  - Load Balancer
    - eBPF Anwendung die Pakete umschreibt
    - Userspace Programm
      - Optional: Config einlesen
      - eBPF Anwendung "starten"

## Quellen

- Link zu Repository oder Zip in Moodle.

## Metriken

- Erreichte Ziele (übererfüllt?)
  - [ ] …
  - [ ] …
- Besonders gut (bis 0,5 besser)? Wählen Sie *eins* aus:
  - [ ] Passende Tests,
  - [ ] Verständlichen Quellcode,
  - [ ] Dokumentation/Kommentare für komplizierte Stellen,
  - [ ] Einfaches Deployment/einfache Ausführbarkeit,
  - [x] Besonders Effizient,
  - [ ] Besonders Gut Skalierbar?
  - [ ] Angriffe in threat model analysiert und erkannte verhindert?
  - [ ] Besonders gute UX,
  - [ ] …
- [ ] Quellcode verständlich? Komplizierte Stellen dokumentiert/kommentiert?
- [ ] Ausführung funktioniert lokal / wurde als Video gezeigt?

## Note: 

```bash
RUST_LOG=info cargo run --config 'target."cfg(all())".runner="sudo -E"' -- --iface eth0
```

If running in docker, create a new network and attach devcontainer and other container to it, then send traffic from the other container.

```bash
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'
```

```bash
docker inspect -f '{{range .NetworkSettings.Networks}}{{.MacAddress}}{{end}}'
```

```bash
docker network create -d macvlan featherlb_net
```
