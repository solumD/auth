-- +goose Up
-- +goose StatementBegin
create or replace function update_modified_column()
returns trigger as $$
begin
NEW.updated_at = now();
return NEW;
end;
$$ language 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
create trigger update_modified_time before update on users for each row execute procedure update_modified_column();
-- +goose StatementEnd