CREATE TABLE IF NOT EXISTS banned_users (
    chat_id   BIGINT PRIMARY KEY,
    banned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    reason    TEXT
);

CREATE OR REPLACE FUNCTION delete_user_on_ban()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM users WHERE chat_id = NEW.chat_id;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_delete_user_on_ban ON banned_users;
CREATE TRIGGER trg_delete_user_on_ban
AFTER INSERT ON banned_users
FOR EACH ROW
EXECUTE PROCEDURE delete_user_on_ban();

CREATE OR REPLACE FUNCTION prevent_banned_signup()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF EXISTS (SELECT 1 FROM banned_users WHERE chat_id = NEW.chat_id) THEN
        RAISE EXCEPTION 'user is banned'
            USING ERRCODE = '42501';
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_prevent_banned_signup ON users;
CREATE TRIGGER trg_prevent_banned_signup
BEFORE INSERT OR UPDATE OF chat_id ON users
FOR EACH ROW
EXECUTE PROCEDURE prevent_banned_signup();
