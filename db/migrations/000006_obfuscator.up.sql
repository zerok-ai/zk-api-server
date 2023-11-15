CREATE TABLE IF NOT EXISTS zk_obfuscation (
   id uuid NOT NULL DEFAULT uuid_generate_v4(),
   org_id character varying(255) NOT NULL,
   rule_name character varying(255) NOT NULL,
   rule_type character varying(255) NOT NULL,
   rule_def bytea,
   created_at timestamp without time zone NOT NULL,
   updated_at timestamp without time zone NOT NULL,
   deleted boolean DEFAULT false,
   disabled boolean DEFAULT false,
   PRIMARY KEY (id)
);
