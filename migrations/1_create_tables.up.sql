CREATE TABLE "order_items" (
  "id" UUID PRIMARY KEY,
  "order_id" UUID,
  "product_id" UUID,
  "quantity" int
);

CREATE TABLE "orders" (
  "id" UUID PRIMARY KEY,
  "delivered" boolean,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "products" (
  "id" UUID PRIMARY KEY,
  "name" varchar,
  "price" float,
  "volume" float
);

CREATE TABLE "shipments" (
  "id" UUID PRIMARY KEY,
  "capacity" float,
  "start_date" timestamp,
  "end_date" timestamp
);

CREATE TABLE "shipment_products" (
  "shipment_id" UUID,
  "order_items_id" UUID,
  "quantity" int
);

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "shipment_products" ADD FOREIGN KEY ("shipment_id") REFERENCES "shipments" ("id");

ALTER TABLE "shipment_products" ADD FOREIGN KEY ("order_items_id") REFERENCES "order_items" ("id");