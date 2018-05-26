# Restvoice

Eine REST API zur Rechnungserstellung für Freiberufler und kleinere Firmen.

## Entry Point

http://restvoice.org/invoice

## CRUD API

POST /invoices

POST /invoices/1234/bookings

PATCH /invoices/1234 {"status": "payment expected"}

PATCH /invoices/1234 {"status": "payed"}

PATCH /invoices/1234 {"status": "archived"}

## Hypermedia API

POST    /invoice               => created

POST    /invoice/1234/booking  => created

PUT     /charge/1234           => payment expected

PUT     /payment/1234          => payed

PUT     /archive/1234          => archived







