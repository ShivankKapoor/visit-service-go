-- Requires the uuid-ossp extension: CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.daily_visit_stats (
	summary_date date NOT NULL,
	total_visits int8 NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	CONSTRAINT daily_visit_stats_pkey PRIMARY KEY (summary_date)
);

CREATE TABLE public.page_visits (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	ip_address text NOT NULL,
	page_visited text NOT NULL,
	device_info text NULL,
	user_agent text NULL,
	"timestamp" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT page_visits_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_ip_address ON public.page_visits USING btree (ip_address);
CREATE INDEX idx_page_ip_combo ON public.page_visits USING btree (page_visited, ip_address);
CREATE INDEX idx_page_visited ON public.page_visits USING btree (page_visited);
CREATE INDEX idx_page_visits_timestamp_asc ON public.page_visits USING btree ("timestamp");
CREATE INDEX idx_page_visits_timestamp_desc ON public.page_visits USING btree ("timestamp" DESC);
