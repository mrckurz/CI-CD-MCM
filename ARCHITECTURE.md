# Übung 2

## Aufgabe 1

Die API ist in mehrere Schichten aufgeteilt. Jede Schicht hat dabei eine bestimmte Aufgabe.

### Handler

Die Handler kümmern sich um eingehende Requests. Dabei werden die Requests auf ihre Gültigkeit geprüft. Wenn ein Request gültig ist, werden über die Stores Daten abgefragt bzw. gespeichert.

### Store

Stores bilden die Abstraktionsebene zur Datenpersistierung. Sie sind dafür verantwortlich, Daten zu lesen bzw. zu schreiben.

### Model

Modelle werden in den verschiedenen Schichten verwendet. Sie repräsentieren Ressourcen.

Das folgende Diagramm zeigt den Ablauf eines Requests bis zur Datenbank.

![Ablauf eines Requests](./assets/task_02.png)

Es gibt zwei Stores: eine In-Memory-Version und eine Datenbank-Version. Die In-Memory-Version eignet sich gut zum Testen, da dabei die Abhängigkeit zu externen Systemen reduziert wird. Die Datenbank-Version ist für den produktiven Einsatz gedacht, da sie die Daten dauerhaft speichert.
