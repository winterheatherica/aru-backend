CREATE TABLE career_vacancies (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  title text NOT NULL,
  slug text UNIQUE NOT NULL,
  employment employment NOT NULL,
  location text NOT NULL,
  is_active boolean DEFAULT true,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  opened_at timestamptz DEFAULT now(),
  closed_at timestamptz,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_career_vacancies_employment ON career_vacancies(employment);
CREATE INDEX idx_career_vacancies_location ON career_vacancies(location);
CREATE INDEX idx_career_vacancies_is_active ON career_vacancies(is_active);
CREATE INDEX idx_career_vacancies_opened_at ON career_vacancies(opened_at);
