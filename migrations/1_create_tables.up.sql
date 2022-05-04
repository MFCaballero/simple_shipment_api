CREATE TABLE "order_items" (
  "order_id" int,
  "product_id" int,
  "quantity" int DEFAULT 1
);

CREATE TABLE "orders" (
  "id" int PRIMARY KEY,
  "delivered" boolean,
  "created_at" datetime DEFAULT (now())
);

CREATE TABLE "products" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "price" float,
  "volume" float
);

CREATE TABLE "shipments" (
  "id" int PRIMARY KEY,
  "capacity" float,
  "start_date" datetime,
  "end_date" datetime
);

CREATE TABLE "shipment_products" (
  "shipment_id" int,
  "product_id" int
);

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "shipment_products" ADD FOREIGN KEY ("shipment_id") REFERENCES "shipments" ("id");

ALTER TABLE "shipment_products" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");