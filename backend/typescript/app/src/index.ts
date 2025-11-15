import Fastify from "fastify";

const server = Fastify();

server.get("/healthz", async () => ({ status: "ok" }));

server.get("/api/info", async () => ({
  name: "Social Scale Frontend Stub",
  version: process.env.APP_VERSION ?? "dev",
}));

async function start() {
  const port = Number(process.env.PORT ?? 3000);
  const host = "0.0.0.0";
  try {
    await server.listen({ port, host });
    console.log(`frontend listening on http://${host}:${port}`);
  } catch (err) {
    server.log.error(err);
    process.exit(1);
  }
}

start();


