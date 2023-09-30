import static_ from '@fastify/static';
import ws from '@fastify/websocket';
import Fastify from 'fastify';
import path from 'path';
import * as url from 'url';

const __dirname = url.fileURLToPath(new URL('.', import.meta.url));
console.log(__dirname);

const app = Fastify({ logger: true });

app
  .register(ws)
  .register(static_, { root: path.join(__dirname, '..', 'public') });

app.get('/', async (_request, reply) => {
  return await reply.sendFile('index.html');
});

app.register(async function (fastify) {
  fastify.get('/ws', { websocket: true }, (connection, request) => {
    const { socket } = connection;

    console.log(
      `${request.ip} connecting to the server at ${new Date().toISOString()}`,
    );

    socket.on('message', (event) => {
      app.log.info(`Received message ${event} from ${request.ip}.`);
      socket.send('echo');
    });

    socket.on('close', (code, reason) => {
      console.log(`${code}, ${reason}`);
      app.log.info(`${request.ip} disconnected from server.`);
    });
  });
});

await app.listen({ port: 8080 });
