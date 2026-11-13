UPDATE categories SET icon = '' WHERE icon IS NULL;
UPDATE bookmarks SET logo_url = '' WHERE logo_url IS NULL;

ALTER TABLE categories ALTER COLUMN icon SET NOT NULL;
ALTER TABLE bookmarks ALTER COLUMN logo_url SET NOT NULL;
