-- Demo data for customers table
INSERT INTO customers (id, name, status, email, phone, address, website, notes) VALUES
('c1a1e1b2-1111-4a1a-9a1a-111111111111', 'Acme Corp', 'active', 'info@acme.com', '+1-555-1234', '123 Main St, Springfield, USA', 'https://acme.com', '- **Key client**
- Prefers email contact
- [x] Onboarded

---

# More info

> Acme Corp is a strategic partner.');
INSERT INTO customers (id, name, status, email, phone, address, website, notes) VALUES
('c2b2e2b3-2222-4b2b-9b2b-222222222222', 'Globex Inc', 'active', 'contact@globex.com', '+1-555-5678', '456 Elm St, Metropolis, USA', 'https://globex.com', '- *Recently onboarded*
- [ ] Needs follow-up

---

## Todo List

- [x] Signed contract
- [ ] Schedule kickoff call');
INSERT INTO customers (id, name, status, email, phone, address, website, notes) VALUES
('c3c3e3b4-3333-4c3c-9c3c-333333333333', 'Initech',  'inactive', 'support@initech.com', '+1-555-8765', '789 Oak St, Gotham, USA', 'https://initech.com', '- [ ] Account on hold
- *Pending review*

---

### Links

See [support docs](https://initech.com/support) for more.');
INSERT INTO customers (id, name, status, email, phone, address, website, notes) VALUES
('c4d4e4b5-4444-4d4d-9d4d-444444444444', 'Umbrella Corp', 'active', 'admin@umbrella.com', '+1-555-4321', '321 Pine St, Raccoon City, USA', 'https://umbrella.com', '- **VIP**
- [x] Priority support

---

## Services

| Service | Status |
|---------|--------|
| Hosting | Active |
| Backup  | Active |');
INSERT INTO customers (id, name, status, email, phone, address, website, notes) VALUES
('c5e5e5b6-5555-4e5e-9e5e-555555555555', 'Wayne Enterprises', 'active', 'hello@wayne.com', '+1-555-2468', '1007 Mountain Dr, Gotham, USA', 'https://wayne.com', '- Long-term partner
- [x] Trusted
- [ ] Review Q3 goals

---

> Wayne Enterprises is a key account.');

-- Demo data for contacts table
INSERT INTO contacts (id, customer_id, name, role, email, phone, avatar, is_primary, notes) VALUES
('a1b2c3d4-aaaa-4aaa-8aaa-aaaaaaaaaaaa', 'c1a1e1b2-1111-4a1a-9a1a-111111111111', 'Jane Smith', 'CTO', 'jane.smith@acme.com', '+1-555-1111', 'jane.png', 1, 'Main technical contact');
INSERT INTO contacts (id, customer_id, name, role, email, phone, avatar, is_primary, notes) VALUES
('b2c3d4e5-bbbb-4bbb-8bbb-bbbbbbbbbbbb', 'c2b2e2b3-2222-4b2b-9b2b-222222222222', 'John Doe', 'CEO', 'john.doe@globex.com', '+1-555-2222', 'john.png', 1, 'Decision maker');
INSERT INTO contacts (id, customer_id, name, role, email, phone, avatar, is_primary, notes) VALUES
('c3d4e5f6-cccc-4ccc-8ccc-cccccccccccc', 'c3c3e3b4-3333-4c3c-9c3c-333333333333', 'Alice Johnson', 'Support Lead', 'alice.j@initech.com', '+1-555-3333', 'alice.png', 1, 'Handles support tickets');
INSERT INTO contacts (id, customer_id, name, role, email, phone, avatar, is_primary, notes) VALUES
('d4e5f6a7-dddd-4ddd-8ddd-dddddddddddd', 'c4d4e4b5-4444-4d4d-9d4d-444444444444', 'Bob Brown', 'Operations', 'bob.brown@umbrella.com', '+1-555-4444', 'bob.png', 0, 'Secondary contact');
INSERT INTO contacts (id, customer_id, name, role, email, phone, avatar, is_primary, notes) VALUES
('e5f6a7b8-eeee-4eee-8eee-eeeeeeeeeeee', 'c5e5e5b6-5555-4e5e-9e5e-555555555555', 'Rachel Green', 'CFO', 'rachel.green@wayne.com', '+1-555-5555', 'rachel.png', 1, 'Finance contact');

-- Demo data for activity_log table
INSERT INTO activity_log (id, customer_id, activity_type, action, description, created_at) VALUES
('f1a2b3c4-aaaa-4aaa-8aaa-aaaaaaaaaaaa', 'c1a1e1b2-1111-4a1a-9a1a-111111111111', 'project', 'project_created', 'New project created for Acme Corp', datetime('now', '-2 hours'));
INSERT INTO activity_log (id, customer_id, activity_type, action, description, created_at) VALUES
('f2b3c4d5-bbbb-4bbb-8bbb-bbbbbbbbbbbb', 'c2b2e2b3-2222-4b2b-9b2b-222222222222', 'subscription', 'subscription_renewed', 'TechStart Inc renewed their subscription', datetime('now', '-5 hours'));
INSERT INTO activity_log (id, customer_id, activity_type, action, description, created_at) VALUES
('f3c4d5e6-cccc-4ccc-8ccc-cccccccccccc', 'c3c3e3b4-3333-4c3c-9c3c-333333333333', 'contact', 'contact_added', 'New contact added for Global Systems', datetime('now', '-1 day'));
INSERT INTO activity_log (id, customer_id, activity_type, action, description, created_at) VALUES
('f4d5e6f7-dddd-4ddd-8ddd-dddddddddddd', 'c4d4e4b5-4444-4d4d-9d4d-444444444444', 'invoice', 'invoice_paid', 'DevCorp paid invoice #5678', datetime('now', '-2 days'));
INSERT INTO activity_log (id, customer_id, activity_type, action, description, created_at) VALUES
('f5e6a7b8-eeee-4eee-8eee-eeeeeeeeeeee', 'c5e5e5b6-5555-4e5e-9e5e-555555555555', 'project', 'project_completed', 'Project completed for Wayne Enterprises', datetime('now', '-3 days'));
