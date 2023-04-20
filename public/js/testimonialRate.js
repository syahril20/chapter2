// const testiData = [
//     {
//         image: "./assets/testimonial/patrick.jpg",
//         quotes: "Berhentilah hidup di masa lalu, itu hanya akan menyakitimu",
//         author: "Patrick Star",
//         rating: 5
//     },
//     {
//         image: "./assets/testimonial/squidward.jpg",
//         quotes: "Orang banyak bicara ketika mereka ingin terlihat pintar",
//         author: "Squidward Tentacles",
//         rating: 4
//     },
//     {
//         image: "./assets/testimonial/mrCrab.jpg",
//         quotes: "Uang lebih manis daripada madu",
//         author: "Eugene Crab",
//         rating: 4
//     },
//     {
//         image: "./assets/testimonial/patrick.jpg",
//         quotes: "Pengetahuan tidak dapat menggantikan persahabatan",
//         author: "Patrick Star",
//         rating: 1
//     },
//     {
//         image: "./assets/testimonial/squidward.jpg",
//         quotes: "Aku benci kalian semua",
//         author: "Squidward Tentacles",
//         rating: 2
//     }
// ]
const testimonial = new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()

    xhr.open("GET", "https://api.npoint.io/d77c61b8e4f36801ff8a", true)

    xhr.onload = function (){
        if(xhr.status == 200) {
            resolve(JSON.parse(xhr.response))
        } else {
            reject("Error loading data")
        }
    }

    xhr.onerror = function() {
        reject("Network Error")
    }

    xhr.send()
})


async function renderData() {
    const response = await testimonial
    console.log(response)
    let testiHtml = ''

    response.forEach((data) => {
        testiHtml += `
<div class="card shadow" id="testi">
    <div class="testiContent">
        <div>
            <img src="${data.image}" alt="" class="testiImg">
        </div>
        <div>
            <p class="quotes">"${data.quotes}"</p>
        </div>
        <div>
            <p class="author">- ${data.author}</p>
        </div>
        <div>
            <p class="rating">${data.rating}<img src="./assets/black-star.png" class="star"></p>
        </div>
    </div>
</div>
`
    })

    document.getElementById("testBody").innerHTML = testiHtml
}
renderData()

async function filterRating(rating) {
    const response =  await testimonial
    let testiHtml = ''

    const filterData = response.filter(function (data) {
        return data.rating === rating
    })
    console.log(filterData)

    if (filterData.length === 0) {
        testiHtml = `<h1>Data Not Found</h1>`
    } else {
        filterData.forEach((data) => {
            testiHtml += `
<div class="card shadow" id="testi">
    <div class="testiContent">
        <div>
            <img src="${data.image}" alt="" class="testiImg">
        </div>
        <div>
            <p class="quotes">"${data.quotes}"</p>
        </div>
        <div>
            <p class="author">- ${data.author}</p>
        </div>
        <div>
            <p class="author">${data.rating}<img src="./assets/black-star.png" class="star"></p>
        </div>
    </div>
</div>
        `
        })
    }

    document.getElementById("testBody").innerHTML = testiHtml
}