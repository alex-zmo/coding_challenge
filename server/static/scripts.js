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

if (typeof module !== 'undefined') {
    module.exports = {
        logout: logout,
        nav: nav,
        upgrade: upgrade,
        submit: submit,
    }
}
