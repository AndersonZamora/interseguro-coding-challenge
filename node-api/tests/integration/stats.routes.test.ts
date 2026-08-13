import request from 'supertest';
import { createApp } from '../../src/app';

const app = createApp();

function getToken() {
  return request(app)
    .post('/auth/token')
    .send({ client_id: process.env.DEMO_CLIENT_ID, client_secret: process.env.DEMO_CLIENT_SECRET });
}

describe('POST /auth/token', () => {
  it('devuelve un access_token con credenciales validas', async () => {
    const res = await getToken();
    expect(res.status).toBe(200);
    expect(res.body.access_token).toEqual(expect.any(String));
    expect(res.body.expires_in).toBe(3600);
  });

  it('devuelve 401 con credenciales invalidas', async () => {
    const res = await request(app).post('/auth/token').send({ client_id: 'x', client_secret: 'y' });
    expect(res.status).toBe(401);
  });

  it('devuelve 400 con payload invalido', async () => {
    const res = await request(app).post('/auth/token').send({ client_id: 'x' });
    expect(res.status).toBe(400);
  });
});

describe('POST /api/v1/stats', () => {
  it('devuelve 401 sin token', async () => {
    const res = await request(app)
      .post('/api/v1/stats')
      .send({ q: [[1, 2]], r: [[3, 4]] });
    expect(res.status).toBe(401);
  });

  it('devuelve 200 con token valido y payload valido', async () => {
    const tokenRes = await getToken();
    const token = tokenRes.body.access_token as string;

    const res = await request(app)
      .post('/api/v1/stats')
      .set('Authorization', `Bearer ${token}`)
      .send({ q: [[1, 2], [3, 4]], r: [[5, 0], [0, 5]] });

    expect(res.status).toBe(200);
    expect(res.body).toMatchObject({
      max: 5,
      min: 0,
      sum: 1 + 2 + 3 + 4 + 5 + 0 + 0 + 5,
      anyDiagonal: true,
    });
  });

  it('devuelve 400 con payload invalido (falta r)', async () => {
    const tokenRes = await getToken();
    const token = tokenRes.body.access_token as string;

    const res = await request(app)
      .post('/api/v1/stats')
      .set('Authorization', `Bearer ${token}`)
      .send({ q: [[1, 2]] });

    expect(res.status).toBe(400);
  });

  it('devuelve 400 con matriz vacia', async () => {
    const tokenRes = await getToken();
    const token = tokenRes.body.access_token as string;

    const res = await request(app)
      .post('/api/v1/stats')
      .set('Authorization', `Bearer ${token}`)
      .send({ q: [], r: [[1]] });

    expect(res.status).toBe(400);
  });

  it('devuelve 401 con token invalido', async () => {
    const res = await request(app)
      .post('/api/v1/stats')
      .set('Authorization', 'Bearer not-a-real-token')
      .send({ q: [[1]], r: [[1]] });
    expect(res.status).toBe(401);
  });
});
