function getData() {
    let name = document.getElementById("name").value
    let email = document.getElementById("email").value
    let phone = document.getElementById("phone").value
    let subject = document.getElementById("subject").value
    let message = document.getElementById("message").value

    if (name == "") {
        return alert("Nama Tidak Boleh Kosong")
    } else if (email == "") {
        return alert("Email Tidak Boleh Kosong")
    } else if (phone == "") {
        return alert("No Tidak Boleh Kosong")
    } else if (subject == "") {
        return alert("Subject Tidak Boleh Kosong")
    } else if (message == "") {
        return alert("Pesan Tidak Boleh Kosong")
    }

    const destination = "rodyulo20@gmail.com"
    let a = document.createElement("a")
    a.href = `mailto:${destination}?subject=${subject}&body=Hallo nama saya ${name} , saya ingin ${message}, bisakah anda menghubungi saya di ${phone}`
    a.click()

    let data = {
        name,
        email,
        phone,
        subject,
        message,
    }

    console.log(data)
}




