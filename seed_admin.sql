-- Seed akun administrator dan organisasi awal.
-- Password awal kedua akun: admin123
-- Segera ganti password setelah login pertama.

BEGIN;

INSERT INTO users (
    id, name, email, password, role, created_at, updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Super Admin',
    'admin@ukmhub.app',
    '$2y$10$q7l.wnJwwyyt6jRKj0IBq.0RC08pD8kmJJ0cs3p0JzM5fVi40SjPS',
    'SUPER_ADMIN',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO UPDATE SET
    name = EXCLUDED.name,
    password = EXCLUDED.password,
    role = EXCLUDED.role,
    updated_at = NOW();

INSERT INTO users (
    id, name, email, password, role, created_at, updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000002',
    'Admin CSA',
    'admin@csa.app',
    '$2y$10$q7l.wnJwwyyt6jRKj0IBq.0RC08pD8kmJJ0cs3p0JzM5fVi40SjPS',
    'ORG_ADMIN',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO UPDATE SET
    name = EXCLUDED.name,
    password = EXCLUDED.password,
    role = EXCLUDED.role,
    updated_at = NOW();

INSERT INTO organizations (
    id, name, slug, description, email, phone, status, created_at, updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000010',
    'Creative Student Association',
    'creative-student-association',
    'Organisasi mahasiswa kreatif dengan divisi Programming dan Multimedia.',
    'csa@ukmhub.app',
    '081234567890',
    'active',
    NOW(),
    NOW()
)
ON CONFLICT (slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    email = EXCLUDED.email,
    phone = EXCLUDED.phone,
    status = EXCLUDED.status,
    updated_at = NOW();

INSERT INTO organization_admins (
    id, user_id, organization_id, created_at, updated_at
) SELECT
    '00000000-0000-0000-0000-000000000020',
    u.id,
    o.id,
    NOW(),
    NOW()
FROM users u
CROSS JOIN organizations o
WHERE u.email = 'admin@csa.app'
  AND o.slug = 'creative-student-association'
ON CONFLICT (user_id, organization_id) DO NOTHING;

COMMIT;

-- Login:
-- Super admin: admin@ukmhub.app / admin123
-- Org admin:   admin@csa.app / admin123
