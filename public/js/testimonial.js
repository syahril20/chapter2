class Blog {
    constructor(image, quotes, author) {
        this._image = image;
        this._quotes = quotes;
        this._author = author;
    }

    // getter dan setter untuk properti image
    get image() {
        return this._image;
    }

    set image(image) {
        this._image = image;
    }

    // getter dan setter untuk properti quotes
    get quotes() {
        return this._quotes;
    }

    set quotes(quotes) {
        this._quotes = quotes;
    }

    // getter dan setter untuk properti author
    get author() {
        return this._author;
    }

    set author(author) {
        this._author = author;
    }

    displayBlog() {
        return `
        <div class="card shadow" id="testi">
            <div class="testiContent">
                <div>
                    <img src="${this._image}" alt="" class="testiImg">
                </div>
                <div>
                    <p class="quotes">${this._quotes}</p>
                </div>
                <div>
                    <p class="author">${this._author}</p>
                </div>
            </div>
        </div>
            `;
    }
}

const test1 = new Blog(
    "./assets/testimonial/patrick.jpg",
    `"Berhentilah hidup di masa lalu, itu hanya akan menyakitimu"`,
    "- Patrick Star"
)
const test2 = new Blog(
    "./assets/testimonial/squidward.jpg",
    `"Orang banyak bicara ketika mereka ingin terlihat pintar"`,
    "- Squidward Tentacles"
)
const test3 = new Blog(
    "./assets/testimonial/mrCrab.jpg",
    `"Uang lebih manis daripada madu"`,
    "- Eugene Crab"
)

let blogs = [test1, test2, test3];

let contentDiv = document.getElementById("testBody");

for (let i = 0; i < blogs.length; i++) {
    contentDiv.innerHTML += blogs[i].displayBlog();
}
