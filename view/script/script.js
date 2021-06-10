
async function startDigitacao() {

    user = document.getElementById('user').value
    pwd = document.getElementById('password').value



    const requestData = {
        StartDigitacao: true,
        Restart: false,
        User: user,
        Passd: pwd

    };

    if (user == "" || pwd == "") {

        let timer = document.createElement('div')
        timer.innerHTML = "Os campos devem ser preenchidos antes de enviar. Tente novamente... "
        div = document.getElementById('notification')
        while (div.firstChild) {
            div.removeChild(div.firstChild);
        }
        div.style.visibility = 'visible'

        div.appendChild(timer)
        setTimeout(() => {
            div.style.visibility = 'hidden'
        }, 2000);

    } else {

        document.getElementById('user').value = '';
        document.getElementById('password').value = '';

        const responseText = await fetch('/report', {
            method: 'POST',
            body: JSON.stringify(requestData)
        }).then(response => response.text());

        let timer = document.createElement('div')
        timer.innerHTML = responseText
        div = document.getElementById('notification')
        while (div.firstChild) {
            div.removeChild(div.firstChild);
        }
        div.style.visibility = 'visible'
        div.appendChild(timer)
    }

}

async function toReport() {

    folder = document.getElementById('homeselect').value
    title = document.getElementById('hometitle').value

    document.getElementById('reduzir').disabled = true;
    document.getElementById('gerarInput').disabled = true;

    spinner = document.getElementById('spinner')
    spinner.style.visibility = 'visible'

    const requestData = {
        Folder: folder,
        Title: title,
    };


    const response = await fetch('/', {
        method: 'POST',
        body: JSON.stringify(requestData)
    });

    let text = await response.text();
    let status = response.status;

    if (status != 200) {
        let timer = document.createElement('div')
        timer.innerHTML = text;
        div = document.getElementById('notification')
        while (div.firstChild) {
            div.removeChild(div.firstChild);
        }
        div.style.visibility = 'visible'
        div.appendChild(timer)
        spinner.style.visibility = 'hidden'
        document.getElementById('reduzir').disabled = false;
        document.getElementById('gerarInput').disabled = false;
    } else {
        window.location.href = '/report';
    }

}

function saveFile() {
    elements = document.getElementsByClassName('alert');
    var arr = [].slice.call(elements);

    for (i = 0; i < arr.length; ++i) {
        element = arr[i];
        console.log('elemento', i)
        element.removeAttribute('class')
    }
    elements = document.getElementsByClassName('badge');
    var arr = [].slice.call(elements);
    for (i = 0; i < arr.length; i++) {
        element = arr[i]
        element.style.visibility = 'hidden'
    }

    element = document.getElementById('buttonsDiv')
    element.style.visibility = 'hidden'

    elements = document.getElementsByClassName('geotag');
    var arr = [].slice.call(elements);
    for (i = 0; i < arr.length; i++) {
        element = arr[i]
        element.style.visibility = 'hidden'
    }

    elements = document.getElementsByClassName('link');
    var arr = [].slice.call(elements);
    for (i = 0; i < arr.length; i++) {
        element = arr[i]
        element.style.visibility = 'hidden'
    }

    document.getElementById('sifamaForm').style.visibility = "hidden";

}


async function restart() {

    document.getElementById('save-file-button').disabled = true;
    document.getElementById('redobutton').disabled = true;

    spinner = document.getElementById('spinner')
    spinner.style.visibility = 'visible'

    const requestData = {
        StartDigitacao: false,
        Restart: true,
    };

    const response = await fetch('/report', {
        method: 'POST',
        body: JSON.stringify(requestData)
    });

    let text = await response.text();
    let status = response.status;


    if (status != 200) {
        let timer = document.createElement('div')
        timer.innerHTML = text;
        div = document.getElementById('notification')
        while (div.firstChild) {
            div.removeChild(div.firstChild);
        }
        div.style.visibility = 'visible'
        div.appendChild(timer)
        spinner.style.visibility = 'hidden'
        document.getElementById('save-file-button').disabled = false;
        document.getElementById('redobutton').disabled = false;
    } else {
        window.location.reload()
    }
}

function compactImages() {
    document.getElementById('reduzir').disabled = true;
    document.getElementById('gerarInput').disabled = true;
    document.getElementById('spinner').style.visibility = "visible";

    compact()

}

async function compact() {
    folder = document.getElementById('homeselect').value
    console.log(folder)
    const requestData = {
        Compact: true,
        Folder: folder,
    };

    const response = fetch('/compact', {
        method: 'POST',
        body: JSON.stringify(requestData)
    });

    status = (await response).status

    if (await status == 200) {
        window.location.reload()
    } else {
        console.log('pagina nao carregou')
        console.log(status)
        throw new Error("não carregou")
    }
}
