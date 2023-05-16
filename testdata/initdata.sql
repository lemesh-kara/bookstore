CREATE EXTENSION fuzzystrmatch;

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    author TEXT NOT NULL,
    year TEXT NOT NULL,
    description TEXT NOT NULL,
    path_to_pdf TEXT NOT NULL,
    review_mark DOUBLE PRECISION NOT NULL CHECK (review_mark >= 0 AND review_mark <= 5),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

INSERT INTO books (name, author, year, description, path_to_pdf, review_mark)
VALUES
    ('The Great Gatsby', 'F. Scott Fitzgerald', '1925', 'A story of love, greed, and tragedy in the Jazz Age.', '/books/the-great-gatsby.pdf', 4.5),
    ('To Kill a Mockingbird', 'Harper Lee', '1960', 'A powerful story of racial injustice and coming of age in the Deep South.', '/books/to-kill-a-mockingbird.pdf', 4.8),
    ('1984', 'George Orwell', '1949', 'A dystopian novel about government surveillance, censorship, and oppression.', '/books/1984.pdf', 3.9),
    ('Pride and Prejudice', 'Jane Austen', '1813', 'A witty and romantic novel about social class and marriage in 19th century England.', '/books/pride-and-prejudice.pdf', 4.2),
    ('Moby-Dick', 'Herman Melville', '1851', 'An epic tale of obsession and revenge on the high seas.', '/books/moby-dick.pdf', 3.5);

