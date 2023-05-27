package testdata

import (
	"bookstore/model"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func InitDatabase(db *gorm.DB) {
	if !isDatabaseEmpty(db) {
		return
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;")

	createTables(db)
}

func CreateTestData(db *gorm.DB) {
	if !isDatabaseEmpty(db) {
		return
	}

	books := createBooks()
	users := createUsers()
	reviews := createReviewsAndUpdateBooks(books, users)
	carts := createCarts(books, users)
	news := createNews()

	saveBooks(db, books)
	saveUsers(db, users)
	saveReviews(db, reviews)
	saveCarts(db, carts)
	saveNews(db, news)
}

func isDatabaseEmpty(db *gorm.DB) bool {
	var count int64
	db.LogMode(false).Model(&model.Book{}).Count(&count)
	return count == 0
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func createTables(db *gorm.DB) {
	db.AutoMigrate(&model.Book{}, &model.User{}, &model.Review{}, &model.Cart{}, &model.News{})

	db.Model(&model.Review{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&model.Review{}).AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")

	db.Model(&model.Cart{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&model.Cart{}).AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")
}

func createBooks() []model.Book {
	rand.Seed(time.Now().UnixNano())
	usedISBNs := make(map[string]bool)
	generateISBN := func() string {
		for {
			isbn := fmt.Sprintf("978%d", rand.Intn(10000000000))
			checksum := 0
			for i, c := range isbn {
				digit := int(c - '0')
				if i%2 == 0 {
					checksum += digit * 1
				} else {
					checksum += digit * 3
				}
			}
			checksum = (10 - (checksum % 10)) % 10
			isbn += strconv.Itoa(checksum)
			if _, ok := usedISBNs[isbn]; !ok {
				return isbn
			}
		}
	}

	generatePublishedData := func() string {
		return strconv.FormatInt(int64(2000+rand.Intn(20)), 10) + " Published by Awesome Store #" + strconv.FormatInt(int64(rand.Intn(6)), 10)
	}

	return []model.Book{
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "The Great Gatsby", Author: "F. Scott Fitzgerald", Year: "1925", PublisherData: generatePublishedData(), Description: "A story of love, greed, and tragedy in the Jazz Age.", PathToPdf: "/pdf/the-great-gatsby.pdf", PathToCover: "/cover/the-great-gatsby.jpg", ReviewMark: 4.5, Price: 1340}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "To Kill a Mockingbird", Author: "Harper Lee", Year: "1960", PublisherData: generatePublishedData(), Description: "A powerful story of racial injustice and coming of age in the Deep South.", PathToPdf: "/pdf/to-kill-a-mockingbird.pdf", PathToCover: "/cover/to-kill-a-mockingbird.jpg", ReviewMark: 4.8, Price: 540}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "1984", Author: "George Orwell", Year: "1949", PublisherData: generatePublishedData(), Description: "A dystopian novel about government surveillance, censorship, and oppression.", PathToPdf: "/pdf/1984.pdf", PathToCover: "/cover/1984.jpg", ReviewMark: 3.9, Price: 1000}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "Pride and Prejudice", Author: "Jane Austen", Year: "1813", PublisherData: generatePublishedData(), Description: "A witty and romantic novel about social class and marriage in 19th century England.", PathToPdf: "/pdf/pride-and-prejudice.pdf", PathToCover: "/cover/pride-and-prejudice.jpg", ReviewMark: 4.2, Price: 2399}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "Moby-Dick", Author: "Herman Melville", Year: "1851", PublisherData: generatePublishedData(), Description: "An epic tale of obsession and revenge on the high seas.", PathToPdf: "/pdf/moby-dick.pdf", PathToCover: "/cover/moby-dick.jpg", ReviewMark: 3.5, Price: 3153}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "The Catcher in the Rye", Author: "J.D. Salinger", Year: "1951", PublisherData: generatePublishedData(), Description: "A controversial novel about teenage angst and rebellion.", PathToPdf: "/pdf/the-catcher-in-the-rye.pdf", PathToCover: "/cover/the-catcher-in-the-rye.jpg", ReviewMark: 3.7, Price: 10100}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "The Grapes of Wrath", Author: "John Steinbeck", Year: "1939", PublisherData: generatePublishedData(), Description: "A classic novel about the struggles of a family during the Great Depression.", PathToPdf: "/pdf/the-grapes-of-wrath.pdf", PathToCover: "/cover/the-grapes-of-wrath.jpg", ReviewMark: 4.1, Price: 1440}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "The Lord of the Rings", Author: "J.R.R. Tolkien", Year: "1954", PublisherData: generatePublishedData(), Description: "An epic fantasy trilogy about a quest to destroy a powerful ring.", PathToPdf: "/pdf/the-lord-of-the-rings.pdf", PathToCover: "/cover/the-lord-of-the-rings.jpg", ReviewMark: 4.9, Price: 1532}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "The Hobbit", Author: "J.R.R. Tolkien", Year: "1937", PublisherData: generatePublishedData(), Description: "A children's novel about a hobbit who goes on an adventure with dwarves and a wizard.", PathToPdf: "/pdf/the-hobbit.pdf", PathToCover: "/cover/the-hobbit.jpg", ReviewMark: 4.5, Price: 1040}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "The Picture of Dorian Gray", Author: "Oscar Wilde", Year: "1890", PublisherData: generatePublishedData(), Description: "A novel about a man who sells his soul for eternal youth and beauty.", PathToPdf: "/pdf/the-picture-of-dorian-gray.pdf", PathToCover: "/cover/the-picture-of-dorian-gray.jpg", ReviewMark: 4.0, Price: 940}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "One Hundred Years of Solitude", Author: "Gabriel Garcia Marquez", Year: "1967", PublisherData: generatePublishedData(), Description: "A novel about the Buendia family and the history of the fictional town of Macondo.", PathToPdf: "/pdf/one-hundred-years-of-solitude.pdf", PathToCover: "/cover/one-hundred-years-of-solitude.jpg", ReviewMark: 4.6, Price: 2540}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "Brave New World", Author: "Aldous Huxley", Year: "1932", PublisherData: generatePublishedData(), Description: "A dystopian novel about a future society that is controlled by technology and drugs.", PathToPdf: "/pdf/brave-new-world.pdf", PathToCover: "/cover/brave-new-world.jpg", ReviewMark: 4.0, Price: 3326}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "Frankenstein", Author: "Mary Shelley", Year: "1818", PublisherData: generatePublishedData(), Description: "A gothic novel about a scientist who creates a monster.", PathToPdf: "/pdf/frankenstein.pdf", PathToCover: "/cover/frankenstein.jpg", ReviewMark: 3.8, Price: 5539}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "The Adventures of Huckleberry Finn", Author: "Mark Twain", Year: "1884", PublisherData: generatePublishedData(), Description: "A novel about the adventures of a boy and a runaway slave on the Mississippi River.", PathToPdf: "/pdf/the-adventures-of-huckleberry-finn.pdf", PathToCover: "/cover/the-adventures-of-huckleberry-finn.jpg", ReviewMark: 4.2, Price: 7783}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "The Adventures of Tom Sawyer", Author: "Mark Twain", Year: "1876", PublisherData: generatePublishedData(), Description: "A novel about a mischievous boy growing up in a small town on the Mississippi River.", PathToPdf: "/pdf/the-adventures-of-tom-sawyer.pdf", PathToCover: "/cover/the-adventures-of-tom-sawyer.jpg", ReviewMark: 3.9, Price: 7920}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "Jane Eyre", Author: "Charlotte Bronte", Year: "1847", PublisherData: generatePublishedData(), Description: "A novel about a young governess who falls in love with her employer, Mr. Rochester.", PathToPdf: "/pdf/jane-eyre.pdf", PathToCover: "/cover/jane-eyre.jpg", ReviewMark: 4.1, Price: 1898}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "Wuthering Heights", Author: "Emily Bronte", Year: "1847", PublisherData: generatePublishedData(), Description: "A gothic novel about the doomed love between Heathcliff and Catherine.", PathToPdf: "/pdf/wuthering-heights.pdf", PathToCover: "/cover/wuthering-heights.jpg", ReviewMark: 3.7, Price: 999}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "The Count of Monte Cristo", Author: "Alexandre Dumas", Year: "1844", PublisherData: generatePublishedData(), Description: "A historical novel about a man who is wrongfully imprisoned and seeks revenge against those who betrayed him.", PathToPdf: "/pdf/the-count-of-monte-cristo.pdf", PathToCover: "/cover/the-count-of-monte-cristo.jpg", ReviewMark: 4.4, Price: 999}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "The Three Musketeers", Author: "Alexandre Dumas", Year: "1844", PublisherData: generatePublishedData(), Description: "An adventure novel about a young man who becomes friends with three musketeers and gets caught up in their escapades.", PathToPdf: "/pdf/the-three-musketeers.pdf", PathToCover: "/cover/the-three-musketeers.jpg", ReviewMark: 4.0, Price: 899}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "Gone with the Wind", Author: "Margaret Mitchell", Year: "1936", PublisherData: generatePublishedData(), Description: "A novel about a Southern belle named Scarlett O'Hara and her struggles during the Civil War and Reconstruction.", PathToPdf: "/pdf/gone-with-the-wind.pdf", PathToCover: "/cover/gone-with-the-wind.jpg", ReviewMark: 4.3, Price: 100}},
		{BookShort: model.BookShort{ISBN: generateISBN(), Name: "The Sun Also Rises", Author: "Ernest Hemingway", Year: "1926", PublisherData: generatePublishedData(), Description: "A novel about a group of expatriates who travel to Pamplona, Spain to watch the running of the bulls.", PathToPdf: "/pdf/the-sun-also-rises.pdf", PathToCover: "/cover/the-sun-also-rises.jpg", ReviewMark: 3.6, Price: 99}},
	}
}

func createUsers() []model.User {
	return []model.User{
		{UserShort: model.UserShort{Username: "admin", Password: hashPassword("123"), Email: "admin@bookstore.com", Role: "admin"}},
		{UserShort: model.UserShort{Username: "user1", Password: hashPassword("111"), Email: "user1@bookstore.com", Role: "admin"}},
		{UserShort: model.UserShort{Username: "user2", Password: hashPassword("222"), Email: "user2@bookstore.com", Role: "user"}},
		{UserShort: model.UserShort{Username: "user3", Password: hashPassword("333"), Email: "user3@bookstore.com", Role: "user"}},
		{UserShort: model.UserShort{Username: "user4", Password: hashPassword("444"), Email: "user4@bookstore.com", Role: "user"}},
	}
}

func createReviewsAndUpdateBooks(books []model.Book, users []model.User) (reviews []model.Review) {
	rand.Seed(time.Now().UnixNano())

	for index := range books {
		mark := float64(rand.Intn(6))
		reviews = append(reviews, model.Review{
			ReviewShort: model.ReviewShort{
				UserID:     1,
				Text:       "Great book to reader, index " + strconv.Itoa(index+1),
				BookID:     uint64(index + 1),
				ReviewMark: mark,
			},
		})
		books[index].ReviewMark = mark
	}

	amountOfMultReview := 3
	for i := 0; i < amountOfMultReview; i++ {
		reviewAmount := float64(1)
		totalMark := float64(books[i].ReviewMark)
		for index := range users {
			mark := float64(rand.Intn(6))
			reviews = append(reviews, model.Review{
				ReviewShort: model.ReviewShort{
					UserID:     uint64(index + 1),
					Text:       "Additional review, index" + strconv.Itoa(index),
					BookID:     uint64(i + 1),
					ReviewMark: mark,
				},
			})
			reviewAmount++
			totalMark += mark
		}
		books[i].ReviewMark = totalMark / reviewAmount
	}

	return
}

func createCarts(books []model.Book, users []model.User) (carts []model.Cart) {
	for index := range users {
		for i := 0; i <= index+1; i++ {
			carts = append(carts, model.Cart{
				CartShort: model.CartShort{
					UserID: uint64(index + 1),
					BookID: uint64(i + 1),
				},
			})
		}
	}

	return
}

func createNews() (news []model.News) {
	for i := 0; i < 10; i++ {
		news = append(news, model.News{
			NewsShort: model.NewsShort{
				Title:         "News title #" + strconv.Itoa(i),
				Text:          strconv.Itoa(i) + " We added a new item here about how we are doing this bookstore and so on, please read the thing, this will really help us to build the store",
				PathToPicture: "/newspic/pic" + strconv.Itoa(i) + ".jpg",
			},
		})
	}

	return
}

func saveBooks(db *gorm.DB, books []model.Book) {
	isTableEmpty := func() bool {
		var book model.Book
		return db.First(&book).Error != nil
	}

	if !isTableEmpty() {
		return
	}

	for _, book := range books {
		db.Create(&book)
	}
}

func saveUsers(db *gorm.DB, users []model.User) {
	isTableEmpty := func() bool {
		var user model.User
		return db.First(&user).Error != nil
	}

	if !isTableEmpty() {
		return
	}

	for _, user := range users {
		db.Create(&user)
	}
}

func saveReviews(db *gorm.DB, reviews []model.Review) {
	isTableEmpty := func() bool {
		var review model.Review
		return db.First(&review).Error != nil
	}

	if !isTableEmpty() {
		return
	}

	for _, review := range reviews {
		db.Create(&review)
	}
}

func saveCarts(db *gorm.DB, carts []model.Cart) {
	isTableEmpty := func() bool {
		var cart model.Cart
		return db.First(&cart).Error != nil
	}

	if !isTableEmpty() {
		return
	}

	for _, cart := range carts {
		db.Create(&cart)
	}
}

func saveNews(db *gorm.DB, news []model.News) {
	isTableEmpty := func() bool {
		var news model.News
		return db.First(&news).Error != nil
	}

	if !isTableEmpty() {
		return
	}

	for _, new := range news {
		db.Create(&new)
	}
}
