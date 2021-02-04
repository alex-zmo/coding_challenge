const nav = (page) => {
    window.location.assign((page === '') ? '/' : '/' + page);
}

async function handleSubmit() {
    const user = {
        'username' : document.getElementById('email').value,
        'password' : document.getElementById('password').value,
    };
    const userJson = JSON.stringify(user);
    return await fetch('https://localhost:443/token/', {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: userJson,
    }).then(async function(response) {
        if (response.ok) {
            console.log('submit passed')
            document.getElementById('X-CSRF-Token-input').value = await response.json();
            await login();
        } else {
            console.log('submit failed')
        }
        return response
    }).catch(function(e) {
        return null;
    })
}

async function login() {
    return await fetch('https://localhost:443/token/', {
        method: 'GET',
        headers: {
            'X-CSRF-Token' : document.getElementById('X-CSRF-Token-input').value
        },
    }).then(async function(response) {
        if (response.ok) {
            console.log('login successful')
            nav('dashboard');
        } else {
            console.log('failed login')
        }
        return response;
    }).catch(function(e) {
        console.log('failed login');
        return null;
    })
}

async function upgradeEnterprise() {
    return await fetch('https://localhost:443/account/', {
        method: 'PATCH',
        headers : {
            'X-CSRF-Token' : document.getElementsByTagName('meta')['csrf-token'].content,
        },
        body: JSON.stringify({
            'plan' : 1,
        }),
    }).then(function(response) {
        if (response.ok) {
            console.log('upgraded');
        } else {
            console.log('upgrade failed')
            nav('');
        }
        return response;
    }).catch(function(e) {
        console.log(e);
        return null;
    })
}

async function logout() {
    return status = await fetch('https://localhost:443/token/', {
        method: 'DELETE',
    }).then(async function(response) {
        if (response.ok) {
            nav('');
            console.log('Logged out');
        } else {
            console.log('failed logout')
        }
        return response;
    }).catch(function (e) {
        console.log(e);
        return null;
    })
}

if(typeof module !== 'undefined') {
    module.exports = {
        logout: logout,
        login: login,
        nav: nav,
        upgradeEnterprise: upgradeEnterprise,
        handleSubmit: handleSubmit,
    }
}










