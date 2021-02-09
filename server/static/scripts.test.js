const fetchMock = require('fetch-mock');
const {logout, nav, upgrade, submit, data, update, show, hide, dash, change} = require('./scripts.js');

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
    it('calls nav and happy path', () => {
        nav('test');
        expect(window.location.assign).toHaveBeenCalledTimes(1);
        expect(window.location.assign).toHaveBeenCalledWith('/test');
    });

    it('calls show and happy path', () => {
        document.body.innerHTML =
            '<div>' +
            '  <span id="element" class="collapse button"/>' +
            '</div>';

        show("element");
        expect(document.getElementById('element').classList[0]).toBe("button");
    });

    it('calls hide and happy path', () => {
        document.body.innerHTML =
            '<div>' +
            '  <span id="element" />' +
            '</div>';

        hide("element");
        expect(document.getElementById('element').classList[0]).toBe("collapse");
    });

    it('calls change and happy path', () => {
        document.body.innerHTML =
            '<div>' +
            '  <span id="element"> before text </span>' +
            '</div>';

        change("element", "after text");
        expect(document.getElementById('element').innerText).toBe('after text');
    });

    it('calls update and happy path', () => {
        document.body.innerHTML =
            '<div>' +
            '  <span id="progress-bar-current"></span>' +
            '</div>';

        update(10);
        expect(document.getElementById('progress-bar-current').style["width"]).toBe('10%');
    });

    it('calls dash and happy path', () => {
        document.body.innerHTML =
            '<div>' +
            '   <span id="plan-header"></span>' +
            '   <span id="upgrade-enterprise"></span>' +
            '   <span id="alert-error"></span>' +
            '   <span id="alert-success" class="collapse button"></span>' +
            '</div>';

        dash();
        expect(document.getElementById('plan-header').innerText).toBe('Enterprise Plan - $1000/Month');
        expect(document.getElementById('upgrade-enterprise').classList[0]).toBe('collapse');
        expect(document.getElementById('alert-error').classList[0]).toBe('collapse');
        expect(document.getElementById('alert-success').classList[0]).toBe('button');

    });

    it('calls submit and returns 200 response', async () => {
        document.body.innerHTML = '<span id="email">random-user</span>' +
            '<span id="password">password</span>' +
            '<meta name="csrf-token" content="random-token">';

        fetchMock.mock('https://localhost/token/', 200);
        console.log = jest.fn();
        const res = await submit();
        expect(res.status).toBe(200)
        expect(fetchMock.lastCall().identifier).toBe('https://localhost/token/');
        expect(fetchMock.lastCall()[1].method).toBe('POST');
        expect(console.log.mock.calls[0][0]).toBe('submit passed');
    });

    it('calls upgrade and returns 200 response', async () => {
        document.body.innerHTML = '<meta name="csrf-token" content="random-token">'
        fetchMock.mock('https://localhost/account/', 200);
        console.log = jest.fn();
        const res = await upgrade();
        expect(res.status).toBe(200)
        expect(fetchMock.lastCall().identifier).toBe('https://localhost/account/');
        expect(fetchMock.lastCall()[1].method).toBe('PATCH');
        expect(console.log.mock.calls[0][0]).toBe('upgraded');
    });

    it('calls data and returns 200 response', async () => {
        document.body.innerHTML = '<meta name="csrf-token" content="random-token">' +
            '<div>' +
            '   <span id="plan-header"></span>' +
            '   <span id="upgrade-enterprise"></span>' +
            '   <span id="alert-error"></span>' +
            '   <span id="alert-success" class="collapse button"></span>' +
            '   <span id="filled-capacity"></span>' +
            '   <span id="progress-bar-current"></span>' +
            '</div>';
        fetchMock.mock('https://localhost/metrics', JSON.stringify({status: 200, metrics_count: 50, plan: 0}));
        const res = await data();
        expect(res.status).toBe(200);
        expect(res.metrics_count).toBe(50);
        expect(res.plan).toBe(0);
        expect(fetchMock.lastCall().identifier).toBe('https://localhost/metrics');
        expect(fetchMock.lastCall()[1].method).toBe('GET');
    });

    it('calls logout and returns 200 response', async () => {
        fetchMock.mock('https://localhost/token/', 200);
        console.log = jest.fn();
        const res = await logout();
        expect(res.status).toBe(200)
        expect(fetchMock.lastCall().identifier).toBe('https://localhost/token/');
        expect(fetchMock.lastCall()[1].method).toBe('DELETE');
        expect(console.log.mock.calls[0][0]).toBe('Logged out');
    });

// UNHAPPY PATHS
    it('calls submit and returns 500 response', async () => {

        document.body.innerHTML = '<meta name="csrf-token" content="random-token">' +
            '<span id="email">random-user</span>' +
            '<span id="password">password</span>';

        fetchMock.mock('https://localhost/token/', 500);
        console.log = jest.fn();
        const res = await submit();
        expect(res.status).toBe(500)
        expect(fetchMock.lastCall().identifier).toBe('https://localhost/token/');
        expect(fetchMock.lastCall()[1].method).toBe('POST');
        expect(console.log.mock.calls[0][0]).toBe('submit failed');
    });

    it('calls upgrade and returns 500 response', async () => {
        document.body.innerHTML = '<meta name="csrf-token" content="random-token">'
        fetchMock.mock('https://localhost/account/', 500);
        console.log = jest.fn();
        const res = await upgrade();
        expect(res.status).toBe(500)
        expect(fetchMock.lastCall().identifier).toBe('https://localhost/account/');
        expect(fetchMock.lastCall()[1].method).toBe('PATCH');
        expect(console.log.mock.calls[0][0]).toBe('upgrade failed');
    });

    it('calls data and returns 500 response', async () => {
        document.body.innerHTML = '<meta name="csrf-token" content="random-token">' +
            '<div>' +
            '   <span id="plan-header"></span>' +
            '   <span id="upgrade-enterprise"></span>' +
            '   <span id="alert-error"></span>' +
            '   <span id="alert-success" class="collapse button"></span>' +
            '   <span id="filled-capacity"></span>' +
            '   <span id="progress-bar-current"></span>' +
            '</div>';
        fetchMock.mock('https://localhost/metrics', 500);
        const res = await data();
        expect(res).toBe(null);
        expect(fetchMock.lastCall().identifier).toBe('https://localhost/metrics');
        expect(fetchMock.lastCall()[1].method).toBe('GET');
    });

    it('calls logout and returns 500 response', async () => {
        fetchMock.mock('https://localhost/token/', 500);
        console.log = jest.fn();
        const res = await logout();
        expect(res.status).toBe(500)
        expect(fetchMock.lastCall().identifier).toBe('https://localhost/token/');
        expect(fetchMock.lastCall()[1].method).toBe('DELETE');
        expect(console.log.mock.calls[0][0]).toBe('failed logout');

    });
});
