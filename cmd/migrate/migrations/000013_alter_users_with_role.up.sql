-- set default role for data already in the table
ALTER TABLE IF EXISTS users
ADD COLUMN IF NOT EXISTS role_id INT REFERENCES roles (id) DEFAULT 1;


-- set default user role
UPDATE users
SET
    role_id = (
        SELECT
            id
        FROM
            roles
        WHERE
            name = 'user'
    );

-- can drop default after filling the table
ALTER TABLE IF EXISTS users
ALTER COLUMN role_id DROP DEFAULT;

-- set not null
ALTER TABLE IF EXISTS users
ALTER COLUMN role_id SET NOT NULL;


-- **postgres會先在舊資料加欄位然後檢查規則，最後才塞default 所以不能先設not null

