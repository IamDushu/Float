CREATE TABLE "users" (
  "user_id" uuid PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "phone_number" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "nurses" (
  "nurse_id" uuid PRIMARY KEY,
  "user_id" uuid UNIQUE NOT NULL,
  "license_number" varchar UNIQUE NOT NULL,
  "specialization" varchar NOT NULL,
  "years_of_experience" integer NOT NULL,
  "zip_code" varchar(10) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "patients" (
  "patient_id" uuid PRIMARY KEY,
  "user_id" uuid UNIQUE NOT NULL,
  "date_of_birth" date NOT NULL,
  "emergency_contact_name" varchar NOT NULL,
  "emergency_contact_phone" varchar NOT NULL,
  "medical_history" text NOT NULL,
  "allergies" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "nurse_availability" (
  "availability_id" uuid PRIMARY KEY,
  "nurse_id" uuid NOT NULL,
  "day_of_week" varchar NOT NULL,
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
  "status" varchar NOT NULL,
  "notes" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "email_verification" (
  "verification_id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "token" text UNIQUE NOT NULL,
  "hashed_otp" varchar NOT NULL,
  "purpose" varchar NOT NULL,
  "attempts" integer NOT NULL DEFAULT 0,
  "expires_at" timestamptz NOT NULL,
  "valid" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "sessions" (
  "session_id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);


CREATE INDEX ON "users" ("email");

CREATE INDEX ON "visits" ("nurse_id");

CREATE INDEX ON "visits" ("patient_id");

CREATE UNIQUE INDEX email_purpose_valid_key ON email_verification (email, purpose) WHERE valid = true;

ALTER TABLE "nurses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "patients" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "nurse_availability" ADD FOREIGN KEY ("nurse_id") REFERENCES "nurses" ("nurse_id");

ALTER TABLE "visits" ADD FOREIGN KEY ("nurse_id") REFERENCES "nurses" ("nurse_id");

ALTER TABLE "visits" ADD FOREIGN KEY ("patient_id") REFERENCES "patients" ("patient_id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("email") REFERENCES "users" ("email");