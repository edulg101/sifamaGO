
async function startDigitacao() {

    user = document.getElementById('user').value
    pwd = document.getElementById('password').value
    spinner = document.getElementById('spinner1')
    errorDiv = document.getElementById('notification')
    successDiv = document.getElementById('sucess')
    spinner.style.visibility = 'visible'

    const requestData = {
        StartDigitacao: true,
        Restart: false,
        User: user,
        Passd: pwd

    };

    if (user == "" || pwd == "") {
        spinner.style.visibility = 'hidden'

        let timer = document.createElement('div')
        timer.innerHTML = "Os campos devem ser preenchidos antes de enviar. Tente novamente... "

        while (div.firstChild) {
            div.removeChild(div.firstChild);
        }
        div.style.display = 'inline-block'

        div.appendChild(timer)
        setTimeout(() => {
            div.style.display = 'none'
        }, 2000);

    } else {

        document.getElementById('user').value = '';
        document.getElementById('password').value = '';

        const response = await fetch('/report', {
            method: 'POST',
            body: JSON.stringify(requestData)
        });

        let responseText = await response.text();
        let status = response.status;

        if (status != 200) {
            spinner.style.visibility = 'hidden'
            let newDiv = document.createElement('div')
            newDiv.innerHTML = responseText
            while (errorDiv.firstChild) {
                errorDiv.removeChild(errorDiv.firstChild);
            }
            errorDiv.style.display = 'inline-block'
            errorDiv.appendChild(newDiv)
        } else {
            spinner.style.visibility = 'hidden'
            let newDiv = document.createElement('div')
            newDiv.innerHTML = responseText
            while (successDiv.firstChild) {
                successDiv.removeChild(successDiv.firstChild);
            }
            successDiv.style.display = 'inline-block'
            successDiv.appendChild(newDiv)

        }
    }

}

async function getResponseFromServer(userId, pwdId, spinnerID, errorDivId, sucessDivId, API, requestData) {

    if (!userId == null) {
        console.log('entrou no null')
        user = document.getElementById(userId).value
        pwd = document.getElementById(pwdId).value

        document.getElementById(userId).value = '';
        document.getElementById(pwdId).value = '';
    }

    spinner = document.getElementById(spinnerID)
    successDiv = document.getElementById(sucessDivId)
    errorDiv = document.getElementById(errorDivId)




    spinner.style.display = 'inline-block'

    const response = await fetch(API, {
        method: 'POST',
        body: JSON.stringify(requestData)
    });

    let responseText = await response.text();
    let status = response.status;
    console.log(status)

    // return

    if (status != 200) {
        spinner.style.display = 'none'
        let newDiv = document.createElement('div')
        responseText = responseText.replace(/\n/g, '</br>')
        newDiv.innerHTML = responseText
        while (errorDiv.firstChild) {
            errorDiv.removeChild(errorDiv.firstChild);
        }
        errorDiv.style.display = 'inline-block';
        errorDiv.appendChild(newDiv)
        spinner.style.display = 'none'
        return false
    } else {
        let newDiv = document.createElement('div')
        spinner.style.display = 'none'
        successDiv.style.display = 'inline-block'
        newDiv.innerHTML = responseText
        while (successDiv.firstChild) {
            successDiv.removeChild(successDiv.firstChild);
        }
        successDiv.appendChild(newDiv)
        return true

    }

}


function formValidator(userId, pwdId, notificationId) {
    notificationDiv = document.getElementById(notificationId)
    user = document.getElementById(userId).value
    pwd = document.getElementById(pwdId).value

    if (user.length < 2 || pwd.length < 2) {
        let timer = document.createElement('div')
        timer.innerHTML = "Os campos devem ser preenchidos antes de enviar. Tente novamente... "
        while (notificationDiv.firstChild) {
            notificationDiv.removeChild(notificationDiv.firstChild);
        }
        notificationDiv.style.display = 'inline-block';
        notificationDiv.appendChild(timer)
        setTimeout(() => {
            notificationDiv.style.display = 'none';
        }, 2000);

        return false
    } else {
        return true
    }

}

async function checkSifama1() {

    console.log('checksifama1 launched')


    valid = formValidator('user', 'password', "notification1")

    if (!valid) {
        return
    }

    let requestData = {
        StartDigitacao: true,
        User: user,
        Passd: pwd
    };

    getResponseFromServer('user',
        'password',
        'spinnerCheck',
        'notification1',
        'sucess',
        '/checkSifama',
        requestData)

}

async function startDigitacao1() {


    userId = 'user'
    passwordId = 'password'
    notificationDivId = 'notification'
    spinnerId = 'spinner1'
    successId = 'sucess'
    API = '/report'

    valid = formValidator(userId, passwordId, notificationDivId)

    if (!valid) {
        return
    }

    const requestData = {
        StartDigitacao: true,
        Restart: false,
        User: user,
        Passd: pwd

    };

    getResponseFromServer(
        userId,
        passwordId,
        spinnerId,
        notificationDivId,
        successId,
        API,
        requestData)

}


// async function checkSifama() {

//     user = document.getElementById('user').value
//     pwd = document.getElementById('password').value
//     notificationDiv = document.getElementById('notification1')
//     notificationDiv.style.display = 'none';
//     spinner = document.getElementById('spinnerCheck')
//     successDiv = document.getElementById('sucess')


//     const requestData = {
//         StartDigitacao: true,
//         User: user,
//         Passd: pwd

//     };


//     if (user == "" || pwd == "") {

//         let timer = document.createElement('div')
//         timer.innerHTML = "Os campos devem ser preenchidos antes de enviar. Tente novamente... "
//         while (notificationDiv.firstChild) {
//             notificationDiv.removeChild(notificationDiv.firstChild);
//         }

//         notificationDiv.style.display = 'inline-block';

//         notificationDiv.appendChild(timer)
//         setTimeout(() => {
//             notificationDiv.style.display = 'none';
//         }, 2000);

//     } else {

//         document.getElementById('user').value = '';
//         document.getElementById('password').value = '';

//         // spinner = document.getElementById('spinnerCheck')
//         spinner.style.display = 'inline-block'

//         const response = await fetch('/checkSifama', {
//             method: 'POST',
//             body: JSON.stringify(requestData)
//         });

//         let responseText = await response.text();
//         let status = response.status;
//         console.log(status)

//         if (status != 200) {
//             spinner.style.display = 'none'
//             let newDiv = document.createElement('div')
//             responseText = responseText.replace(/\n/g, '</br>')
//             newDiv.innerHTML = responseText
//             while (notificationDiv.firstChild) {
//                 notificationDiv.removeChild(notificationDiv.firstChild);
//             }
//             notificationDiv.style.display = 'inline-block';
//             notificationDiv.appendChild(newDiv)
//         } else {
//             let newDiv = document.createElement('div')
//             spinner.style.display = 'none'
//             successDiv.style.display = 'inline-block'
//             newDiv.innerHTML = responseText
//             while (successDiv.firstChild) {
//                 successDiv.removeChild(successDiv.firstChild);
//             }
//             successDiv.appendChild(newDiv)

//         }
//         spinner.style.display = 'none'
//     }
// }

async function toReport1() {

    console.log('toReport1 launched')

    userId = null
    passwordId = null
    notificationDivId = 'notification'
    spinnerId = 'spinner'
    successId = 'sucess'
    API = '/'



    folder = document.getElementById('homeselect').value
    title = document.getElementById('hometitle').value

    document.getElementById('reduzir').disabled = true;
    document.getElementById('gerarInput').disabled = true;



    // spinner = document.getElementById('spinner')
    // spinner.style.visibility = 'visible'

    const requestData = {
        Folder: folder,
        Title: title,
    };



    sucess = await getResponseFromServer(
        userId,
        passwordId,
        spinnerId,
        notificationDivId,
        successId,
        API,
        requestData);

    document.getElementById('reduzir').disabled = false;
    document.getElementById('gerarInput').disabled = false;

    console.log(sucess)

    if (sucess) {
        window.location.href = '/report';
    }
}

async function toReport() {

    folder = document.getElementById('homeselect').value
    title = document.getElementById('hometitle').value

    document.getElementById('reduzir').disabled = true;
    document.getElementById('gerarInput').disabled = true;

    notificationDiv = document.getElementById('notification')
    notificationDiv.style.display = 'none'


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
        while (notificationDiv.firstChild) {
            notificationDiv.removeChild(notificationDiv.firstChild);
        }
        notificationDiv.style.display = 'inline-block'
        notificationDiv.appendChild(timer)
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

    notificationDiv = document.getElementById('notification')

    notificationDiv.style.display = 'none'
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
        while (notificationDiv.firstChild) {
            notificationDiv.removeChild(notificationDiv.firstChild);
        }
        notificationDiv.style.display = 'inline-block'
        notificationDiv.appendChild(timer)
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
    errorDiv = document.getElementById('notification')
    errorDiv.style.display = 'none'
    successDiv = document.getElementById('sucess1')
    spinner = document.getElementById('spinner')
    spinner.style.visibility = 'visible'

    const requestData = {
        Compact: true,
        Folder: folder,
    };

    const response = await fetch('/compact', {
        method: 'POST',
        body: JSON.stringify(requestData)
    });

    let responseText = await response.text();
    let status = response.status;

    if (status != 200) {
        spinner.style.visibility = 'hidden'
        let newDiv = document.createElement('div')
        newDiv.innerHTML = responseText
        while (errorDiv.firstChild) {
            errorDiv.removeChild(errorDiv.firstChild);
        }
        errorDiv.style.display = 'inline-block'
        errorDiv.appendChild(newDiv)
    } else {
        spinner.style.visibility = 'hidden'
        let newDiv = document.createElement('div')
        newDiv.innerHTML = responseText
        while (successDiv.firstChild) {
            successDiv.removeChild(successDiv.firstChild);
        }
        successDiv.style.display = 'inline-block'
        successDiv.appendChild(newDiv)

    }
}
