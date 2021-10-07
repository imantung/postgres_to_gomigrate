DROP INDEX public.books_title_idx;
ALTER TABLE ONLY public.books DROP CONSTRAINT books_pkey;
ALTER TABLE public.books ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE public.books_id_seq;
DROP TABLE public.books;