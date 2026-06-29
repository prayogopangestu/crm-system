-- Jalankan setelah migration selesai.
-- Login demo:
--   email    : admin@example.com
--   password : Admin123!

CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $seed$
DECLARE
  v_org_id uuid;
  v_admin_id uuid;
  v_sales_one_id uuid;
  v_sales_two_id uuid;
BEGIN
  SELECT id, organization_id
    INTO v_admin_id, v_org_id
  FROM users
  WHERE lower(email) = lower('admin@example.com')
    AND revoked_at IS NULL
  LIMIT 1;

  IF v_admin_id IS NULL THEN
    v_org_id := '00000000-0000-4000-8000-000000000100';
    v_admin_id := '00000000-0000-4000-8000-000000000101';

    INSERT INTO organizations (id, name)
    VALUES (v_org_id, 'CRM Demo')
    ON CONFLICT (id) DO UPDATE
      SET name = EXCLUDED.name,
          updated_at = now();

    INSERT INTO users (
      id, organization_id, first_name, last_name, email, password_hash, role, status
    )
    VALUES (
      v_admin_id, v_org_id, 'Demo', 'Admin', 'admin@example.com',
      crypt('Admin123!', gen_salt('bf', 12)), 'Admin', 'Aktif'
    );
  ELSE
    UPDATE users
    SET first_name = 'Demo',
        last_name = 'Admin',
        password_hash = crypt('Admin123!', gen_salt('bf', 12)),
        role = 'Admin',
        status = 'Aktif',
        updated_at = now()
    WHERE id = v_admin_id;
  END IF;

  SELECT id INTO v_sales_one_id
  FROM users
  WHERE lower(email) = lower('budi.sales@example.com')
  LIMIT 1;

  IF v_sales_one_id IS NULL THEN
    v_sales_one_id := '00000000-0000-4000-8000-000000000102';
    INSERT INTO users (
      id, organization_id, first_name, last_name, email, password_hash, role, status
    )
    VALUES (
      v_sales_one_id, v_org_id, 'Budi', 'Santoso', 'budi.sales@example.com',
      crypt('Sales123!', gen_salt('bf', 12)), 'Staf Sales', 'Aktif'
    );
  ELSE
    UPDATE users
    SET organization_id = v_org_id,
        first_name = 'Budi',
        last_name = 'Santoso',
        role = 'Staf Sales',
        status = 'Aktif',
        revoked_at = NULL,
        updated_at = now()
    WHERE id = v_sales_one_id;
  END IF;

  SELECT id INTO v_sales_two_id
  FROM users
  WHERE lower(email) = lower('siti.sales@example.com')
  LIMIT 1;

  IF v_sales_two_id IS NULL THEN
    v_sales_two_id := '00000000-0000-4000-8000-000000000103';
    INSERT INTO users (
      id, organization_id, first_name, last_name, email, password_hash, role, status
    )
    VALUES (
      v_sales_two_id, v_org_id, 'Siti', 'Rahma', 'siti.sales@example.com',
      crypt('Sales123!', gen_salt('bf', 12)), 'Staf Sales', 'Aktif'
    );
  ELSE
    UPDATE users
    SET organization_id = v_org_id,
        first_name = 'Siti',
        last_name = 'Rahma',
        role = 'Staf Sales',
        status = 'Aktif',
        revoked_at = NULL,
        updated_at = now()
    WHERE id = v_sales_two_id;
  END IF;

  INSERT INTO pipeline_stages (id, organization_id, key, name, color, position, is_system)
  VALUES
    ('00000000-0000-4000-8000-000000000201', v_org_id, 'lead', 'Lead Masuk', 'bg-primary-container', 1, true),
    ('00000000-0000-4000-8000-000000000202', v_org_id, 'contacted', 'Dihubungi', 'bg-secondary-container', 2, true),
    ('00000000-0000-4000-8000-000000000203', v_org_id, 'meeting', 'Meeting', 'bg-tertiary-container', 3, true),
    ('00000000-0000-4000-8000-000000000204', v_org_id, 'negotiation', 'Negosiasi', 'bg-primary-fixed', 4, true),
    ('00000000-0000-4000-8000-000000000205', v_org_id, 'won', 'Deal Won', 'bg-surface-tint', 5, true),
    ('00000000-0000-4000-8000-000000000206', v_org_id, 'lost', 'Deal Lost', 'bg-error-container', 6, true)
  ON CONFLICT (organization_id, key) DO UPDATE
    SET name = EXCLUDED.name,
        color = EXCLUDED.color,
        position = EXCLUDED.position,
        is_system = EXCLUDED.is_system,
        updated_at = now();

  INSERT INTO contacts (
    id, organization_id, owner_id, name, email, company, role, status, last_contacted_at
  )
  VALUES
    ('00000000-0000-4000-8000-000000000301', v_org_id, v_admin_id, 'Budi Wijaya', 'budi.w@telkomsel.co.id', 'PT Telkomsel', 'VP Sales', 'Negosiasi', now() - interval '2 hours'),
    ('00000000-0000-4000-8000-000000000302', v_org_id, v_sales_one_id, 'Anita Sari', 'anita.sari@gojek.com', 'Gojek Indonesia', 'Procurement Manager', 'Menang', now() - interval '1 day'),
    ('00000000-0000-4000-8000-000000000303', v_org_id, v_sales_two_id, 'Dian Pratama', 'dian.pratama@bca.co.id', 'Bank BCA', 'Head of IT Infrastructure', 'Prospek Awal', now() - interval '3 days'),
    ('00000000-0000-4000-8000-000000000304', v_org_id, v_admin_id, 'Hendra Kusuma', 'hendra@astra.co.id', 'PT Astra International', 'Supply Chain Director', 'Proposal', now() - interval '5 days'),
    ('00000000-0000-4000-8000-000000000305', v_org_id, v_sales_one_id, 'Rina Novita', 'rina@unilever.co.id', 'Unilever Indonesia', 'Marketing Manager', 'Kualifikasi', now() - interval '7 days')
  ON CONFLICT (id) DO UPDATE
    SET organization_id = EXCLUDED.organization_id,
        owner_id = EXCLUDED.owner_id,
        name = EXCLUDED.name,
        email = EXCLUDED.email,
        company = EXCLUDED.company,
        role = EXCLUDED.role,
        status = EXCLUDED.status,
        last_contacted_at = EXCLUDED.last_contacted_at,
        deleted_at = NULL,
        updated_at = now();

  INSERT INTO deals (
    id, organization_id, assignee_id, title, company, value, priority, stage_key, lost_reason, created_at, updated_at
  )
  VALUES
    ('00000000-0000-4000-8000-000000000401', v_org_id, v_admin_id, 'PT Telkomsel', 'Implementasi CRM Enterprise', 45000000, 'High', 'lead', '', now() - interval '9 days', now() - interval '9 days'),
    ('00000000-0000-4000-8000-000000000402', v_org_id, v_sales_one_id, 'Bank BCA', 'Upgrade sistem pipeline sales', 120000000, 'Medium', 'meeting', '', now() - interval '18 days', now() - interval '8 days'),
    ('00000000-0000-4000-8000-000000000403', v_org_id, v_sales_two_id, 'Gojek Indonesia', 'Lisensi CRM tahunan', 250000000, 'High', 'won', '', date_trunc('month', now()) + interval '2 days', now() - interval '2 days'),
    ('00000000-0000-4000-8000-000000000404', v_org_id, v_sales_one_id, 'PT Astra International', 'Integrasi notifikasi sales', 80000000, 'Medium', 'negotiation', '', now() - interval '6 days', now() - interval '1 day'),
    ('00000000-0000-4000-8000-000000000405', v_org_id, v_admin_id, 'Unilever Indonesia', 'Pilot project CRM mobile', 35000000, 'Low', 'lost', 'Harga terlalu tinggi', now() - interval '20 days', now() - interval '4 days')
  ON CONFLICT (id) DO UPDATE
    SET organization_id = EXCLUDED.organization_id,
        assignee_id = EXCLUDED.assignee_id,
        title = EXCLUDED.title,
        company = EXCLUDED.company,
        value = EXCLUDED.value,
        priority = EXCLUDED.priority,
        stage_key = EXCLUDED.stage_key,
        lost_reason = EXCLUDED.lost_reason,
        deleted_at = NULL,
        updated_at = EXCLUDED.updated_at;

  INSERT INTO tasks (
    id, organization_id, assignee_id, title, company, due_date, due_time, type, priority, notes, completed
  )
  VALUES
    ('00000000-0000-4000-8000-000000000501', v_org_id, v_admin_id, 'Telepon Budi Wijaya', 'PT Telkomsel', current_date, '14:00', 'Call', 'Tinggi', 'Konfirmasi jadwal demo dan kebutuhan modul sales.', false),
    ('00000000-0000-4000-8000-000000000502', v_org_id, v_sales_one_id, 'Kirim revisi proposal', 'Bank BCA', current_date, '16:30', 'Proposal', 'Sedang', 'Lampirkan opsi pembayaran tahunan dan diskon implementasi.', false),
    ('00000000-0000-4000-8000-000000000503', v_org_id, v_sales_two_id, 'Meeting evaluasi lisensi', 'Gojek Indonesia', current_date - 1, '10:00', 'Meeting', 'Tinggi', 'Tugas terlewat untuk follow-up kontrak lisensi.', false),
    ('00000000-0000-4000-8000-000000000504', v_org_id, v_admin_id, 'Follow-up integrasi Telegram', 'PT Astra International', current_date + 1, '09:30', 'Other', 'Rendah', 'Pastikan token bot dan chat ID sudah tersedia.', false)
  ON CONFLICT (id) DO UPDATE
    SET organization_id = EXCLUDED.organization_id,
        assignee_id = EXCLUDED.assignee_id,
        title = EXCLUDED.title,
        company = EXCLUDED.company,
        due_date = EXCLUDED.due_date,
        due_time = EXCLUDED.due_time,
        type = EXCLUDED.type,
        priority = EXCLUDED.priority,
        notes = EXCLUDED.notes,
        completed = EXCLUDED.completed,
        deleted_at = NULL,
        updated_at = now();

  INSERT INTO activities (id, organization_id, actor_id, actor_name, action, target, is_highlight, created_at)
  VALUES
    ('00000000-0000-4000-8000-000000000601', v_org_id, v_admin_id, 'Demo Admin', 'menambahkan kontak baru', 'Budi Wijaya', false, now() - interval '15 minutes'),
    ('00000000-0000-4000-8000-000000000602', v_org_id, v_sales_one_id, 'Budi Santoso', 'memindahkan deal ke meeting', 'Bank BCA', false, now() - interval '2 hours'),
    ('00000000-0000-4000-8000-000000000603', v_org_id, v_sales_two_id, 'Siti Rahma', 'memindahkan deal ke won', 'Gojek Indonesia', true, now() - interval '1 day')
  ON CONFLICT (id) DO UPDATE
    SET organization_id = EXCLUDED.organization_id,
        actor_id = EXCLUDED.actor_id,
        actor_name = EXCLUDED.actor_name,
        action = EXCLUDED.action,
        target = EXCLUDED.target,
        is_highlight = EXCLUDED.is_highlight,
        created_at = EXCLUDED.created_at;

  INSERT INTO performance_goals (id, organization_id, month, goal)
  VALUES
    ('00000000-0000-4000-8000-000000000701', v_org_id, date_trunc('month', current_date)::date, 500000000),
    ('00000000-0000-4000-8000-000000000702', v_org_id, (date_trunc('month', current_date) - interval '1 month')::date, 400000000),
    ('00000000-0000-4000-8000-000000000703', v_org_id, (date_trunc('month', current_date) - interval '2 months')::date, 350000000)
  ON CONFLICT (organization_id, month) DO UPDATE
    SET goal = EXCLUDED.goal,
        updated_at = now();

  INSERT INTO notifications (id, organization_id, user_id, title, message, created_at)
  VALUES
    ('00000000-0000-4000-8000-000000000801', v_org_id, v_admin_id, 'Deal Won!', 'Gojek Indonesia berhasil masuk tahap Deal Won.', now() - interval '1 day')
  ON CONFLICT (id) DO UPDATE
    SET organization_id = EXCLUDED.organization_id,
        user_id = EXCLUDED.user_id,
        title = EXCLUDED.title,
        message = EXCLUDED.message,
        read_at = NULL,
        created_at = EXCLUDED.created_at;
END;
$seed$;
