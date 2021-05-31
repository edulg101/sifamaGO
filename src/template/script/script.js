
function startDigitacao1() {
    const requestData = {
        StartDigitacao: true,
        Restart: false,
    };

    const usersPromise = fetch('/', {
        method: 'POST',
        body: JSON.stringify(requestData)
    }).then(response => {
        if (!response.ok) {
            throw new Error("Got non-2XX response from API server.");
        }
        return response.json();
    });
}

function startDigitacao() {
    const requestData = {
        StartDigitacao: true,
        Restart: false,
    };

    const usersPromise = fetch('/report', {
        method: 'POST',
        body: JSON.stringify(requestData)
    });
}

function sendFirstForm() {

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





}

function toReport() {
    console.log('entrou')
    element = document.getElementById('folderselect').value

}

async function restart() {
    const requestData = {
        StartDigitacao: false,
        Restart: true,
    };

    const response = fetch('/report', {
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
