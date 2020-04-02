## /company

POST /company {company_name, opening_hours ....} (add)
PATCH /company (update)
GET /company

POST /company/login

## With company token

PATCH /company/slots
POST /company/slots
GET /company/slots

GET /company/slots/id -- ger slot med visst id

### To scan

GET /company/{alfanumeric code} -> {status: true/false}

## /public

POST /reserveBook

GET /stores {location, radius=1km}

GET /stores/slots {day}

GET /search?q="hello" -> Elasticsearch

POST /book {phonenum}
POST /book/confirm {code}
POST /unbook -- rate limit

GET /book/{alfanumerisk kod} -> qr + alfanumerisk kod

-- NOT MVP /company/verify
