# Shipment Api
This is a simple api server in Go for handling shipments with orders.

### Running the server

- make postgres
- make migrate
- go run cmd/main.go

You can run make adminer and consult the database from http://localhost:8080.

### Endpoints
##### Products

    POST /products
body required example: 
{
    "name": "smart tv",
    "price": 30.5,
    "volume": 20
}

    GET /products
response value example:
[{"id":"cdd5af0c-d7b4-4083-8697-269acdf5742c","name":"table","price":20,"volume":25},{"id":"6d530a61-a058-455a-b932-9e0c848e00c2","name":"smart tv","price":30.5,"volume":20}]

##### Orders
    POST /orders
body required example:
[
    {"product_id": "cdd5af0c-d7b4-4083-8697-269acdf5742c", "quantity": 3},
    {"product_id": "6d530a61-a058-455a-b932-9e0c848e00c2", "quantity": 2}
]

    GET /orders/{orderId}
response value example:
{"products":[{"order_items_id":"44120c42-344e-42c9-b8ed-66ead93432e8","quantity":3,"product_id":"cdd5af0c-d7b4-4083-8697-269acdf5742c","name":"table","price":20,"volume":25},{"order_items_id":"9702bf1c-0f55-4e3e-98cd-1462391cd057","quantity":2,"product_id":"6d530a61-a058-455a-b932-9e0c848e00c2","name":"smart tv","price":30.5,"volume":20}],"shipment_itemId_44120c42-344e-42c9-b8ed-66ead93432e8":[{"capacity":25,"start_date":"2022-05-11T11:15:45.683124Z","end_date":null,"quantity":3,"order_items_id":"44120c42-344e-42c9-b8ed-66ead93432e8","shipment_id":"4bdcba1b-6f95-41eb-bb90-95c5d1337e97"},{"capacity":25,"start_date":"2022-05-11T11:39:37.267575Z","end_date":null,"quantity":3,"order_items_id":"44120c42-344e-42c9-b8ed-66ead93432e8","shipment_id":"bb407118-b3ba-4fd2-912f-88973c239a79"}],"shipment_itemId_9702bf1c-0f55-4e3e-98cd-1462391cd057":[{"capacity":25,"start_date":"2022-05-11T11:15:45.683124Z","end_date":null,"quantity":2,"order_items_id":"9702bf1c-0f55-4e3e-98cd-1462391cd057","shipment_id":"4bdcba1b-6f95-41eb-bb90-95c5d1337e97"},{"capacity":25,"start_date":"2022-05-11T11:39:37.267575Z","end_date":null,"quantity":5,"order_items_id":"9702bf1c-0f55-4e3e-98cd-1462391cd057","shipment_id":"bb407118-b3ba-4fd2-912f-88973c239a79"}]}

##### Shipments
    POST /shipments?capacity=value
body required example:
[
    {"order_items_id": "44120c42-344e-42c9-b8ed-66ead93432e8", "quantity": 3},
    {"order_items_id": "9702bf1c-0f55-4e3e-98cd-1462391cd057", "quantity": 5}
]

    GET /shipments/{shipmentId}
response value example:
{"id":"4bdcba1b-6f95-41eb-bb90-95c5d1337e97","capacity":25,"start_date":"2022-05-11T11:15:45.683124Z","end_date":null}
    