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
INSERT INTO customers (id, name, status, email, phone, address, website, notes) VALUES
('c6f6f6b7-6666-4f6f-9f6f-666666666666', 'Stark Industries', 'active', 'contact@stark.com', '+1-555-1112', '200 Park Ave, New York, USA', 'https://stark.com', '**VIP client**
> "Changing the world, one suit at a time."
- Innovation leader
- Prefers quarterly reviews
- See [project archive](https://stark.com/projects)
');
INSERT INTO customers (id, name, status, email, phone, address, website, notes) VALUES
('c7a7a7b8-7777-4a7a-9a7a-777777777777', 'Hooli', 'active', 'info@hooli.com', '+1-555-2223', '500 Tech Rd, Silicon Valley, USA', 'https://hooli.com', '### Internal Notes
- *Disruptive*
- Interested in cloud migration
- Contact: Gavin Belson
---
| Priority | Next Step      |
|----------|----------------|
| High     | Schedule demo  |
');
INSERT INTO customers (id, name, status, email, phone, address, website, notes) VALUES
('c8b8b8b9-8888-4b8b-9b8b-888888888888', 'Soylent Corp', 'inactive', 'hello@soylent.com', '+1-555-3334', '404 Green St, New York, USA', 'https://soylent.com', '> *"It''s people!"*
- Account inactive
- Awaiting compliance review
- [Sustainability report](https://soylent.com/sustainability)
');
INSERT INTO customers (id, name, status, email, phone, address, website, notes) VALUES
('c9c9c9c0-9999-4c9c-9c9c-999999999999', 'Cyberdyne Systems', 'active', 'admin@cyberdyne.com', '+1-555-4445', '101 AI Blvd, Los Angeles, USA', 'https://cyberdyne.com', '**Strategic partner**
- Robotics division
- *Sensitive data protocols required*
- See [AI policy](https://cyberdyne.com/ai-policy)
');
INSERT INTO customers (id, name, status, email, phone, address, website, notes) VALUES
('c0d0d0d1-0000-4d0d-9d0d-000000000000', 'Tyrell Corporation', 'active', 'contact@tyrell.com', '+1-555-5556', '1 Replicant Way, Los Angeles, USA', 'https://tyrell.com', '#### Replicant Program
- New client
- *Requires NDA*
> "More human than human."
- [Intro deck](https://tyrell.com/intro)
');

-- Demo data for contacts table
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('a1b2c3d4-aaaa-4aaa-8aaa-aaaaaaaaaaaa', 'c1a1e1b2-1111-4a1a-9a1a-111111111111', 'Jane Smith', 'CTO', 'jane.smith@acme.com', '+1-555-1111', 1, '- **Main technical contact**
- Handles escalations
- [Profile](https://acme.com/team/jane)
> "Always available for urgent issues."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('b2c3d4e5-bbbb-4bbb-8bbb-bbbbbbbbbbbb', 'c2b2e2b3-2222-4b2b-9b2b-222222222222', 'John Doe', 'CEO', 'john.doe@globex.com', '+1-555-2222', 1, '- Decision maker
- *Strategic planning*
- [LinkedIn](https://linkedin.com/in/johndoe)
> "Vision for growth."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('c3d4e5f6-cccc-4ccc-8ccc-cccccccccccc', 'c3c3e3b4-3333-4c3c-9c3c-333333333333', 'Alice Johnson', 'Support Lead', 'alice.j@initech.com', '+1-555-3333', 1, '- Handles support tickets
- *Customer advocate*
- [Support Portal](https://initech.com/support)
> "Always ready to help."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('d4e5f6a7-dddd-4ddd-8ddd-dddddddddddd', 'c4d4e4b5-4444-4d4d-9d4d-444444444444', 'Bob Brown', 'Operations', 'bob.brown@umbrella.com', '+1-555-4444', 0, '- Secondary contact
- *Logistics coordinator*
- [Ops Dashboard](https://umbrella.com/ops)
> "Keeps things running."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('e5f6a7b8-eeee-4eee-8eee-eeeeeeeeeeee', 'c5e5e5b6-5555-4e5e-9e5e-555555555555', 'Rachel Green', 'CFO', 'rachel.green@wayne.com', '+1-555-5555', 1, '- Finance contact
- *Budget planning*
- [Finance Portal](https://wayne.com/finance)
> "Keeps the books balanced."');

INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('a1b2c3d4-aaaa-4aaa-8aaa-aaaaaaaaaaab', 'c1a1e1b2-1111-4a1a-9a1a-111111111111', 'Tom Allen', 'Project Manager', 'tom.allen@acme.com', '+1-555-1112', 0, '- *Coordinates project timelines*
- Oversees deliverables
- [Project dashboard](https://acme.com/projects)
> "Keeps everyone on track."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('a1b2c3d4-aaaa-4aaa-8aaa-aaaaaaaaaaac', 'c1a1e1b2-1111-4a1a-9a1a-111111111111', 'Lisa Wu', 'Support Engineer', 'lisa.wu@acme.com', '+1-555-1113', 0, '- First responder for support tickets
- **Expert in cloud systems**
- [Support docs](https://acme.com/support)
> "Solves problems fast."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('b2c3d4e5-bbbb-4bbb-8bbb-bbbbbbbbbbbc', 'c2b2e2b3-2222-4b2b-9b2b-222222222222', 'Maria Garcia', 'Account Manager', 'maria.garcia@globex.com', '+1-555-2223', 0, '- Manages client relationships
- [Contact Maria](mailto:maria.garcia@globex.com)
- **Excellent communicator**
> "Your go-to for account questions."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('b2c3d4e5-bbbb-4bbb-8bbb-bbbbbbbbbbbd', 'c2b2e2b3-2222-4b2b-9b2b-222222222222', 'Steve Kim', 'Technical Lead', 'steve.kim@globex.com', '+1-555-2224', 0, '- Oversees tech stack
- *DevOps specialist*
- [Tech Wiki](https://globex.com/wiki)
> "Ensures uptime."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('c3d4e5f6-cccc-4ccc-8ccc-cccccccccccd', 'c3c3e3b4-3333-4c3c-9c3c-333333333333', 'Mark Patel', 'Product Owner', 'mark.patel@initech.com', '+1-555-3334', 0, '- Guides product vision
- [Product Roadmap](https://initech.com/roadmap)
- **Agile expert**
> "Drives innovation."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('c3d4e5f6-cccc-4ccc-8ccc-ccccccccccce', 'c3c3e3b4-3333-4c3c-9c3c-333333333333', 'Emily Chen', 'QA Specialist', 'emily.chen@initech.com', '+1-555-3335', 0, '- Tests releases
- *Quality assurance*
- [QA Docs](https://initech.com/qa)
> "Finds the bugs before you do."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('d4e5f6a7-dddd-4ddd-8ddd-ddddddddddde', 'c4d4e4b5-4444-4d4d-9d4d-444444444444', 'Sarah Lee', 'Legal Counsel', 'sarah.lee@umbrella.com', '+1-555-4445', 0, '- Handles contracts
- [Legal Docs](https://umbrella.com/legal)
- **Compliance expert**
> "Ensures everything is above board."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('d4e5f6a7-dddd-4ddd-8ddd-dddddddddddf', 'c4d4e4b5-4444-4d4d-9d4d-444444444444', 'Mike Davis', 'IT Security', 'mike.davis@umbrella.com', '+1-555-4446', 0, '- Manages security protocols
- *Incident response*
- [Security Policy](https://umbrella.com/security)
> "Protects your data."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('e5f6a7b8-eeee-4eee-8eee-eeeeeeeeeeef', 'c5e5e5b6-5555-4e5e-9e5e-555555555555', 'Bruce Wayne', 'CEO', 'bruce.wayne@wayne.com', '+1-555-5556', 0, '- Executive decisions
- [Profile](https://wayne.com/bruce)
- **Visionary leader**
> "Always in the boardroom."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('e5f6a7b8-eeee-4eee-8eee-eeeeeeeeeeeg', 'c5e5e5b6-5555-4e5e-9e5e-555555555555', 'Alfred Pennyworth', 'Executive Assistant', 'alfred@wayne.com', '+1-555-5557', 0, '- Schedules meetings
- *Trusted confidant*
- [Contact Alfred](mailto:alfred@wayne.com)
> "Handles everything with grace."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('f6a7b8c9-ffff-4fff-8fff-fffffffffff1', 'c6f6f6b7-6666-4f6f-9f6f-666666666666', 'Tony Stark', 'CEO', 'tony.stark@stark.com', '+1-555-6661', 1, '- Leads innovation
- [Profile](https://stark.com/tony)
- **Inventor**
> "Genius, billionaire, playboy, philanthropist."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('f6a7b8c9-ffff-4fff-8fff-fffffffffff2', 'c6f6f6b7-6666-4f6f-9f6f-666666666666', 'Pepper Potts', 'COO', 'pepper.potts@stark.com', '+1-555-6662', 0, '- Manages operations
- [Contact Pepper](mailto:pepper@stark.com)
- *Organizational expert*
> "Keeps Tony on schedule."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('f6a7b8c9-ffff-4fff-8fff-fffffffffff3', 'c6f6f6b7-6666-4f6f-9f6f-666666666666', 'Happy Hogan', 'Head of Security', 'happy.hogan@stark.com', '+1-555-6663', 0, '- Security lead
- [Security Team](https://stark.com/security)
- **Trusted protector**
> "Always on guard."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('a7b8c9d0-aaaa-4aaa-8aaa-aaaaaaaaaaa1', 'c7a7a7b8-7777-4a7a-9a7a-777777777777', 'Gavin Belson', 'CEO', 'gavin.belson@hooli.com', '+1-555-7771', 1, '- Company vision
- [Profile](https://hooli.com/gavin)
- *Disruptive thinker*
> "Making the world a better place."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('a7b8c9d0-aaaa-4aaa-8aaa-aaaaaaaaaaa2', 'c7a7a7b8-7777-4a7a-9a7a-777777777777', 'Monica Hall', 'Product Manager', 'monica.hall@hooli.com', '+1-555-7772', 0, '- Oversees product launches
- [Product Docs](https://hooli.com/products)
- **Detail-oriented**
> "Ensures product success."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('a7b8c9d0-aaaa-4aaa-8aaa-aaaaaaaaaaa3', 'c7a7a7b8-7777-4a7a-9a7a-777777777777', 'Jared Dunn', 'Operations', 'jared.dunn@hooli.com', '+1-555-7773', 0, '- Manages day-to-day ops
- [Ops Wiki](https://hooli.com/ops)
- *Process improvement*
> "Optimizes everything."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('b8b9c0d1-bbbb-4bbb-8bbb-bbbbbbbbbbb1', 'c8b8b8b9-8888-4b8b-9b8b-888888888888', 'Frank Thorn', 'Compliance Officer', 'frank.thorn@soylent.com', '+1-555-8881', 1, '- Ensures regulatory compliance
- [Compliance Docs](https://soylent.com/compliance)
- *Detail-focused*
> "Keeps us legal."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('b8b9c0d1-bbbb-4bbb-8bbb-bbbbbbbbbbb2', 'c8b8b8b9-8888-4b8b-9b8b-888888888888', 'Tab Fielding', 'Operations', 'tab.fielding@soylent.com', '+1-555-8882', 0, '- Handles logistics
- [Ops Portal](https://soylent.com/ops)
- **Efficient manager**
> "Smooth operations."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('b8b9c0d1-bbbb-4bbb-8bbb-bbbbbbbbbbb3', 'c8b8b8b9-8888-4b8b-9b8b-888888888888', 'Shirl', 'Customer Success', 'shirl@soylent.com', '+1-555-8883', 0, '- Customer onboarding
- [Success Stories](https://soylent.com/success)
- *Empathetic communicator*
> "Clients love her."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('c9c0d1e2-cccc-4ccc-8ccc-ccccccccccf1', 'c9c9c9c0-9999-4c9c-9c9c-999999999999', 'Miles Dyson', 'Director of R&D', 'miles.dyson@cyberdyne.com', '+1-555-9991', 1, '- Leads research
- [R&D Wiki](https://cyberdyne.com/rd)
- **Innovator**
> "Building the future."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('c9c0d1e2-cccc-4ccc-8ccc-ccccccccccf2', 'c9c9c9c0-9999-4c9c-9c9c-999999999999', 'Sarah Connor', 'Security Consultant', 'sarah.connor@cyberdyne.com', '+1-555-9992', 0, '- Security protocols
- [Security Docs](https://cyberdyne.com/security)
- *Risk management*
> "Prepared for anything."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('c9c0d1e2-cccc-4ccc-8ccc-ccccccccccf3', 'c9c9c9c0-9999-4c9c-9c9c-999999999999', 'John Connor', 'AI Ethics', 'john.connor@cyberdyne.com', '+1-555-9993', 0, '- Oversees AI policy
- [AI Policy](https://cyberdyne.com/ai-policy)
- **Ethics advocate**
> "Protects humanity."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('d0d1e2f3-dddd-4ddd-8ddd-ddddddddddf1', 'c0d0d0d1-0000-4d0d-9d0d-000000000000', 'Eldon Tyrell', 'CEO', 'eldon.tyrell@tyrell.com', '+1-555-0001', 1, '- Executive decisions
- [Profile](https://tyrell.com/eldon)
- *Visionary*
> "More human than human."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('d0d1e2f3-dddd-4ddd-8ddd-ddddddddddf2', 'c0d0d0d1-0000-4d0d-9d0d-000000000000', 'Rachael', 'Replicant Liaison', 'rachael@tyrell.com', '+1-555-0002', 0, '- Client relations
- [Liaison Info](https://tyrell.com/liaison)
- **Empathetic**
> "Understands replicants."');
INSERT INTO contacts (id, customer_id, name, role, email, phone, is_primary, notes) VALUES
('d0d1e2f3-dddd-4ddd-8ddd-ddddddddddf3', 'c0d0d0d1-0000-4d0d-9d0d-000000000000', 'Leon Kowalski', 'Field Agent', 'leon.kowalski@tyrell.com', '+1-555-0003', 0, '- Field operations
- [Agent Profile](https://tyrell.com/leon)
- *Problem solver*
> "Gets the job done."');

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

-- Demo data for subscriptions table
INSERT INTO subscriptions (id, customer_id, description, amount, term, billing_cadence, start_date, end_date, status, created_at, updated_at)
VALUES
('s1a1e1b2-1111-4a1a-9a1a-111111111111', 'c1a1e1b2-1111-4a1a-9a1a-111111111111', 'Acme SaaS Platform', 99.00, 'monthly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s1a1e1b2-1111-4a1a-9a1a-111111111112', 'c1a1e1b2-1111-4a1a-9a1a-111111111111', 'Acme Premium Support', 1200.00, 'yearly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s1a1e1b2-1111-4a1a-9a1a-111111111113', 'c1a1e1b2-1111-4a1a-9a1a-111111111111', 'Acme Data Backup', 300.00, 'yearly', 'yearly', '2024-01-01T00:00:00Z', NULL, 'paused', datetime('now'), datetime('now')),
('s2b2e2b3-2222-4b2b-9b2b-222222222221', 'c2b2e2b3-2222-4b2b-9b2b-222222222222', 'Globex Cloud Storage', 250.00, 'yearly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s2b2e2b3-2222-4b2b-9b2b-222222222222', 'c2b2e2b3-2222-4b2b-9b2b-222222222222', 'Globex API Access', 50.00, 'monthly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s2b2e2b3-2222-4b2b-9b2b-222222222223', 'c2b2e2b3-2222-4b2b-9b2b-222222222222', 'Globex Analytics', 600.00, 'yearly', 'yearly', '2024-01-01T00:00:00Z', NULL, 'cancelled', datetime('now'), datetime('now')),
('s3c3e3b4-3333-4c3c-9c3c-333333333331', 'c3c3e3b4-3333-4c3c-9c3c-333333333333', 'Initech Support Retainer', 500.00, 'yearly', 'monthly', '2024-01-01T00:00:00Z', '2024-12-31T00:00:00Z', 'active', datetime('now'), datetime('now')),
('s3c3e3b4-3333-4c3c-9c3c-333333333332', 'c3c3e3b4-3333-4c3c-9c3c-333333333333', 'Initech DevOps', 75.00, 'monthly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s3c3e3b4-3333-4c3c-9c3c-333333333333', 'c3c3e3b4-3333-4c3c-9c3c-333333333333', 'Initech QA', 200.00, 'yearly', 'yearly', '2024-01-01T00:00:00Z', NULL, 'paused', datetime('now'), datetime('now')),
('s4d4e4b5-4444-4d4d-9d4d-444444444441', 'c4d4e4b5-4444-4d4d-9d4d-444444444444', 'Umbrella Hosting', 400.00, 'yearly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s4d4e4b5-4444-4d4d-9d4d-444444444442', 'c4d4e4b5-4444-4d4d-9d4d-444444444444', 'Umbrella Backup', 100.00, 'monthly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'paused', datetime('now'), datetime('now')),
('s5e5e5b6-5555-4e5e-9e5e-555555555551', 'c5e5e5b6-5555-4e5e-9e5e-555555555555', 'Wayne Cloud Suite', 150.00, 'monthly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s5e5e5b6-5555-4e5e-9e5e-555555555552', 'c5e5e5b6-5555-4e5e-9e5e-555555555555', 'Wayne Security Monitoring', 1800.00, 'yearly', 'yearly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s5e5e5b6-5555-4e5e-9e5e-555555555553', 'c5e5e5b6-5555-4e5e-9e5e-555555555555', 'Wayne Data Analytics', 500.00, 'quarterly', 'quarterly', '2024-01-01T00:00:00Z', NULL, 'cancelled', datetime('now'), datetime('now')),
('s6f6f6b7-6666-4f6f-9f6f-666666666661', 'c6f6f6b7-6666-4f6f-9f6f-666666666666', 'Stark Innovation Platform', 250.00, 'monthly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s6f6f6b7-6666-4f6f-9f6f-666666666662', 'c6f6f6b7-6666-4f6f-9f6f-666666666666', 'Stark R&D Retainer', 5000.00, 'yearly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s6f6f6b7-6666-4f6f-9f6f-666666666663', 'c6f6f6b7-6666-4f6f-9f6f-666666666666', 'Stark VIP Support', 1200.00, 'yearly', 'yearly', '2024-01-01T00:00:00Z', NULL, 'paused', datetime('now'), datetime('now')),
('s6f6f6b7-6666-4f6f-9f6f-666666666664', 'c6f6f6b7-6666-4f6f-9f6f-666666666666', 'Stark IoT Monitoring', 300.00, 'quarterly', 'quarterly', '2024-01-01T00:00:00Z', NULL, 'cancelled', datetime('now'), datetime('now')),
('s7a7a7b8-7777-4a7a-9a7a-777777777771', 'c7a7a7b8-7777-4a7a-9a7a-777777777777', 'Hooli Cloud Migration', 800.00, 'yearly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s8b8b8b9-8888-4b8b-9b8b-888888888881', 'c8b8b8b9-8888-4b8b-9b8b-888888888888', 'Soylent Compliance Suite', 200.00, 'monthly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s8b8b8b9-8888-4b8b-9b8b-888888888882', 'c8b8b8b9-8888-4b8b-9b8b-888888888888', 'Soylent Sustainability Reports', 900.00, 'yearly', 'yearly', '2024-01-01T00:00:00Z', NULL, 'cancelled', datetime('now'), datetime('now')),
('s9c9c9c0-9999-4c9c-9c9c-999999999991', 'c9c9c9c0-9999-4c9c-9c9c-999999999999', 'Cyberdyne Robotics Platform', 1000.00, 'yearly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s9c9c9c0-9999-4c9c-9c9c-999999999992', 'c9c9c9c0-9999-4c9c-9c9c-999999999999', 'Cyberdyne AI Ethics Review', 300.00, 'quarterly', 'quarterly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s9c9c9c0-9999-4c9c-9c9c-999999999993', 'c9c9c9c0-9999-4c9c-9c9c-999999999999', 'Cyberdyne Data Security', 150.00, 'monthly', 'monthly', '2024-01-01T00:00:00Z', NULL, 'paused', datetime('now'), datetime('now')),
('s0d0d0d1-0000-4d0d-9d0d-000000000001', 'c0d0d0d1-0000-4d0d-9d0d-000000000000', 'Tyrell Replicant Program', 2000.00, 'yearly', 'yearly', '2024-01-01T00:00:00Z', NULL, 'active', datetime('now'), datetime('now')),
('s0d0d0d1-0000-4d0d-9d0d-000000000002', 'c0d0d0d1-0000-4d0d-9d0d-000000000000', 'Tyrell Compliance Monitoring', 400.00, 'quarterly', 'quarterly', '2024-01-01T00:00:00Z', NULL, 'cancelled', datetime('now'), datetime('now'));
