-- 顧客レコードを直接作成
INSERT INTO customers (id, user_id, first_name, last_name, phone, birth_date, gender, created_at, updated_at)
VALUES (
  '$(CUSTOMER_ID)',
  '$(USER_ID)',
  '太郎',
  '山田',
  '090-1234-5678',
  '1990-01-15',
  'male',
  NOW(),
  NOW()
);
