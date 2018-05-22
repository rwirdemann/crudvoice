# Restvoice

Eine API zur Rechnungserstellung für Freiberufler und kleinere Firmen.

## Ressourcen

### POST /invoices
```
curl -d '{"year": 2018, "month":12, "customerId"}' \
    -H "Content-Type: application/json" -X POST \
    http://localhost:8080/invoices
```

### POST /invoices/bookings

### PUT /invoices/charge/1234