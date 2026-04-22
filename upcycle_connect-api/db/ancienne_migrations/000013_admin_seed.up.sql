INSERT INTO USER (id_user, email, password_user, first_name, last_name, upcycling_score, premium, status) VALUES
('7cd909a9-a27a-4849-a867-55e39b871c66', 'admin@upcycle.com', '$2a$10$JHayhaAb0AYcARtBUWkyuuxf2fN64DydnGgbHSVPQfA60uvY5mo5K', 'Admin', 'Upcycle', 0, 0, 'active');

INSERT INTO USER_ROLE (id_user, id_role) VALUES
('7cd909a9-a27a-4849-a867-55e39b871c66', '77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab');
