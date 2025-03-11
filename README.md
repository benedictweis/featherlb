# featherlb - An eBPF based Load Balancer

Zur Bewertung:

- Fertiges Projekt mit allen wichtigen Anforderungen erfüllt ist eine 1,5.
- Fertiges Projekt mit nicht erfüllten wichtigen Anforderungen ist eine 2,5.
- Besonders Gutes verbessert die Note um 0,25 bis 0,5.
- Wenn es nur mit extremem Aufwand funktioniert o.ä., verschlechtert das die Note um bis zu 0,5.


## Leute

- Benedict Weis

## Beschreibung

## Estimation-Eintrag

- 10h Implement Load Balancer using eBPF with aya-rs

## Planungsnotizen

Anforderungen:

- …
- …

Details:

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
