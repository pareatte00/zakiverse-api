CREATE OR REPLACE FUNCTION sanitize_card_tag_name()
RETURNS TRIGGER AS $$
BEGIN
    NEW.name := LOWER(TRIM(REGEXP_REPLACE(NEW.name, '\s{2,}', ' ', 'g')));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_card_tag_sanitize_name ON card_tag;

CREATE TRIGGER trg_card_tag_sanitize_name
    BEFORE INSERT OR UPDATE OF name ON card_tag
    FOR EACH ROW
    EXECUTE FUNCTION sanitize_card_tag_name();
