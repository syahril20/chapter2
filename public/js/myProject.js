let blogs = []

function test() {
    const today = new Date();
    const monthSet = new Date().setMonth(today.getMonth() + 3);
    const date = today.toLocaleDateString();
    let data = {
        title: "Phyton Jadi Bahasa Pemrograman Terpopuler di Dunia, Ini Alasannya",
        start: date,
        end: monthSet,
        postAt: today,
        description: `PARBOABOA - Tiobe merilis peringkat 10 bahasa pemograman paling populer di dunia untuk Oktober 2021. Dalam laporan bertajuk "Tiobe Programming Community index", Phyton dinobatkan sebagai bahasa programming terpopuler di dunia saat ini.`,
        image: "../assets/content.png",
        node: true,
        next: true,
        react: true,
        typescript: true,
    }
    blogs.push(data)
    renderBlog()
}

function getBlog(event) {
    event.preventDefault()

    let title = document.getElementById("title").value
    let start = document.getElementById("start").value
    let end = document.getElementById("end").value
    let description = document.getElementById("description").value
    let image = document.getElementById("img").files
    let node = document.getElementById("nodeBox").checked
    let next = document.getElementById("nextBox").checked
    let react = document.getElementById("reactBox").checked
    let typescript = document.getElementById("typeBox").checked

    image = URL.createObjectURL(image[0]);

    if (title == "") {
        return alert("title Tidak Boleh Kosong")
    } else if (start == "") {
        return alert("Start Tidak Boleh Kosong")
    } else if (end == "") {
        return alert("End Tidak Boleh Kosong")
    } else if (start > end) {
        return alert("End Date Tidak Boleh Kurang dari Start Date")
    } else if (description == "") {
        return alert("description Tidak Boleh Kosong")
    }

    let blog = {
        title,
        start,
        end,
        postAt: new Date(),
        description,
        image,
        node,
        next,
        react,
        typescript,
    }

    blogs.push(blog)
    console.log(blogs)
    renderBlog()
}

function renderBlog() {
    document.getElementById("content").innerHTML = ""
    for (let i = 0; i < blogs.length; i++) {
        document.getElementById("content").innerHTML += `
        <div class="card-box">
                <div class="header">
                    <img src="${blogs[i].image}" alt="">
                    <a href="projectDetail.html">${blogs[i].title}</a>
                    <label>durasi : ${createDuration(blogs[i].start, blogs[i].end)}</label>
                    <label>${getPosted(blogs[i].postAt)}</label>
                </div>
                <div class="content">
                    <p>${blogs[i].description}</p>
                </div>
                <div class="footer">
                    <div class="logo-area">
                    ${blogs[i].node ? '<img src="../assets/node.png" alt="logo" class="tech-logo">' : ''}
                    ${blogs[i].next ? '<img src="../assets/next.png" alt="logo" class="tech-logo">' : ''}
                    ${blogs[i].react ? '<img src="../assets/react.png" alt="logo" class="tech-logo">' : ''}
                    ${blogs[i].typescript ? '<img src="../assets/typescript.png" alt="logo" class="tech-logo">' : ''}
                    </div>
                    <div class="button-area">
                        <div>
                            <a href="" class="foot-button">edit</a>
                        </div>
                        <div>
                            <a href="" class="foot-button">delete</a>
                        </div>
                    </div>
                </div>
            </div>`
    }
}

function createDuration(start, end) {
    let distance = new Date(end) - new Date(start)
    // Valid
    const dayDistance = Math.floor(distance / (24 * 60 * 60 * 1000))
    const monthDistance = Math.floor(distance / (30 * 24 * 60 * 60 * 1000))
    const yearDistance = Math.floor(distance / (12 * 30 * 24 * 60 * 60 * 1000))

    if (dayDistance > 0 && dayDistance <= 29) {
        return `${dayDistance} Hari `
    } else if (dayDistance >= 30 && monthDistance <= 12) {
        return monthDistance + " Bulan"
    } else if (monthDistance >= 12) {
        return yearDistance + " Tahun"
    } else if (dayDistance >= 0 && dayDistance <= 24) {
        return "1 hari"
    }
}

function getPosted(time) {
    const posted = new Date() - new Date(time)
    const dayPosted = Math.floor(posted / (24 * 60 * 60 * 1000))
    if (dayPosted > 0) {
        return dayPosted + " day ago"
    } else {
        const hourPosted = Math.floor(posted / (60 * 60 * 1000))
        if (hourPosted > 0) {
            return hourPosted + " hour ago"
        } else {
            const minutesPosted = Math.floor(posted / (60 * 1000))
            if (minutesPosted > 0) {
                return minutesPosted + " minutes ago"
            } else {
                const secondPosted = Math.floor(posted / 1000)
                if (secondPosted > 0) {
                    return secondPosted + " second ago"
                }
            }
        }
    }
}

setInterval(renderBlog, 1000)