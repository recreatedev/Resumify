-- RESUMES
CREATE TABLE resumes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id TEXT NOT NULL, -- from Clerk
  title TEXT NOT NULL,
  theme TEXT DEFAULT 'default',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- SECTIONS CONTROL
CREATE TABLE resume_sections (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  resume_id UUID NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
  name TEXT NOT NULL,        -- e.g. 'education', 'experience'
  display_name TEXT,
  is_visible BOOLEAN DEFAULT TRUE,
  order_index INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- EDUCATION
CREATE TABLE education (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  resume_id UUID NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
  institution TEXT,
  degree TEXT,
  field_of_study TEXT,
  start_date DATE,
  end_date DATE,
  grade TEXT,
  description TEXT,
  order_index INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- EXPERIENCE
CREATE TABLE experience (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  resume_id UUID NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
  company TEXT,
  position TEXT,
  start_date DATE,
  end_date DATE,
  location TEXT,
  description TEXT,
  order_index INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- PROJECTS
CREATE TABLE projects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  resume_id UUID NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
  name TEXT,
  role TEXT,
  description TEXT,
  link TEXT,
  technologies TEXT[],
  order_index INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- SKILLS
CREATE TABLE skills (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  resume_id UUID NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
  name TEXT,
  level TEXT, -- e.g. 'Beginner', 'Intermediate', 'Expert'
  category TEXT,
  order_index INT DEFAULT 0
);

-- CERTIFICATIONS
CREATE TABLE certifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  resume_id UUID NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
  name TEXT,
  organization TEXT,
  issue_date DATE,
  expiry_date DATE,
  credential_id TEXT,
  credential_url TEXT,
  order_index INT DEFAULT 0
);



-- Resume lookups by user
CREATE INDEX idx_resumes_user_id ON resumes(user_id);

-- Section ordering and visibility
CREATE INDEX idx_resume_sections_resume_id_order ON resume_sections(resume_id, order_index);
CREATE INDEX idx_resume_sections_resume_id_visible ON resume_sections(resume_id, is_visible);

-- Education ordering
CREATE INDEX idx_education_resume_id_order ON education(resume_id, order_index);

-- Experience ordering  
CREATE INDEX idx_experience_resume_id_order ON experience(resume_id, order_index);

-- Projects ordering
CREATE INDEX idx_projects_resume_id_order ON projects(resume_id, order_index);

-- Skills by category
CREATE INDEX idx_skills_resume_id_category ON skills(resume_id, category);
CREATE INDEX idx_skills_resume_id_order ON skills(resume_id, order_index);

-- Certifications ordering
CREATE INDEX idx_certifications_resume_id_order ON certifications(resume_id, order_index);