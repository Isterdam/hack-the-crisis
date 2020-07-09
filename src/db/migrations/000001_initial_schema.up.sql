--
-- Name: booking_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.booking_id_seq;

--
-- Name: bookings; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE IF NOT EXISTS public.bookings (
    id integer DEFAULT nextval('public.booking_id_seq'::regclass) PRIMARY KEY NOT NULL,
    slot_id integer NOT NULL,
    phone_number character varying(50) NOT NULL,
    code character varying(50) NOT NULL,
    first_name character varying(50) NOT NULL,
    last_name character varying(50) NOT NULL,
    visitee character varying(50) NOT NULL,
    message text DEFAULT ''::text,
    status character varying(50) DEFAULT ''::character varying
);

--
-- Name: company_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.company_id_seq;

--
-- Name: company; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE IF NOT EXISTS public.company (
    id integer DEFAULT nextval('public.company_id_seq'::regclass) PRIMARY KEY NOT NULL,
    name character varying(50) NOT NULL,
    adress character varying(50) DEFAULT ' '::character varying NOT NULL,
    city character varying(50) DEFAULT ' '::character varying NOT NULL,
    country character varying(50) DEFAULT ' '::character varying NOT NULL,
    post_code character varying(50) DEFAULT ' '::character varying NOT NULL,
    contact_firstname character varying(50) DEFAULT ' '::character varying NOT NULL,
    contact_number character varying(50) DEFAULT ' '::character varying NOT NULL,
    contact_lastname character varying(50) DEFAULT ' '::character varying NOT NULL,
    verified boolean DEFAULT false NOT NULL,
    email character varying(50) DEFAULT ' '::character varying NOT NULL,
    password character varying(200) DEFAULT ' '::character varying NOT NULL,
    lon numeric,
    lat numeric,
    contact_email character varying(50)
);

--
-- Name: slot_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.slot_id_seq;

--
-- Name: slots; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE IF NOT EXISTS public.slots (
    id integer DEFAULT nextval('public.slot_id_seq'::regclass) PRIMARY KEY NOT NULL,
    start_time timestamp with time zone NOT NULL,
    end_time timestamp with time zone NOT NULL,
    max integer NOT NULL,
    company_id integer NOT NULL,
    booked integer DEFAULT 0 NOT NULL
);

--
-- Name: slots FK_20; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.slots
    ADD CONSTRAINT "FK_20" FOREIGN KEY (company_id) REFERENCES public.company(id);


ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT fk_20 FOREIGN KEY (slot_id) REFERENCES public.slots(id) ON DELETE CASCADE;
