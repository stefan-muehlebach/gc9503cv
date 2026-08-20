# gc9503cv

Seit geraumer Zeit liegen bei mir zwei TFT Flachmonitore herum, die beide eine etwas aussergewöhnliche Auflösung von 360x960 Pixel haben.
Angesteuert werden sie über ein 40-adriges FPC-Flachbandkabel, über welches im Wesentlichen eine 3Wire SPI-Verbindung (für die Konfiguration des Controllers) und eine 18 Bit RGB-Schnittstelle (für die Übermittlung der Bilddaten) realisiert ist.
Der Controller im Monitor ist ein GC9503CV von GalaxyCore, dessen hervorzuhebende Eigenschaft das fehlende GRAM (Graphical RAM) ist.
Es gibt also keinen graphischen Speicher auf Seiten Monitor, der das darzustellende Bild enthält.
Damit unterscheidet sich dieser Monitor von den meisten mir bekannten Monitoren (mit Controllern wie ILI9341, ILI9468, HX8357, etc.), die über ein eingebautes GRAM verfügen, ausschliesslich über eine SPI-Schnittstelle angesprochen werden können und für dies auch reichlich Software gibt.

Ich hatte vor, die Ansteuerung über einen Microcontroller (Arduino, Pico, etc) oder über einen RaspberryPi zu realisieren.
Als Programmiersprache war Go vorgesehen, da ich von Anfang an geplant hatte, dass dies ein Go-Projekt werden sollte.

## Pin-Belegung des Monitors

| Pin-Nr | Pin-Name | Gruppe | Beschreibung |
| ---: | --- | --- | --- |
| 1 | LEDA | LED | Anode der Hintergrundbeleuchtung |
| 2 | LEDK | LED | Kathode der Hintergrundbeleuchtung |
| 3 | LEDK | LED | Kathode der Hintergrundbeleuchtung |
| 4 | GND  | | Common Ground |
| 5 | VDD  | | Stromversorgung des Controllers (+3.2V) |
| 6 | RST  | | Reset Input Pin (low active) |
| 7 | NC   | | not connected |
| 8 | NC   | | not connected |
| 9 | SDA  | SPI | Datenleitung (Input/Output) des SPI-Interfaces |
| 10| SCK  | SPI | Clock-Leitung des SPI-Interfaces |
| 11| CS   | SPI | Chip-Select des SPI-Interfaces (low active) |
| 12| PCLK | RGB | Pixel-Clock Signal des RGB-Interfaces (rising edge) |
| 13| DE   | RGB | Data-Enable Signal des RGB-Interfaces (high active) |
| 14| VS   | RGB | Vertical-Sync Signal des RGB-Interfaces (low active) |
| 15| HS   | RGB | Horizontal-Sync Signal des RGB-Interfaces (low active) |
| 16-33 | DB0-DB17 | RGB | 18 Bit parallel RGB-Interface |
| 34| GND  | | Common Ground |
| 35-39 | TP | | Vorgesehen für ein kapazitives Touch-Panel |
| 40 | GND | | Common Ground |




