ALTER TABLE invitation
ADD COLUMN IF NOT EXISTS exp TIMESTAMP(0) WITH TIME ZONE NOT NULL;
