-- Seed: Romanian HR / Recruiter interview questions
-- Categories:
--   acomodare              -> Acomodare
--   educatie_experienta    -> Educație / Experiență
--   abilitati_personale    -> Abilități și caracteristici personale

-- Acomodare
INSERT INTO questions (question_text, question_type, difficulty_level, expected_duration) VALUES
('Cum a fost drumul până aici? Ai mai fost în această clădire? Cum ți-ai început ziua?', 'acomodare', 'easy', 90);

-- Educație / Experiență
INSERT INTO questions (question_text, question_type, difficulty_level, expected_duration) VALUES
('Ce anume te-a determinat să alegi specializarea de Psihologie/Resurse Umane?', 'educatie_experienta', 'easy', 120);

-- Abilități și caracteristici personale
INSERT INTO questions (question_text, question_type, difficulty_level, expected_duration) VALUES
('Ce anume te atrage la rolul de asistent recrutor? Ce înseamnă pentru tine acest rol?', 'abilitati_personale', 'medium', 150),
('Care sunt așteptările tale de la rolul de junior în recrutare?', 'abilitati_personale', 'medium', 120),
('Îmi poți spune care sunt principalii pași dintr-un proces de recrutare?', 'abilitati_personale', 'medium', 180),
('Cum ți-ai organiza ziua atunci când ai avea mai multe interviuri?', 'abilitati_personale', 'medium', 150),
('Cum ai gestiona situația în care ai mai multe roluri deschise în același timp?', 'abilitati_personale', 'medium', 150),
('Cum ai învăța despre un rol pentru care ai de recrutat, însă nu știi nimic despre el?', 'abilitati_personale', 'medium', 150),
('Ce ai face atunci când nu înțelegi un termen HR sau o cerință de job?', 'abilitati_personale', 'easy', 120),
('Ai putea să îmi povestești o situație în care ai avut de transmis unui coleg/unei colege o informație mai greu de înțeles?', 'abilitati_personale', 'medium', 180),
('Ai putea să îmi spui o situație în care ai avut un coleg nou sau o colegă nouă la facultate și ați colaborat la un proiect?', 'abilitati_personale', 'medium', 180),
('Povestește-mi despre o situație în care ai lucrat la un proiect în echipă la facultate și cum v-ați organizat?', 'abilitati_personale', 'medium', 180),
('Ai putea să îmi povestești o situație de la facultate când ai încercat să rezolvi ceva, dar soluția ta nu a funcționat?', 'abilitati_personale', 'medium', 180),
('Ai putea să-mi povestești de un moment când un profesor de la facultate ți-a dat un proiect de făcut și ulterior a modificat cerințele cerute inițial?', 'abilitati_personale', 'medium', 180),
('Cum crezi că te pot ajuta cunoștințele acumulate până acum în rolul de recrutor?', 'abilitati_personale', 'medium', 150);
