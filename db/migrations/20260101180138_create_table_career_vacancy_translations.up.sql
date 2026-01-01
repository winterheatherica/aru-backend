CREATE TABLE career_vacancy_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  vacancy_id uuid NOT NULL REFERENCES career_vacancies(id) ON DELETE CASCADE,
  language language NOT NULL,
  description text,
  meta_title text,
  meta_description text,
  meta_keywords text[],
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (vacancy_id, language)
);

CREATE INDEX idx_career_vacancy_translations_language ON career_vacancy_translations(language);
CREATE INDEX idx_career_vacancy_translations_vacancy ON career_vacancy_translations(vacancy_id);
