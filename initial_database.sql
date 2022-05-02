CREATE TABLE books (
    id serial PRIMARY KEY,
    title VARCHAR (255) NOT NULL,
    author VARCHAR (255) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

create index books_title_idx on books(title);

insert into books (title, author) values ('title1', 'author1');
insert into books (title, author) values ('title2', 'author2');
insert into books (title, author) values ('title3', 'author3');