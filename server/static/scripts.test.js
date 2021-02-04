const fetchMock = require('fetch-mock');
const {logout, login, nav, upgradeEnterprise, handleSubmit} = require('./scripts.js');

const oldWindowLocation = window.location;

beforeAll(() => {
    delete window.location
    window.location = Object.defineProperties(
        {},
        {
            ...Object.getOwnPropertyDescriptors(oldWindowLocation),
            assign: {
                configurable: true,
                value: jest.fn(),
            },
        },
    )
})

beforeEach(() => {
    fetchMock.restore();
    window.location.assign.mockReset();
})

afterAll(() => {
    window.location = oldWindowLocation;
})

describe('testing api', () => {
    // HAPPY PATH
    it('calls nav and happy path', async() => {
        nav("test");
        expect(window.location.assign).toHaveBeenCalledTimes(1);
        expect(window.location.assign).toHaveBeenCalledWith('/test');
    });

    it('calls handleSubmit and returns 200 response', async () => {
    document.body.innerHTML = '<span id="email">random-user</span>' +
        '<span id="password">password</span>' +
        '<span id="X-CSRF-Token-input">123</span>';

        fetchMock.mock('https://localhost/token/', 200);
        console.log = jest.fn();
        const res = await handleSubmit();
        expect(fetchMock.lastCall().identifier).toBe("https://localhost/token/")
        expect(fetchMock.lastCall()[1].method).toBe("POST")
        expect(console.log.mock.calls[0][0]).toBe('submit passed');
    });

    it('calls login and returns 200 response', async () => {
        document.body.innerHTML = '<input id="X-CSRF-Token-input" name="X-CSRF-Token" />';
        document.getElementById("X-CSRF-Token-input").value = "random-token";
        fetchMock.mock('https://localhost/token/', 200);
        console.log = jest.fn();
        const res = await login();
        expect(res.status).toBe(200)
        expect(fetchMock.lastCall().identifier).toBe("https://localhost/token/")
        expect(fetchMock.lastCall()[1].method).toBe("GET")
        expect(console.log.mock.calls[0][0]).toBe('login successful');
    });

    it('calls upgradeEnterprise and returns 200 response', async () => {
        document.body.innerHTML = '<meta name="csrf-token" content="random-token">'
        fetchMock.mock('https://localhost/account/', 200);
        console.log = jest.fn();
        const res = await upgradeEnterprise();
        expect(res.status).toBe(200)
        expect(fetchMock.lastCall().identifier).toBe("https://localhost/account/")
        expect(fetchMock.lastCall()[1].method).toBe("PATCH")
        expect(console.log.mock.calls[0][0]).toBe('upgraded');
    });

    it('calls logout and returns 200 response', async () => {
        fetchMock.mock('https://localhost/token/', 200);
        console.log = jest.fn();
        const res = await logout();
        expect(res.status).toBe(200)
        expect(fetchMock.lastCall().identifier).toBe("https://localhost/token/")
        expect(fetchMock.lastCall()[1].method).toBe("DELETE")
        expect(console.log.mock.calls[0][0]).toBe('Logged out');
    });

// UNHAPPY PATHS
    it('calls handleSubmit and returns 500 response', async () => {
        document.body.innerHTML = '<span id="email">random-user</span>' +
            '<span id="password">password</span>' +
            '<span id="X-CSRF-Token-input">123</span>';

        fetchMock.mock('https://localhost/token/', 500);
        console.log = jest.fn();
        const res = await handleSubmit();
        expect(res.status).toBe(500)
        expect(fetchMock.lastCall().identifier).toBe("https://localhost/token/")
        expect(fetchMock.lastCall()[1].method).toBe("POST")
        expect(console.log.mock.calls[0][0]).toBe('submit failed');
    });

    it('calls login and returns 500 response', async () => {
        document.body.innerHTML = '<input id="X-CSRF-Token-input" name="X-CSRF-Token" />';
        document.getElementById("X-CSRF-Token-input").value = "random-token";
        fetchMock.mock('https://localhost/token/', 500);
        console.log = jest.fn();
        const res = await login();
        expect(res.status).toBe(500)
        expect(fetchMock.lastCall().identifier).toBe("https://localhost/token/")
        expect(fetchMock.lastCall()[1].method).toBe("GET")
        expect(console.log.mock.calls[0][0]).toBe('failed login');
    });

    it('calls upgradeEnterprise and returns 500 response', async () => {
        document.body.innerHTML = '<meta name="csrf-token" content="random-token">'
        fetchMock.mock('https://localhost/account/', 500);
        console.log = jest.fn();
        const res = await upgradeEnterprise();
        expect(res.status).toBe(500)
        expect(fetchMock.lastCall().identifier).toBe("https://localhost/account/")
        expect(fetchMock.lastCall()[1].method).toBe("PATCH")
        expect(console.log.mock.calls[0][0]).toBe('upgrade failed');
    });

    it('calls logout and returns 500 response', async () => {
        fetchMock.mock('https://localhost/token/', 500);
        console.log = jest.fn();
        const res = await logout();
        expect(res.status).toBe(500)
        expect(fetchMock.lastCall().identifier).toBe("https://localhost/token/")
        expect(fetchMock.lastCall()[1].method).toBe("DELETE")
        expect(console.log.mock.calls[0][0]).toBe('failed logout');

    })
});
