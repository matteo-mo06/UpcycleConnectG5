SET FOREIGN_KEY_CHECKS=0;

ALTER TABLE payement_event                      DROP FOREIGN KEY fk_payement_event_payement;
ALTER TABLE payement_formation                  DROP FOREIGN KEY fk_payement_formation_payement;
ALTER TABLE payement_project                    DROP FOREIGN KEY fk_payement_project_payement;
ALTER TABLE user_formation_inscription_payement DROP FOREIGN KEY fk_ufip_payement;
ALTER TABLE user_payement                       DROP FOREIGN KEY fk_user_payement_payement;

ALTER TABLE payement                            MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;
ALTER TABLE payement_event                      MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;
ALTER TABLE payement_formation                  MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;
ALTER TABLE payement_project                    MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;
ALTER TABLE user_formation_inscription_payement MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;
ALTER TABLE user_payement                       MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;

ALTER TABLE payement_event                      ADD CONSTRAINT fk_payement_event_payement      FOREIGN KEY (id_payement) REFERENCES payement(id_payement);
ALTER TABLE payement_formation                  ADD CONSTRAINT fk_payement_formation_payement   FOREIGN KEY (id_payement) REFERENCES payement(id_payement);
ALTER TABLE payement_project                    ADD CONSTRAINT fk_payement_project_payement     FOREIGN KEY (id_payement) REFERENCES payement(id_payement);
ALTER TABLE user_formation_inscription_payement ADD CONSTRAINT fk_ufip_payement                 FOREIGN KEY (id_payement) REFERENCES payement(id_payement);
ALTER TABLE user_payement                       ADD CONSTRAINT fk_user_payement_payement        FOREIGN KEY (id_payement) REFERENCES payement(id_payement);

SET FOREIGN_KEY_CHECKS=1;
