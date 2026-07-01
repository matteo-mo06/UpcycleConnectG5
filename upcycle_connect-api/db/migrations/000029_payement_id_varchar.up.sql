SET FOREIGN_KEY_CHECKS=0;
ALTER TABLE payement                            MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;
ALTER TABLE payement_event                      MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;
ALTER TABLE payement_formation                  MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;
ALTER TABLE payement_project                    MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;
ALTER TABLE user_formation_inscription_payement MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;
ALTER TABLE user_payement                       MODIFY COLUMN id_payement VARCHAR(255) NOT NULL;
SET FOREIGN_KEY_CHECKS=1;
