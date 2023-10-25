--
-- PostgreSQL database dump
--

-- Dumped from database version 16.0
-- Dumped by pg_dump version 16.0

-- Started on 2023-10-24 23:08:16

--
-- TOC entry 4863 (class 1262 OID 16455)
-- Name: test-service-database; Type: DATABASE; Schema: -; Owner: -
--

CREATE DATABASE "test-service-database";


\c "test-service-database";

--
-- TOC entry 220 (class 1259 OID 16508)
-- Name: disaster_list_table; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.disaster_list_table (
    disaster_id integer NOT NULL,
    disaster_type_id integer NOT NULL,
    location_id integer NOT NULL,
    timestamp_start timestamp without time zone NOT NULL,
    timestamp_end timestamp without time zone NOT NULL,
    severity character varying(50) NOT NULL,
    description text NOT NULL
);


--
-- TOC entry 219 (class 1259 OID 16507)
-- Name: disaster_list_table_disaster_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.disaster_list_table_disaster_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
    CYCLE;


--
-- TOC entry 4864 (class 0 OID 0)
-- Dependencies: 219
-- Name: disaster_list_table_disaster_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.disaster_list_table_disaster_id_seq OWNED BY public.disaster_list_table.disaster_id;


--
-- TOC entry 216 (class 1259 OID 16471)
-- Name: disaster_type_table; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.disaster_type_table (
    disaster_type_id integer NOT NULL,
    disaster_name character varying(255) NOT NULL,
    disaster_description text NOT NULL
);


--
-- TOC entry 218 (class 1259 OID 16489)
-- Name: location_table; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.location_table (
    location_id integer NOT NULL,
    country character varying(255) NOT NULL,
    city character varying(255) NOT NULL,
    longitude numeric(10,4) NOT NULL,
    latitude numeric(10,4) NOT NULL
);


--
-- TOC entry 217 (class 1259 OID 16488)
-- Name: location_table_location_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.location_table_location_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
    CYCLE;


--
-- TOC entry 4865 (class 0 OID 0)
-- Dependencies: 217
-- Name: location_table_location_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.location_table_location_id_seq OWNED BY public.location_table.location_id;


--
-- TOC entry 215 (class 1259 OID 16470)
-- Name: natural_disaster_type_table_disaster_type_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.natural_disaster_type_table_disaster_type_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
    CYCLE;


--
-- TOC entry 4866 (class 0 OID 0)
-- Dependencies: 215
-- Name: natural_disaster_type_table_disaster_type_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.natural_disaster_type_table_disaster_type_id_seq OWNED BY public.disaster_type_table.disaster_type_id;


--
-- TOC entry 4700 (class 2604 OID 16511)
-- Name: disaster_list_table disaster_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.disaster_list_table ALTER COLUMN disaster_id SET DEFAULT nextval('public.disaster_list_table_disaster_id_seq'::regclass);


--
-- TOC entry 4698 (class 2604 OID 16497)
-- Name: disaster_type_table disaster_type_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.disaster_type_table ALTER COLUMN disaster_type_id SET DEFAULT nextval('public.natural_disaster_type_table_disaster_type_id_seq'::regclass);


--
-- TOC entry 4699 (class 2604 OID 16492)
-- Name: location_table location_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.location_table ALTER COLUMN location_id SET DEFAULT nextval('public.location_table_location_id_seq'::regclass);


--
-- TOC entry 4857 (class 0 OID 16508)
-- Dependencies: 220
-- Data for Name: disaster_list_table; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- TOC entry 4853 (class 0 OID 16471)
-- Dependencies: 216
-- Data for Name: disaster_type_table; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (1, 'Earthquake', 'Sudden shaking of the Earth''s surface due to tectonic movements.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (2, 'Hurricane', 'A violent storm with strong winds and heavy rain that can cause widespread damage.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (3, 'Tornado', 'A destructive rotating column of air characterized by a funnel shape.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (4, 'Flood', 'An overflow of water onto normally dry land, often caused by heavy rainfall.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (5, 'Wildfire', 'An uncontrolled fire that spreads rapidly through vegetation.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (6, 'Volcanic Eruption', 'The release of molten lava, ash, and gases from a volcano.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (7, 'Drought', 'A prolonged period of abnormally low rainfall leading to water shortages.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (8, 'Blizzard', 'A severe snowstorm with strong winds and reduced visibility.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (9, 'Tsunami', 'A series of large ocean waves caused by underwater seismic activity.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (10, 'Landslide', 'The sudden movement of rock, soil, and debris down a slope.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (11, 'Heatwave', 'A prolonged period of excessively hot weather.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (12, 'Avalanche', 'A mass of snow, ice, and debris rapidly descending down a mountain.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (13, 'Hailstorm', 'A severe weather event with hailstones that can damage crops and property.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (14, 'Cyclone', 'A large-scale weather system characterized by strong winds and heavy rain.') ON CONFLICT DO NOTHING;
INSERT INTO public.disaster_type_table (disaster_type_id, disaster_name, disaster_description) VALUES (15, 'Sinkhole', 'A depression or hole in the ground caused by the collapse of underlying material.') ON CONFLICT DO NOTHING;


--
-- TOC entry 4855 (class 0 OID 16489)
-- Dependencies: 218
-- Data for Name: location_table; Type: TABLE DATA; Schema: public; Owner: -
--


INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (1, 'Moldova', 'Chisinau', 28.8353, 47.0228) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (2, 'Moldova', 'Tiraspol', 29.6433, 46.8403) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (3, 'Moldova', 'Balti', 27.9167, 47.7667) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (4, 'Moldova', 'Bender', 29.4833, 46.8333) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (5, 'Moldova', 'Ungheni', 27.8167, 47.2167) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (6, 'Moldova', 'Cahul', 28.1836, 45.9167) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (7, 'Moldova', 'Soroca', 28.3000, 48.1667) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (8, 'Moldova', 'Orhei', 28.8167, 47.3833) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (9, 'Moldova', 'Comrat', 28.6667, 46.3167) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (10, 'Moldova', 'Straseni', 28.6167, 47.1333) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (11, 'Moldova', 'Causeni', 29.4000, 46.6333) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (12, 'Moldova', 'Edinet', 27.3167, 48.1667) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (13, 'Moldova', 'Drochia', 27.7500, 48.0333) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (14, 'Moldova', 'Ialoveni', 28.7833, 46.9500) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (15, 'Moldova', 'Hincesti', 28.5833, 46.8167) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (16, 'Moldova', 'Singerei', 28.1500, 47.6333) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (17, 'Moldova', 'Taraclia', 28.6689, 45.9000) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (18, 'Moldova', 'Falesti', 27.7139, 47.5722) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (19, 'Moldova', 'Floresti', 28.3014, 47.8933) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (20, 'Moldova', 'Cimislia', 28.7833, 46.5167) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (21, 'Moldova', 'Rezina', 28.9500, 47.7333) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (22, 'Moldova', 'Anenii Noi', 29.2167, 46.8833) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (23, 'Moldova', 'Calarasi', 28.3000, 47.2500) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (24, 'Moldova', 'Nisporeni', 28.1833, 47.0833) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (25, 'Moldova', 'Riscani', 27.5539, 47.9572) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (26, 'Moldova', 'Glodeni', 27.5167, 47.7667) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (27, 'Moldova', 'Basarabeasca', 28.9614, 46.3336) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (28, 'Moldova', 'Leova', 28.2500, 46.4833) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (29, 'Moldova', 'Briceni', 27.0839, 48.3611) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (30, 'Moldova', 'Ocnita', 27.4392, 48.3853) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (31, 'Moldova', 'Telenesti', 28.3667, 47.5028) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (32, 'Moldova', 'Donduseni', 27.5833, 48.2167) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (33, 'Moldova', 'Stefan Voda', 29.6631, 46.5153) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (34, 'Moldova', 'Criuleni', 29.1667, 47.2167) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (35, 'Moldova', 'Soldanesti', 28.8000, 47.8167) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (36, 'Moldova', 'Cantemir', 28.2167, 46.2667) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (37, 'Moldova', 'Cocieri', 29.1167, 47.3000) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (38, 'Moldova', 'Ribnita', 29.0000, 47.7667) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (39, 'Moldova', 'Dubasari', 29.1667, 47.2667) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (40, 'Moldova', 'Slobozia', 29.7000, 46.7333) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (41, 'Moldova', 'Durlesti', 28.9500, 47.0333) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (42, 'Moldova', 'Codru', 28.8194, 46.9753) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (43, 'Moldova', 'Ceadir-Lunga', 28.8333, 46.0500) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (44, 'Romania', 'Moldova Noua', 21.6639, 44.7178) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (45, 'Moldova', 'Vulcanesti', 28.4028, 45.6833) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (46, 'Moldova', 'Cricova', 28.8667, 47.1333) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (47, 'Moldova', 'Bacioi', 28.8839, 46.9122) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (48, 'Moldova', 'Congaz', 28.5972, 46.1083) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (49, 'Moldova', 'Truseni', 28.6833, 47.0667) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (50, 'Moldova', 'Costesti', 28.7689, 46.8678) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (51, 'Moldova', 'Dnestrovsc', 29.9167, 46.6167) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (52, 'Moldova', 'Singera', 28.9708, 46.9139) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (53, 'Moldova', 'Borogani', 28.5167, 46.3667) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (54, 'Moldova', 'Grigoriopol', 29.2925, 47.1503) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (55, 'Moldova', 'Stauceni', 28.8703, 47.0875) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (56, 'Moldova', 'Peresecina', 28.7667, 47.2500) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (57, 'Moldova', 'Copceac', 28.6944, 45.8500) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (58, 'Moldova', 'Chitcani', 29.6167, 46.7833) ON CONFLICT DO NOTHING;
INSERT INTO public.location_table (location_id, country, city, longitude, latitude) VALUES (59, 'Moldova', 'Camenca', 28.7167, 48.0167) ON CONFLICT DO NOTHING;


--
-- TOC entry 4867 (class 0 OID 0)
-- Dependencies: 219
-- Name: disaster_list_table_disaster_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.disaster_list_table_disaster_id_seq', 1, true);


--
-- TOC entry 4868 (class 0 OID 0)
-- Dependencies: 217
-- Name: location_table_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.location_table_location_id_seq', 60, false);


--
-- TOC entry 4869 (class 0 OID 0)
-- Dependencies: 215
-- Name: natural_disaster_type_table_disaster_type_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.natural_disaster_type_table_disaster_type_id_seq', 1, false);


--
-- TOC entry 4706 (class 2606 OID 16551)
-- Name: disaster_list_table disaster_list_table_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.disaster_list_table
    ADD CONSTRAINT disaster_list_table_pk PRIMARY KEY (disaster_id);


--
-- TOC entry 4704 (class 2606 OID 16527)
-- Name: location_table location_table_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.location_table
    ADD CONSTRAINT location_table_pk PRIMARY KEY (location_id);


--
-- TOC entry 4702 (class 2606 OID 16499)
-- Name: disaster_type_table natural_disaster_type_table_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.disaster_type_table
    ADD CONSTRAINT natural_disaster_type_table_pkey PRIMARY KEY (disaster_type_id);


--
-- TOC entry 4707 (class 2606 OID 16521)
-- Name: disaster_list_table disaster_list_table_disaster_type_table_disaster_type_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.disaster_list_table
    ADD CONSTRAINT disaster_list_table_disaster_type_table_disaster_type_id_fk FOREIGN KEY (disaster_type_id) REFERENCES public.disaster_type_table(disaster_type_id);


--
-- TOC entry 4708 (class 2606 OID 16528)
-- Name: disaster_list_table disaster_list_table_location_table_location_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.disaster_list_table
    ADD CONSTRAINT disaster_list_table_location_table_location_id_fk FOREIGN KEY (location_id) REFERENCES public.location_table(location_id);


-- Completed on 2023-10-24 23:08:16

--
-- PostgreSQL database dump complete
--

