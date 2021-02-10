// Navigates to a page.
const nav = function navigateTo(page) {
    window.location.assign((page === '') ? '/' : '/' + page);
}

// Submits username and password and attempts to login.
const submit = function handleSubmitAndLogin() {
    let userJson = JSON.stringify({
        'username' : document.getElementById('email').value,
        'password' : document.getElementById('password').value,}
        );

    return fetch('https://localhost:443/token/', {
        method: 'POST',
        headers: {
            'Content-Type' : 'application/json',
            'X-CSRF-Token' : document.getElementsByTagName('meta')['csrf-token'].content,
        },
        body: userJson,
    }).then(response => {
        if (response.ok) {
            console.log('submit passed');
            nav('dashboard');
        } else {
            console.log('submit failed');
        }
        return response;
    }).catch(err => {
        console.log(err);
        return null;
    })
}

// Upgrades user plan to Enterprise plan.
const upgrade = function upgradeEnterprise() {
    return fetch('https://localhost:443/account/', {
        method: 'PATCH',
        headers : {
            'Content-Type' : 'application/json',
            'X-CSRF-Token' : document.getElementsByTagName('meta')['csrf-token'].content,
        },
        body: JSON.stringify({
            'plan' : 1,
        }),
    }).then(response => {
        if (response.ok) {
            console.log('upgraded');
        } else {
            console.log('upgrade failed');
            nav('');
        }
        return response;
    }).catch(err => {
        console.log(err);
        return null;
    })
}

// Logs user out from the dashboard.
const logout = function () {
    return fetch('https://localhost:443/token/', {
        method: 'DELETE',
    }).then(response => {
        if (response.ok) {
            nav('');
            console.log('Logged out');
        } else {
            console.log('failed logout');
        }
        return response;
    }).catch(err => {
        console.log(err);
        return null;
    })
}

// Updates the progress bar by percent.
const update = function updateProgressBar(percent) {
    document.getElementById('progress-bar-current').style['width'] = String(percent) + '%';
}

// Displays an element.
const show = function showElement(elementID) {
    document.getElementById(elementID).classList.remove('collapse');
}

// Hides an element.
const hide = function hideElement(elementID) {
    document.getElementById(elementID).classList.add('collapse');
}

// Changes the inner text of an element.
const change = function setInnerText(elementID, text) {
    document.getElementById(elementID).innerText = text;
}

// Sends fetch request to get Data.
const data = function getData() {
    return fetch('https://localhost:443/metrics', {
        method : 'GET',
        headers : {
            'Content-Type': 'application/json',
            'X-CSRF-Token' : document.getElementsByTagName('meta')['csrf-token'].content,
        },
    }).then(response => {
        if (response.ok) {
            return response.json();
        } else {
            throw new Error(response.statusText)
        }
    }).then(data => {
        const {metrics_count, plan} = data;
        let total = plan === 1 ? 1000 : 100;
        let percent = metrics_count/total;
        if (plan === 1) {
            dash();
        }
        change('filled-capacity', 'Users: ' + metrics_count +  ((plan === 1) ? '/1000' : '/100'));
        update(percent * 100);
        if (metrics_count === 100 && !(plan === 1)) {
            show('alert-error');
            show('upgrade-enterprise');
        }
        return data;
    }).catch(err => {
        console.log(err);
        nav('');
        return null;
    })
}

// Displays Enterprise Dash
const dash = function showEnterpriseDash() {
    change('plan-header', 'Enterprise Plan - $1000/Month');
    hide('upgrade-enterprise');
    hide('alert-error');
    show('alert-success');
}

// Polls at an interval of 1000 ms
const poll = function startLiveUpdate() {
    data();
    setInterval(data, 1000);
}

if (typeof module !== 'undefined') {
    module.exports = {
        logout: logout,
        nav: nav,
        upgrade: upgrade,
        submit: submit,
        data : data,
        show: show,
        hide: hide,
        update: update,
        dash: dash,
        change: change,
    }
}
