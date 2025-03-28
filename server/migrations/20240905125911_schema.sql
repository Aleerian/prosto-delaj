-- +goose Up
-- +goose StatementBegin
CREATE TABLE "User"
(
    "email"        character varying NOT NULL,
    "name"         character varying NOT NULL,
    CONSTRAINT "Account_pk" PRIMARY KEY ("email")
) WITHOUT OIDS;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "User";
-- +goose StatementEnd
