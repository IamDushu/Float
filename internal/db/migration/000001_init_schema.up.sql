CREATE TABLE "users" (
  "user_id" uuid PRIMARY KEY,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "first_name" VARCHAR(100) NOT NULL,
  "last_name" VARCHAR(100) NOT NULL,
  "phone_number" VARCHAR(20) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "nurses" (
  "nurse_id" uuid PRIMARY KEY,
  "user_id" uuid UNIQUE NOT NULL,
  "license_number" VARCHAR(50) UNIQUE NOT NULL,
  "specialization" VARCHAR(100) NOT NULL,
  "years_of_experience" INTEGER NOT NULL,
  "zip_code" VARCHAR(10) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "patients" (
  "patient_id" uuid PRIMARY KEY,
  "user_id" uuid UNIQUE NOT NULL,
  "date_of_birth" DATE NOT NULL,
  "emergency_contact_name" VARCHAR(200) NOT NULL,
  "emergency_contact_phone" VARCHAR(20) NOT NULL,
  "medical_history" TEXT NOT NULL,
  "allergies" TEXT NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "nurse_availability" (
  "availability_id" uuid PRIMARY KEY,
  "nurse_id" uuid NOT NULL,
  "day_of_week" varchar(20) NOT NULL,
  "start_time" time NOT NULL,
  "end_time" time NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "visits" (
  "visit_id" uuid PRIMARY KEY,
  "nurse_id" uuid,
  "patient_id" uuid NOT NULL,
  "scheduled_at" timestamptz NOT NULL,
  "completed_at" timestamptz NOT NULL,
  "status" VARCHAR(20) NOT NULL,
  "notes" TEXT NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE INDEX "idx_users_email" ON "users" ("email");

CREATE INDEX "idx_visits_nurse_id" ON "visits" ("nurse_id");

CREATE INDEX "idx_prescriptions_patient_id" ON "visits" ("patient_id");

ALTER TABLE "nurses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "patients" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "nurse_availability" ADD FOREIGN KEY ("nurse_id") REFERENCES "nurses" ("nurse_id");

ALTER TABLE "visits" ADD FOREIGN KEY ("nurse_id") REFERENCES "nurses" ("nurse_id");

ALTER TABLE "visits" ADD FOREIGN KEY ("patient_id") REFERENCES "patients" ("patient_id");
