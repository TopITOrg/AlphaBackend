-- name: GetUserByEmail :one
SELECT
    users.id,
    users.full_name,
    users.social_network_link,
    users.phone_number,
    users.email,
    users.birth_date,
    users.role,
    users.password,
    users.created_at,
    users.updated_at,
    (groups.prefix || '-' || extract(YEAR FROM age(now(), (groups.enrollment_year::text || '-09-01 00:00:00')::timestamptz)) + 1 || lpad(groups.group_number::text, 2, '0') || group_types.name || '-' || substring(enrollment_year::text FROM 3 FOR 2))::text as group_name
FROM users
LEFT JOIN groups on users.group_id = groups.id AND groups.is_deleted = false
LEFT JOIN group_types on groups.group_type_id = group_types.id and group_types.is_deleted = false
WHERE users.email = @email
LIMIT 1;

-- name: GetUserById :one
SELECT
    users.id,
    users.full_name,
    users.social_network_link,
    users.phone_number,
    users.email,
    users.birth_date,
    users.role,
    users.password,
    users.created_at,
    users.updated_at,
    (groups.prefix || '-' || extract(YEAR FROM age(now(), (groups.enrollment_year::text || '-09-01 00:00:00')::timestamptz)) + 1 || lpad(groups.group_number::text, 2, '0') || group_types.name || '-' || substring(enrollment_year::text FROM 3 FOR 2))::text as group_name
FROM users
         LEFT JOIN groups on users.group_id = groups.id AND groups.is_deleted = false
         LEFT JOIN group_types on groups.group_type_id = group_types.id and group_types.is_deleted = false
WHERE users.id = @id
    LIMIT 1;

-- name: CreateUser :one
WITH user_info AS (
    INSERT INTO users
    (full_name, social_network_link, phone_number, email, birth_date, role, password, group_id)
    VALUES (@full_name, @social_network_link, @phone_number, @email, @birth_date, @role, @password, @group_id)
    RETURNING *
)
SELECT 
    user_info.*,
    CASE 
        WHEN groups.id IS NOT NULL THEN
            (groups.prefix || '-' || 
             extract(YEAR FROM age(now(), (groups.enrollment_year::text || '-09-01 00:00:00')::timestamptz)) + 1 || 
             lpad(groups.group_number::text, 2, '0') || 
             group_types.name || '-' || 
             substring(groups.enrollment_year::text FROM 3 FOR 2)
            )::text
        ELSE ''
    END as group_name
FROM user_info
LEFT JOIN groups ON groups.id = user_info.group_id AND groups.is_deleted = false
LEFT JOIN group_types ON groups.group_type_id = group_types.id AND group_types.is_deleted = false;

-- name: DeleteUser :exec
UPDATE users
SET
    is_deleted = true
WHERE email = @email and is_deleted = false;