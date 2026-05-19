# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[8.1].define(version: 2026_05_18_171008) do
  # These are extensions that must be enabled in order to support this database
  enable_extension "pg_catalog.plpgsql"

  create_table "api_events", id: :serial, force: :cascade do |t|
    t.text "client_id", null: false
    t.timestamptz "created_at", default: -> { "now()" }, null: false
    t.text "data_source", null: false
    t.text "endpoint", null: false
    t.integer "status_code", null: false
    t.boolean "success", null: false
    t.timestamptz "timestamp", null: false
    t.index ["client_id"], name: "idx_api_events_client_id"
    t.index ["client_id"], name: "index_api_events_on_client_id"
    t.index ["timestamp"], name: "idx_api_events_timestamp"
    t.index ["timestamp"], name: "index_api_events_on_timestamp"
  end

  create_table "batch_student_results", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.uuid "batch_student_id", null: false
    t.timestamptz "created_at", default: -> { "CURRENT_TIMESTAMP" }
    t.boolean "found_enrollment", default: false, null: false
    t.jsonb "results"
    t.index ["batch_student_id"], name: "idx_batch_student_results_batch_student_id"
  end

  create_table "batch_students", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.uuid "batch_db_id", null: false
    t.timestamptz "created_at", default: -> { "CURRENT_TIMESTAMP" }
    t.text "date_of_birth", null: false
    t.text "first_name", null: false
    t.text "last_name", null: false
    t.text "record_id", null: false
    t.text "ssn", null: false
    t.text "status", null: false
    t.index ["batch_db_id"], name: "idx_batch_students_batch_db_id"
  end

  create_table "education_batch_student_results", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.datetime "created_at", null: false
    t.uuid "education_batch_student_id", null: false
    t.boolean "found_enrollment"
    t.jsonb "results"
    t.datetime "updated_at", null: false
    t.index ["education_batch_student_id"], name: "idx_on_education_batch_student_id_662efa214a", unique: true
  end

  create_table "education_batch_students", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.datetime "created_at", null: false
    t.string "date_of_birth", null: false
    t.uuid "education_enrollment_batch_id", null: false
    t.string "first_name", null: false
    t.string "last_name", null: false
    t.string "record_id", null: false
    t.string "ssn", null: false
    t.string "status", null: false
    t.datetime "updated_at", null: false
    t.index ["education_enrollment_batch_id"], name: "idx_on_education_enrollment_batch_id_787e79cd79"
  end

  create_table "education_enrollment_batches", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.string "batch_id", null: false
    t.string "callback_url", null: false
    t.datetime "created_at", null: false
    t.string "status", null: false
    t.string "submitted_by", null: false
    t.datetime "updated_at", null: false
    t.index ["batch_id"], name: "index_education_enrollment_batches_on_batch_id", unique: true
  end

  create_table "enrollment_batches", id: :uuid, default: -> { "gen_random_uuid()" }, force: :cascade do |t|
    t.text "batch_id", null: false
    t.text "callback_url", null: false
    t.timestamptz "created_at", default: -> { "CURRENT_TIMESTAMP" }
    t.text "status", null: false
    t.text "submitted_by", null: false

    t.unique_constraint ["batch_id"], name: "enrollment_batches_batch_id_key"
  end

  add_foreign_key "batch_student_results", "batch_students", name: "batch_student_results_batch_student_id_fkey", on_delete: :cascade
  add_foreign_key "batch_students", "enrollment_batches", column: "batch_db_id", name: "batch_students_batch_db_id_fkey", on_delete: :cascade
  add_foreign_key "education_batch_student_results", "education_batch_students"
  add_foreign_key "education_batch_students", "education_enrollment_batches"
end
