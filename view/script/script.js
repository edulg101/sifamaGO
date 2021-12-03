window.onload = changeMapHeight();

async function getResponseFromServer(
    userId,
    pwdId,
    spinnerID,
    errorDivId,
    successDivId,
    API,
    requestData
) {
    if (userId.length > 1) {
        console.log('entrou no null');
        clearForm([userId, pwdId]);
    }

    spinner = document.getElementById(spinnerID);
    successDiv = document.getElementById(successDivId);
    errorDiv = document.getElementById(errorDivId);
    errorDiv.style.display = 'none';
    successDiv.style.display = 'none';

    spinner.style.display = 'inline-block';

    const response = await fetch(API, {
        method: 'POST',
        body: JSON.stringify(requestData),
    });

    let responseText = await response.text();
    let status = response.status;

    console.log(status);
    console.log(responseText);

    if (status != 200) {
        spinner.style.display = 'none';
        const newDiv = document.createElement('div');
        responseText = responseText.replace(/\n/g, '</br>');
        newDiv.innerHTML = responseText;
        while (errorDiv.firstChild) {
            errorDiv.removeChild(errorDiv.firstChild);
        }
        errorDiv.style.display = 'inline-block';
        errorDiv.appendChild(newDiv);
        spinner.style.display = 'none';
        return false;
    } else {
        let newDiv = document.createElement('div');
        spinner.style.display = 'none';
        successDiv.style.display = 'inline-block';
        newDiv.innerHTML = responseText;
        while (successDiv.firstChild) {
            successDiv.removeChild(successDiv.firstChild);
        }
        successDiv.appendChild(newDiv);
        return true;
    }
}

function formValidator(userId, pwdId, notificationId) {
    notificationDiv = document.getElementById(notificationId);
    const user = document.getElementById(userId).value;
    const pwd = document.getElementById(pwdId).value;

    if (user.length < 2 || pwd.length < 2) {
        let timer = document.createElement('div');
        timer.innerHTML =
            'Os campos devem ser preenchidos antes de enviar. Tente novamente... ';
        while (notificationDiv.firstChild) {
            notificationDiv.removeChild(notificationDiv.firstChild);
        }
        notificationDiv.style.display = 'inline-block';
        notificationDiv.appendChild(timer);
        setTimeout(() => {
            notificationDiv.style.display = 'none';
        }, 2000);

        return false;
    } else {
        return true;
    }
}

function clearForm(fieldsId) {
    for (x in fieldsId) {
        document.getElementById(fieldsId[x]).value = '';
    }
}

async function checkSifama1() {
    console.log('entrou checksifama1');
    const userId = 'user';
    const pwdId = 'password';

    const user = document.getElementById(userId).value;
    const pwd = document.getElementById(pwdId).value;

    const valid = formValidator('user', 'password', 'notification1');

    if (!valid) {
        return;
    }

    const requestData = {
        StartDigitacao: true,
        User: user,
        Passd: pwd,
    };

    getResponseFromServer(
        userId,
        pwdId,
        'spinnerCheck',
        'notification1',
        'success',
        '/checkSifama',
        requestData
    );
}

async function startDigitacao1() {
    userId = 'user';
    passwordId = 'password';
    notificationDivId = 'notification';
    spinnerId = 'spinner1';
    successId = 'success';
    API = '/report';

    // document.getElementById(spinnerId).style.display = 'inline-block'

    user = document.getElementById(userId).value;
    pwd = document.getElementById(passwordId).value;

    valid = formValidator(userId, passwordId, notificationDivId);

    if (!valid) {
        return;
    }

    const requestData = {
        StartDigitacao: true,
        Restart: false,
        User: user,
        Passd: pwd,
    };

    getResponseFromServer(
        userId,
        passwordId,
        spinnerId,
        notificationDivId,
        successId,
        API,
        requestData
    );
}

async function toReport1() {
    userId = '';
    passwordId = '';
    notificationDivId = 'notification';
    spinnerId = 'spinner';
    successId = 'success';
    API = '/';

    folder = document.getElementById('homeselect').value;
    title = document.getElementById('hometitle').value;

    document.getElementById('reduzir').disabled = true;
    document.getElementById('gerarInput').disabled = true;

    const requestData = {
        Folder: folder,
        Title: title,
    };

    success = await getResponseFromServer(
        userId,
        passwordId,
        spinnerId,
        notificationDivId,
        successId,
        API,
        requestData
    );

    document.getElementById('reduzir').disabled = false;
    document.getElementById('gerarInput').disabled = false;

    console.log(success);

    if (success) {
        window.location.href = '/report';
    }
}

function saveFile() {
    elements = document.getElementsByClassName('alert');
    var arr = [].slice.call(elements);

    for (i = 0; i < arr.length; ++i) {
        element = arr[i];
        console.log('elemento', i);
        element.removeAttribute('class');
    }
    elements = document.getElementsByClassName('badge');
    var arr = [].slice.call(elements);
    for (i = 0; i < arr.length; i++) {
        element = arr[i];
        element.style.visibility = 'hidden';
    }

    element = document.getElementById('buttonsDiv');
    element.style.visibility = 'hidden';

    // elements = document.getElementsByClassName('geotag');
    // var arr = [].slice.call(elements);
    // for (i = 0; i < arr.length; i++) {
    //     element = arr[i]
    //     element.style.visibility = 'hidden'
    // }

    // elements = document.getElementsByClassName('link');
    // var arr = [].slice.call(elements);
    // for (i = 0; i < arr.length; i++) {
    //     element = arr[i]
    //     element.style.visibility = 'hidden'
    // }

    document.getElementById('sifamaForm').style.visibility = 'hidden';
}

async function restart() {
    document.getElementById('save-file-button').disabled = true;
    document.getElementById('redobutton').disabled = true;

    notificationDiv = document.getElementById('notification');

    notificationDiv.style.display = 'none';
    spinner = document.getElementById('spinner');
    spinner.style.visibility = 'visible';

    const requestData = {
        StartDigitacao: false,
        Restart: true,
    };

    const response = await fetch('/report', {
        method: 'POST',
        body: JSON.stringify(requestData),
    });

    let text = await response.text();
    let status = response.status;

    if (status != 200) {
        let timer = document.createElement('div');
        timer.innerHTML = text;
        while (notificationDiv.firstChild) {
            notificationDiv.removeChild(notificationDiv.firstChild);
        }
        notificationDiv.style.display = 'inline-block';
        notificationDiv.appendChild(timer);
        spinner.style.visibility = 'hidden';

        document.getElementById('save-file-button').disabled = false;
        document.getElementById('redobutton').disabled = false;
    } else {
        window.location.reload();
    }
}

function compactImages() {
    document.getElementById('reduzir').disabled = true;
    document.getElementById('gerarInput').disabled = true;
    // document.getElementById('spinner').style.visibility = "visible";
    compact1();
    document.getElementById('reduzir').disabled = false;
    document.getElementById('gerarInput').disabled = false;
}

async function compact1() {
    const folder = document.getElementById('homeselect').value;
    const errorDivId = 'notification';
    const successId = 'success1';
    const spinnerId = 'spinner';
    const API = '/compact';

    const userId = '';
    const passwordId = '';

    const requestData = {
        Compact: true,
        Folder: folder,
    };

    getResponseFromServer(
        userId,
        passwordId,
        spinnerId,
        errorDivId,
        successId,
        API,
        requestData
    );
}

async function compact() {
    folder = document.getElementById('homeselect').value;
    errorDiv = document.getElementById('notification');
    errorDiv.style.display = 'none';
    successDiv = document.getElementById('success1');
    spinner = document.getElementById('spinner');
    spinner.style.visibility = 'visible';
    const requestData = {
        Compact: true,
        Folder: folder,
    };

    const response = await fetch('/compact', {
        method: 'POST',
        body: JSON.stringify(requestData),
    });

    let responseText = await response.text();
    let status = response.status;
    console.log(status);
    console.log(responseText);

    if (status != 200) {
        spinner.style.visibility = 'hidden';
        let newDiv = document.createElement('div');
        newDiv.innerHTML = responseText;
        while (errorDiv.firstChild) {
            errorDiv.removeChild(errorDiv.firstChild);
        }
        errorDiv.style.display = 'inline-block';
        errorDiv.appendChild(newDiv);
    } else {
        spinner.style.visibility = 'hidden';
        let newDiv = document.createElement('div');
        newDiv.innerHTML = responseText;
        while (successDiv.firstChild) {
            successDiv.removeChild(successDiv.firstChild);
        }
        successDiv.style.display = 'inline-block';
        successDiv.appendChild(newDiv);
    }
}

function changeMapHeight() {
    height = window.innerHeight;
    div = document.getElementById('map');
    div.style.height = height + 'px';
    console.log(height);
}

function onClick(element) {
    console.log('clicou na imagem');
    document.getElementById('img01').src = element.src;
    document.getElementById('img01').style.display = 'block';
    document.getElementById('modal01').style.display = 'block';

    document.getElementById('modal01').style.verticalAlign = 'middle';
    document.getElementById('modal01').style.padding = 'auto';

    document.getElementById('img01').style.verticalAlign = 'middle';
    // document.getElementById('img01').style.marginRight = 'auto';
    // document.getElementById('img01').style.marginLeft = 'auto';
}